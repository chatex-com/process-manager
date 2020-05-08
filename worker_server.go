package process

import (
	"context"
	"errors"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	ErrServerIsAlreadyRun = errors.New("server is already run")
	ErrServerIsNotRun = errors.New("server isn't run")
)

type ServerWorker struct {
	name      string
	server    *http.Server
	isRunning bool
	lock      sync.Locker
}

func NewServerWorker(name string, server *http.Server) *ServerWorker {
	w := ServerWorker{
		name:   name,
		server: server,
		lock:   &sync.Mutex{},
	}

	return &w
}

func (w *ServerWorker) Start() error {
	w.lock.Lock()

	if w.isRunning {
		w.lock.Unlock()

		return ErrServerIsAlreadyRun
	}

	w.isRunning = true

	w.lock.Unlock()

	log.WithFields(log.Fields{
		"listen": w.server.Addr,
		"worker": w.name,
	}).Info("start server worker")

	err := w.server.ListenAndServe()

	log.WithField("worker", w.name).Info("server worker has been stopped")

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (w *ServerWorker) Stop() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if !w.isRunning {
		return ErrServerIsNotRun
	}

	err := w.server.Shutdown(context.Background())
	if err == http.ErrServerClosed {
		err = nil
	}

	w.isRunning = false

	log.WithField("worker", w.name).Info("worker has got signal for stopping")

	return err
}
