package gebus

import "fmt"

type HandlerNotFound struct {
	Event interface{}
}

func (e *HandlerNotFound) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("event handler for event of type %T was not found", e.Event)
}
