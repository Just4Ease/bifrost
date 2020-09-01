package eventStore

type SubscriptionHandler func(event Event)

type EventStore interface {
	Publish(topic string, message []byte) error
	Subscribe(topic string, handler SubscriptionHandler) error
	SetServiceName(name string)
	GetServiceName() string
}
