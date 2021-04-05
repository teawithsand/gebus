package gebus_test

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/teawithsand/gebus"
)

type arg1 struct{}
type arg2 struct{}
type timedArg struct {
	Time time.Duration
}
type argNotExist struct{}

type callStats struct {
	C1 int32
	C2 int32
}

func makeJobHandler(t *testing.T, cs *callStats) *gebus.FuncEventBus {
	h := gebus.FuncEventBus{}

	h.EventHandlers = append(h.EventHandlers, reflect.ValueOf(func(ctx context.Context, arg *arg1) (err error) {
		if cs != nil {
			// cs.C1 += 1 // TODO(teawithsand): use atomic here
			atomic.AddInt32(&cs.C1, 1)
		}
		return
	}))
	h.EventHandlers = append(h.EventHandlers, reflect.ValueOf(func(ctx context.Context, arg *arg2) (err error) {
		if cs != nil {
			atomic.AddInt32(&cs.C2, 1)
		}
		return
	}))
	h.EventHandlers = append(h.EventHandlers, reflect.ValueOf(func(ctx context.Context, arg *timedArg) (err error) {
		time.Sleep(arg.Time)
		return
	}))

	err := h.CheckHandlers()
	if err != nil {
		t.Error(err)
		return nil
	}

	return &h
}
