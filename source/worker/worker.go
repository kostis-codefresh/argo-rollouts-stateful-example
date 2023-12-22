package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sethvargo/go-retry"
)

func (interestApp *InterestApplication) retryConnecting() {
	ctx := context.Background()
	if err := retry.Fibonacci(ctx, 1*time.Second, func(ctx context.Context) error {
		if err := interestApp.startReadingMessages(); err != nil {
			// This marks the error as retryable
			fmt.Println("Cannot connect, will try later ", err)
			return retry.RetryableError(err)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func (interestApp *InterestApplication) startReadingMessages() (err error) {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", interestApp.RabbitHost, interestApp.RabbitPort)
	fmt.Printf("Connecting to %s:%s\n", interestApp.RabbitHost, interestApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NopLogger{},
	)
	if err != nil {
		return err
	}

	messages, err := subscriber.Subscribe(context.Background(), interestApp.RabbitReadQueue)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	interestApp.stopNow = cancelFunc

	go interestApp.process(messages, ctx)
	fmt.Printf("Ready to receive messages at %s\n", interestApp.RabbitReadQueue)
	return nil
}

func (interestApp *InterestApplication) publishMessage() {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", interestApp.RabbitHost, interestApp.RabbitPort)
	fmt.Printf("Sending dummy message on queue %s at %s:%s\n", interestApp.RabbitReadQueue, interestApp.RabbitHost, interestApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NopLogger{})
	if err != nil {
		fmt.Println("Could not connect to queue", err)
		return
	}

	//Just to distinguish messages from each other show the time that each message was sent
	t := time.Now()
	messageText := fmt.Sprintf("Dummy message sent at %s", t.Format("15:04:05"))

	msg := message.NewMessage(watermill.NewUUID(), []byte(messageText))

	if err := publisher.Publish(interestApp.RabbitReadQueue, msg); err != nil {
		fmt.Println("Could not publish message", err)
		return
	}
	interestApp.dummyCounter++
}

func (interestApp *InterestApplication) process(messages <-chan *message.Message, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("New configuration loaded - stopped reading from old queue")
			return
		case msg := <-messages:
			fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

			interestApp.mu.RLock()
			interestApp.LastMessages.PushFront(msg.Payload)
			if interestApp.LastMessages.Len() > 5 {
				oldest := interestApp.LastMessages.Back()
				interestApp.LastMessages.Remove(oldest)
			}

			// we need to Acknowledge that we received and processed the message,
			// otherwise, it will be resent over and over again.
			msg.Ack()
			interestApp.MessagesProcessed++
			interestApp.mu.RUnlock()
		}
	}
}
