package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

// var amqpURI = "amqp://guest:guest@rabbitmq:5672/"
var amqpURI = "amqp://guest:guest@localhost:5672/"

func (interestApp *InterestApplication) startReadingMessages() {
	fmt.Printf("Connecting to %s:%s\n", interestApp.RabbitHost, interestApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NopLogger{},
	)
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), "example.topic")
	if err != nil {
		panic(err)
	}

	go interestApp.process(messages)
	fmt.Printf("Ready to receive messages at %s\n", interestApp.RabbitReadQueue)
}

func (interestApp *InterestApplication) publishMessage() {
	fmt.Printf("Sending dummy message on queue %s at %s:%s\n", interestApp.RabbitReadQueue, interestApp.RabbitHost, interestApp.RabbitPort)
	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NopLogger{})
	if err != nil {
		panic(err)
	}
	msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

	if err := publisher.Publish("example.topic", msg); err != nil {
		panic(err)
	}
	interestApp.dummyCounter++
}

func (interestApp *InterestApplication) process(messages <-chan *message.Message) {
	interestApp.mu.RLock()
	defer interestApp.mu.RUnlock()
	for msg := range messages {
		fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		interestApp.LastMessages.PushBack(msg.Payload)
		if interestApp.LastMessages.Len() > 5 {
			oldest := interestApp.LastMessages.Front()
			interestApp.LastMessages.Remove(oldest)
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
		interestApp.MessagesProcessed++
	}
}
