package process

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Manager struct {
	wg            sync.WaitGroup
	workers       []Worker
	workersLock   sync.RWMutex
	isRunning     bool
	isRunningLock sync.RWMutex
}

func NewManagerWithWorkers(workers []Worker) *Manager {
	manager := &Manager{
		workers: workers,
	}

	return manager
}

func NewManager() *Manager {
	return NewManagerWithWorkers([]Worker{})
}

func (m *Manager) AddWorker(worker Worker) {
	m.workersLock.Lock()
	defer m.workersLock.Unlock()

	m.workers = append(m.workers, worker)

	if m.IsRunning() {
		m.wg.Add(1)
		go m.startWorker(worker)
	}
}

func (m *Manager) StartAll() {
	m.isRunningLock.Lock()
	defer m.isRunningLock.Unlock()

	if m.isRunning {
		log.WithError(errors.New("manager is already run")).Error("manager is already run")
		return
	}

	for _, w := range m.workers {
		m.wg.Add(1)
		go m.startWorker(w)
	}

	m.isRunning = true
}

func (m *Manager) startWorker(w Worker) {
	defer m.StopAll()

	err := w.Start()
	if err != nil {
		log.WithError(err).Error("the channel raised an error")
	}
}

func (m *Manager) stop() bool {
	m.isRunningLock.Lock()
	defer m.isRunningLock.Unlock()

	if !m.isRunning {
		return false
	}

	m.isRunning = false

	return true
}

func (m *Manager) IsRunning() bool {
	m.isRunningLock.RLock()
	defer m.isRunningLock.RUnlock()

	return m.isRunning
}

func (m *Manager) StopAll() {
	if !m.stop() {
		return
	}

	for _, w := range m.workers {
		err := w.Stop()
		if err != nil {
			log.WithError(err).Error("the channel raised an error")
		}
		m.wg.Done()
	}
}

func (m *Manager) AwaitAll() {
	m.wg.Wait()
	log.Info("all background processes were stopped")
}
