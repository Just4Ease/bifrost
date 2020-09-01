package eventStore

type Event interface {
	Ack()
	Data() []byte
	Topic() string
}
