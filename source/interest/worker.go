package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

// var amqpURI = "amqp://guest:guest@rabbitmq:5672/"
var amqpURI = "amqp://guest:guest@localhost:5672/"

func (interestApp *InterestApplication) startReadingMessages() {
	fmt.Printf("Processing message from queue %s at %s:%s\n", interestApp.RabbitReadQueue, interestApp.RabbitHost, interestApp.RabbitPort)
	interestApp.LastMessages[0] = "dfdf"

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(
		// This config is based on this example: https://www.rabbitmq.com/tutorials/tutorial-two-go.html
		// It works as a simple queue.
		//
		// If you want to implement a Pub/Sub style service instead, check
		// https://watermill.io/pubsubs/amqp/#amqp-consumer-groups
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

	go process(messages)

	// go publishMessages(publisher)
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
}
func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

		if err := publisher.Publish("example.topic", msg); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
