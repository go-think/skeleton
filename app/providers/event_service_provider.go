package providers

import (
	"app/app/events"

	"github.com/go-think/think/container"
	"github.com/go-think/think/contract"
)

type EventServiceProvider struct{}

func (p *EventServiceProvider) Register(app *container.Container) {
}

func (p *EventServiceProvider) Boot(app *container.Container) {
	dispatcher := app.Make[contract.EventDispatcher]("contract.EventDispatcher")
	if dispatcher == nil {
		dispatcher = app.Make[contract.EventDispatcher]("events")
	}

	if dispatcher == nil {
		return
	}

	// Register user.registered event listener
	dispatcher.Listen("user.registered", func(payload interface{}) {
		events.Tracker.Record("user.registered:" + payload.(string))
	})

	// Register order.creating event listener (for Until)
	dispatcher.Listen("order.creating", func(payload interface{}) interface{} {
		if amount, ok := payload.(float64); ok && amount > 10000 {
			return "Amount exceeds risk threshold"
		}
		return nil
	})

	// Register kernel lifecycle listeners
	dispatcher.Listen("kernel.handled", func(payload interface{}) {
		events.Tracker.Record("kernel.handled")
	})

	dispatcher.Listen("kernel.terminating", func(payload interface{}) {
		events.Tracker.Record("kernel.terminating")
	})
}
