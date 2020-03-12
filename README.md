# Process Manager

<a href="https://opensource.org/licenses/Apache-2.0" rel="nofollow"><img src="https://img.shields.io/badge/license-Apache%202-blue" alt="License" style="max-width:100%;"></a>
![unit-tests](https://github.com/chatex-com/process-manager/workflows/unit-tests/badge.svg)

## Example

```go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chatex-com/process-manager"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
	manager := process_manager.NewManager()

	// Create a callback worker
	manager.AddWorker(process_manager.NewCallbackWorker("test", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}))

	// Create a callback worker with retries
	callbackWorker := process_manager.NewCallbackWorker("test with error", func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Millisecond):
				return errors.New("test error")
			}
		}
	}, true)
	callbackWorker.Retries = 10 // When this param is missed manager will try restart in infinity loop
	callbackWorker.RetryTimeout = time.Second // When it is omitted manager will try to run it immediately 
	manager.AddWorker(callbackWorker)


	// Create an example of server worker for prometheus
	handler := mux.NewRouter()
	handler.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:    ":2112",
		Handler: handler,
	}
	manager.AddWorker(process_manager.NewServerWorker("prometheus", server))

	manager.StartAll()

	WaitShutdown(manager)
}

func WaitShutdown(manager *process_manager.Manager) {
	go func(manager *process_manager.Manager) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan

		manager.StopAll()
	}(manager)

	manager.AwaitAll()
}

```
