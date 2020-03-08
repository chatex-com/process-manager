package process_manager

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

	log.WithField("listen", w.server.Addr).Infof("%s: start server", w.name)

	err := w.server.ListenAndServe()
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

	log.Infof("%s: server has been stopped", w.name)

	return err
}
