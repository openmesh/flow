package eventbus

import (
	"github.com/openmesh/flow"
	"sync"
)

func New() flow.EventBus {
	return &eventBus{
		subscribers: nil,
		rm:          sync.RWMutex{},
	}
}

type eventBus struct {
	subscribers map[string][]flow.Channel
	rm          sync.RWMutex
}

func (b *eventBus) Subscribe(topic string, ch flow.Channel) error {
	b.rm.Lock()
	if prev, found := b.subscribers[topic]; found {
		b.subscribers[topic] = append(prev, ch)
	} else {
		b.subscribers[topic] = append([]flow.Channel{}, ch)
	}
	b.rm.Unlock()
	return nil
}

func (b *eventBus) Publish(topic string, payload interface{}) error {
	b.rm.RLock()
	if chans, found := b.subscribers[topic]; found {
		channels := append([]flow.Channel{}, chans...)
		go func(ev flow.Event, channels []flow.Channel) {
			for _, ch := range channels {
				ch <- ev
			}
		}(flow.Event{Payload: payload}, channels)
	}
	b.rm.Unlock()
	return nil
}
