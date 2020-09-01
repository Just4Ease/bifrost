package pulse

import (
	"github.com/roava/eventStore"
	"testing"
	"time"
)

var topic = "test-topic-xxx" // for random topic name.
var store, _ = Init(eventStore.Options{
	Address: "pulsar://localhost:6650",
})

func TestStore_Publish(t *testing.T) {
	store.SetServiceName("test-service")
	if err := store.Publish(topic, []byte("Hello World!")); err != nil {
		t.Errorf("Failed to publish data to event store topic %s. Failed with error: %v", topic, err)
	}
}

func TestStore_Subscribe(t *testing.T) {
	store.SetServiceName("test-service")
	timer := time.AfterFunc(3*time.Second, func() {
		if err := store.Subscribe(topic, func(event eventStore.Event) {
			data := event.Data()

			eventTopic := event.Topic()

			if topic != eventTopic {
				t.Errorf("Event topic is not the same as subscription topic. Why?: Expected %s, instead got: %s \n", topic, eventTopic)
				return
			}

			t.Logf("Received data: %s on topic: %s \n", string(data), eventTopic)
			event.Ack() // Acknowledge event.
			return
		}); err != nil {
			t.Errorf("Failed to subscribe to topic: %s, with the following error: %v \n", topic, err)
			return
		}
	})

	defer timer.Stop()
}
