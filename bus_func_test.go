package gebus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/teawithsand/gebus"
)

// TODO(teawithsand): tests for invalid handlers

func TestFuncEventhandler(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		cs := callStats{}
		h := makeJobHandler(t, &cs)

		err := h.HandleEvent(context.Background(), &arg1{})
		if err != nil {
			t.Error(err)
		}

		err = h.HandleEvent(context.Background(), &arg2{})
		if err != nil {
			t.Error(err)
		}

		if cs.C1 != 1 {
			t.Error("Fn1 called invalid amount of times")
		}
		if cs.C2 != 1 {
			t.Error("Fn1 called invalid amount of times")
		}
	})

	t.Run("handles_not_existing", func(t *testing.T) {
		h := makeJobHandler(t, nil)
		err := h.HandleEvent(context.Background(), &argNotExist{})
		var verr *gebus.HandlerNotFound
		if !errors.As(err, &verr) {
			t.Error(err)
		}
	})
}
