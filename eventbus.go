package flow

type Event struct {
	Payload interface{}
	Topic   string
}

type Channel chan Event

type EventBus interface {
	Subscribe(topic string, ch Channel) error
	Publish(topic string, payload interface{}) error
}
