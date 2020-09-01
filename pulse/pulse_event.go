package pulse

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/roava/eventStore"
)

type event struct {
	raw *pulsar.ConsumerMessage
}

func NewEvent(message pulsar.ConsumerMessage) eventStore.Event {
	return &event{raw: &message}
}

func (e *event) Data() []byte {
	return e.raw.Payload()
}

func (e *event) Topic() string {
	t := e.raw.Topic()
	// Manually agree what we want the topic to look like from pulsar.
	return t
}

func (e *event) Ack() {
	e.raw.AckID(e.raw.ID())
}
