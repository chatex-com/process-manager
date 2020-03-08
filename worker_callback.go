package process_manager

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type CallbackFunc func(ctx context.Context) error

type CallbackWorker struct {
	ctx       context.Context
	cancel    context.CancelFunc
	name      string
	cb        CallbackFunc
	isRunning bool
	lock      sync.Locker
}

func NewCallbackWorker(name string, cb CallbackFunc) *CallbackWorker {
	ctx, cancel := context.WithCancel(context.Background())

	return &CallbackWorker{
		ctx:    ctx,
		cancel: cancel,
		name:   name,
		cb:     cb,
		lock:   &sync.Mutex{},
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

	log.Infof("start %s worker", w.name)

	err := w.cb(w.ctx)
	if err != nil && err != context.Canceled {
		return err
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

	log.Infof("worker %s has been stopped", w.name)

	return nil
}
