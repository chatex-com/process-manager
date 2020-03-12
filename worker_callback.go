package process_manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type CallbackFunc func(ctx context.Context) error

type CallbackWorker struct {
	ctx          context.Context
	cancel       context.CancelFunc
	name         string
	cb           CallbackFunc
	isRunning    bool
	lock         sync.Locker
	RetryOnError bool
	Retries      uint
	RetryTimeout time.Duration
	errors       uint
}

func NewCallbackWorker(name string, cb CallbackFunc, retryOnError ...bool) *CallbackWorker {
	ctx, cancel := context.WithCancel(context.Background())

	var retry bool
	if len(retryOnError) > 0 {
		retry = retryOnError[0]
	}

	return &CallbackWorker{
		ctx:          ctx,
		cancel:       cancel,
		name:         name,
		cb:           cb,
		lock:         &sync.Mutex{},
		RetryOnError: retry,
	}
}

func (w *CallbackWorker) Start() error {
	w.lock.Lock()
	if w.isRunning {
		w.lock.Unlock()
		return fmt.Errorf("worker %s is already run", w.name)
	}

	w.isRunning = true
	w.lock.Unlock()

	log.WithField("worker", w.name).Info("start callback worker")

	for w.isRunning {
		err := w.cb(w.ctx)

		if err == nil {
			return nil
		}

		if err == context.Canceled {
			return nil
		}

		w.errors++
		if w.Retries > 0 && w.errors >= w.Retries {
			return err
		}

		if !w.RetryOnError {
			return err
		}

		log.WithError(err).WithField("worker", w.name).Error("retrying execution of callback during error")
		<-time.After(w.RetryTimeout)
	}

	return nil
}

func (w *CallbackWorker) Stop() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if !w.isRunning {
		return fmt.Errorf("worker %s isn't run", w.name)
	}

	w.cancel()
	w.isRunning = false

	log.WithField("worker", w.name).Info("worker has been stopped")

	return nil
}
