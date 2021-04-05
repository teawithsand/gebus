package gebus

import (
	"context"
	"errors"
	"sync"
)

// ACKingAdapter is adapter, whcih uses receive-and-ack semantics.
//
// Once job is received and then handled, ACK is send to remote server to notify that handling was success
// and there is no need to reschedule this job later / on different server.
// This is used to make sure that no event is ever lost and all handlers were executed.
//
// Note: adapters are not guaranteed to be thread safe.
// Each goroutine should have it's own adapter.
type ACKingAdapter interface {
	// Requests next event from server adapter.
	// Should not block longer than until context times out.
	NextEvent(ctx context.Context) (event interface{}, err error)

	// Returns value of current evnet.
	// Nil if there is no any.
	GetCurrentEvent() interface{}

	// ACKPositive notifies server that event was handled properly and there is no need to reschedule it for later handling.
	ACKPositive() (err error)

	// Close closes adapter finalizing stuff associated with current job data(if any).
	Close() (err error)
}

// ACKingServer takes ACKingAdapter and makes it into working server.
//
// Note: EventHandler passed to this server should be single threaded and synchronous.
// If parallelism is required simply run many servers using many adapters.
// The reason for this is that ACKing has to be done once job is handled, since during execution there is still risk of some
// failure like power loss.
//
// Note #2: Returning error from event handler causes server to stop, so this should not be done too often.
// Also please note that context error's MUST be returned from handler if some occur.
type ACKingServer struct {
	Adapter ACKingAdapter
	Handler EventHandler

	isClosed bool
	lock     *sync.Mutex
}

func (srv *ACKingServer) Initialize() (err error) {
	if srv.lock != nil {
		err = errors.New("server already initialized")
		return
	}

	return
}

func (srv *ACKingServer) RunServer(ctx context.Context) (err error) {
	// Lock is used to make sure that single instance runs at a time.
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.isClosed {
		err = errors.New("server was already closed. Use new instance to rerun it")
		return
	}

	defer srv.Adapter.Close()
	defer func() {
		srv.isClosed = true
	}()

	for {
		var event interface{}
		event, err = srv.Adapter.NextEvent(ctx)
		if err != nil {
			return
		}

		err = srv.Handler.HandleEvent(ctx, event)
		if err != nil {
			return
		}

		err = srv.Adapter.ACKPositive()
		if err != nil {
			return
		}
	}
}
