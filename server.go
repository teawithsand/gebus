package gebus

import (
	"context"
)

// Some handlers require servers to work because producers and consumers of events are splitted either logically
// or even work on separate machines.
//
// For this purpose on client there exists ClientEventHandler and on server there is Server, which takes incoming jobs
// and handles them using underlying event handler.
//
// HandlerServer is interface abstracting away single instance of any server.
type HandlerServer interface {
	// Run runs server infinitely until context times out or something crashes.
	RunServer(ctx context.Context) (err error)
}
