package pulse

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"eventStore"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

type pulsarStore struct {
	serviceName string
	client      pulsar.Client
	opts        eventStore.Options
}

func Init(opts eventStore.Options) (eventStore.EventStore, error) {
	if opts.Secure && opts.TLSConfig == nil {
		return nil, errors.New("secure connection must have valid tls configuration")
	}

	p, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: opts.Address,
		//Authentication:,
	})

	if err != nil {
		log.Print("Unable to connect with Pulsar secrets. Failed with error: ", err)
		return nil, err
	}

	return &pulsarStore{client: p}, nil
}

func (s *pulsarStore) SetServiceName(name string) {
	s.serviceName = name
}

func (s *pulsarStore) GetServiceName() string {
	return s.serviceName
}

func (s *pulsarStore) Publish(topic string, message []byte) error {
	sn := s.GetServiceName()
	__topic__ := fmt.Sprintf("%s.%s", sn, topic) // eventRoot is: io.roava.serviceName, topic is whatever is passed.

	producer, err := s.client.CreateProducer(pulsar.ProducerOptions{
		Topic: __topic__,
		Name:  sn,
	})
	if err != nil {
		log.Println("Failed to create new producer with the following error", err)
		return err
	}

	id, e := producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload:   message,
		EventTime: time.Now(),
	})

	if e != nil {
		log.Println(e, " Failed to send message.")
		return e
	}

	log.Printf("Published message to %s id ==>> %s", __topic__, byteToHex(id.Serialize()))
	return nil
}

// Manually put the fqdn of your topics.
func (s *pulsarStore) Subscribe(topic string, handler eventStore.SubscriptionHandler) error {
	serviceName := s.GetServiceName()
	consumer, err := s.client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       topic,
		AutoDiscoveryPeriod:         0,
		SubscriptionName:            serviceName,
		Type:                        pulsar.Shared,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
		Name:                        serviceName,
	})
	if err != nil {
		fmt.Print(err, " Subscribing to unknown topic.")
		return err
	}

	for {
		if val, ok := <-consumer.Chan(); ok {
			event := NewEvent(val)
			// TODO: Ensure event struct is according to the Roava Ecosystem.
			go handler(event)
			// TODO: Decide if we want to add something to stream this data to another place as backup.
		}
	}
}

func byteToHex(b []byte) string {
	var out struct{}
	_ = json.Unmarshal(b, &out)
	return hex.EncodeToString(b)
}
