package gebus

import "context"

// EventHandler is something capable of handling events.
// It may be handler itself or event bus.
type EventHandler interface {
	HandleEvent(ctx context.Context, event interface{}) (err error)
}

type EventHandlerFunc func(ctx context.Context, event interface{}) (err error)

func (f EventHandlerFunc) HandleEvent(ctx context.Context, event interface{}) (err error) {
	return f(ctx, event)
}
