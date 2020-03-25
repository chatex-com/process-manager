package process

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewManager(t *testing.T) {
	Convey("Create empty manager", t, func() {
		manager := NewManager()

		So(manager, ShouldNotBeNil)
		So(manager, ShouldHaveSameTypeAs, &Manager{})
		So(manager.isRunning, ShouldBeFalse)
		So(manager.workers, ShouldBeEmpty)
	})

	Convey("Create manager with workers", t, func() {
		worker := NewCallbackWorker("test", func(ctx context.Context) error {
			return nil
		})
		manager := NewManagerWithWorkers([]Worker{worker})

		So(manager, ShouldNotBeNil)
		So(manager.isRunning, ShouldBeFalse)
		So(manager.workers, ShouldHaveLength, 1)
		So(manager.workers[0], ShouldEqual, worker)
	})
}

func TestAddWorker(t *testing.T) {
	Convey("Adding worker to stopped manager", t, func() {
		manager := NewManager()

		So(manager.workers, ShouldBeEmpty)
		worker := NewCallbackWorker("test", func(ctx context.Context) error {
			return nil
		})

		manager.AddWorker(worker)

		So(manager.workers, ShouldHaveLength, 1)
		So(manager.workers[0], ShouldEqual, worker)
	})

	Convey("Adding worker to run manager", t, func() {
		manager := NewManager()
		manager.StartAll()

		So(manager.workers, ShouldBeEmpty)
		worker := NewCallbackWorker("test", func(ctx context.Context) error {
			return nil
		})

		manager.AddWorker(worker)

		So(manager.workers, ShouldHaveLength, 1)
		So(manager.workers[0], ShouldEqual, worker)
	})
}
