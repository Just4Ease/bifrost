package bifrost

type Event interface {
	Ack()
	Data() []byte
	Topic() string
}
