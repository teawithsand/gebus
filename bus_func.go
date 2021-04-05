package gebus

import (
	"context"
	"errors"
	"reflect"
)

// FuncEventBus uses functions of reflect.Value type in order to handle
// incoming event depending on data type.
//
// Event handler will be called once matching one is found.
// Match is triggered when event value is assignable to 2nd function argument.
//
// Note: If some event matches more than one handler then first one is picked.
type FuncEventBus struct {
	EventHandlers []reflect.Value
}

// CheckHandlers checks handlers and panics if some is invalid.
func (h *FuncEventBus) CheckHandlers() (err error) {
	for _, v := range h.EventHandlers {
		if v.Kind() != reflect.Func {
			err = errors.New("kind other than func provided")
			return
		}
		if v.Type().NumIn() != 2 {
			err = errors.New("function provided does not accept two arguments")
			return
		}
		if v.Type().NumOut() != 1 {
			err = errors.New("function provided does not return one value")
			return
		}
		if !reflect.TypeOf((*context.Context)(nil)).Elem().AssignableTo(v.Type().In(0)) {
			err = errors.New("can't assign context to 1st argument")
			return
		}
		// 2nd arg may be any type
		/*
			if v.Type().In(1).Kind() != reflect.Ptr || v.Type().In(1).Elem().Kind() != reflect.Struct {
				err = errors.New("second argument is not pointer to struct type")
				return
			}
		*/
		if !v.Type().Out(0).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			err = errors.New("can't assign 1st output value to error")
			return
		}
	}
	return
}

func (h *FuncEventBus) HandleEvent(ctx context.Context, event interface{}) (err error) {
	et := reflect.TypeOf(event)
	for _, reflectHandler := range h.EventHandlers {
		// do not run check here
		// reflect will panic anyway
		/*
			if reflectHandler.Kind() != reflect.Func {
				panic("invalid job handler type here; Call CheckHandlers to catch this panic befor running HandleEvent")
			}
		*/

		ty := reflectHandler.Type().In(1)
		if et.AssignableTo(ty) {
			reflectHandler.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(event)})
			return
		}
	}

	err = &HandlerNotFound{Event: event}
	return
}
