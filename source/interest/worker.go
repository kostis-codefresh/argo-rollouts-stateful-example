package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func (interestApp *InterestApplication) startReadingMessages() {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", interestApp.RabbitHost, interestApp.RabbitPort)
	fmt.Printf("Connecting to %s:%s\n", interestApp.RabbitHost, interestApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NopLogger{},
	)
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), interestApp.RabbitReadQueue)
	if err != nil {
		panic(err)
	}

	go interestApp.process(messages)
	fmt.Printf("Ready to receive messages at %s\n", interestApp.RabbitReadQueue)
}

func (interestApp *InterestApplication) publishMessage() {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", interestApp.RabbitHost, interestApp.RabbitPort)
	fmt.Printf("Sending dummy message on queue %s at %s:%s\n", interestApp.RabbitReadQueue, interestApp.RabbitHost, interestApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NopLogger{})
	if err != nil {
		panic(err)
	}

	//Just to distinguish messages from each other show the time that each message was sent
	t := time.Now()
	messageText := fmt.Sprintf("Dummy message sent at %s", t.Format("15:04:05"))

	msg := message.NewMessage(watermill.NewUUID(), []byte(messageText))

	if err := publisher.Publish(interestApp.RabbitReadQueue, msg); err != nil {
		panic(err)
	}
	interestApp.dummyCounter++
}

func (interestApp *InterestApplication) process(messages <-chan *message.Message) {
	interestApp.mu.RLock()
	defer interestApp.mu.RUnlock()
	for msg := range messages {
		fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		interestApp.LastMessages.PushFront(msg.Payload)
		if interestApp.LastMessages.Len() > 5 {
			oldest := interestApp.LastMessages.Back()
			interestApp.LastMessages.Remove(oldest)
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
		interestApp.MessagesProcessed++
	}
}
