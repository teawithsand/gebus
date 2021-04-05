package gebus

import (
	"context"
	"errors"
	"sync"
)

// WorkerPoolEventHandler runs up to specified amount of goroutines simultaneously.
// for underlying event handler.
//
// It never returns error from handler. Instead ErrorHandler should be used to handle errors if there is such requirement.
//
// If concurrent job limit is zero then infinite amount of goroutines is allowed.
type WorkerPoolEventHandler struct {
	Handler                  EventHandler
	ConcurrentJobLimit       int
	BackgroundHandlerContext context.Context
	ErrorHandler             func(ctx context.Context, data interface{}, err error)

	runVar      *sync.Cond
	runningJobs int
}

func (eh *WorkerPoolEventHandler) Initialize() (err error) {
	if eh.runVar != nil {
		err = errors.New("this WorkerPoolEventHandler is already initialized")
		return
	}

	eh.runVar = sync.NewCond(&sync.Mutex{})

	return

}

func (eh *WorkerPoolEventHandler) HandleEvent(ctx context.Context, event interface{}) (err error) {
	handler := eh.Handler
	backgroundContext := eh.BackgroundHandlerContext
	if backgroundContext == nil {
		backgroundContext = context.Background()
	}
	errorHandler := eh.ErrorHandler
	limit := eh.ConcurrentJobLimit

	// note: running jobs counter has to be incremented before call to this function
	runJob := func() {
		defer func() {
			eh.runVar.L.Lock()
			eh.runningJobs -= 1
			eh.runVar.L.Unlock()

			eh.runVar.Broadcast() // Is signal sufficient here always? Broadcast is safer to use, might slower though.
		}()

		err := handler.HandleEvent(backgroundContext, event)
		if err != nil && errorHandler != nil {
			errorHandler(backgroundContext, event, err)
		}
	}

	// note: this function requires eh.runVar.L to be acquited
	scheduleJob := func() {
		eh.runningJobs += 1
		go runJob()
	}

	if limit == 0 {
		eh.runVar.L.Lock()
		scheduleJob()
		eh.runVar.L.Unlock()
		return
	}

	doneChan := make(chan struct{})
	jobScheduled := false
	interrupted := false
	go func() {
	schedLoop:
		for {
			eh.runVar.L.Lock()
			if eh.runningJobs >= limit {
				eh.runVar.Wait()
				// wait can lock as well, so add interrupt check here, not above if
				if interrupted {
					// already interrupted, error was returned so do not schedule job now

					eh.runVar.L.Unlock()
					break schedLoop
				}

				eh.runVar.L.Unlock()
				continue
			} else {
				if interrupted {
					// already interrupted, error was returned so do not schedule job now

					eh.runVar.L.Unlock()
					break schedLoop
				}

				scheduleJob()
				jobScheduled = true

				eh.runVar.L.Unlock()
				break schedLoop
			}
		}
		doneChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		// notify all, so one which owns this context won't leak goroutine
		eh.runVar.L.Lock()
		if jobScheduled { // if race cond occurred and job was scheduled then do not return error from context timeout, because scheduling succeed
			err = nil
		} else {
			err = ctx.Err()
		}
		interrupted = true
		eh.runVar.L.Unlock()

		// notify all listeners so that one, which is responsible for this job is able to catch interrupted flag and
		// exit gorutine quite quickly without leaking it
		eh.runVar.Broadcast()
		return
	case <-doneChan:
		return
	}
}
