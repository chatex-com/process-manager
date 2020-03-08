package process_manager

type Worker interface {
	// Run worker process in a current goroutine. An implementation
	// of this method mustn't return the execution context to the caller until finishing
	Start() error

	// Graceful stop of the worker process. It have to finish all active
	// processes and return execution context to the caller.
	Stop() error
}
