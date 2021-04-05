package gebus

/*
// CallAllBus runs all owned handlers on each event until some returns non-nil error.
type CallAllBus struct {
	Handlers []EventHandler
}

func (cab *CallAllBus) HandleEvent(ctx context.Context, event interface{}) (err error) {
	for _, h := range cab.Handlers {
		err = h.HandleEvent(ctx, event)
		if err != nil {
			return
		}
	}
	return
}
*/
