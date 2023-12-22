package main

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func (testerApp *TesterApplication) publishProductionMessage() {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", testerApp.RabbitHost, testerApp.RabbitPort)
	fmt.Printf("Sending production message on queue %s at %s:%s\n", testerApp.RabbitQueue, testerApp.RabbitHost, testerApp.RabbitPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NopLogger{})
	if err != nil {
		panic(err)
	}

	//Just to distinguish messages from each other show the time that each message was sent
	t := time.Now()
	messageText := fmt.Sprintf("Production message sent at %s", t.Format("15:04:05"))

	msg := message.NewMessage(watermill.NewUUID(), []byte(messageText))

	if err := publisher.Publish(testerApp.RabbitQueue, msg); err != nil {
		panic(err)
	}
	testerApp.ProductionMessagesSent++
}

func (testerApp *TesterApplication) publishPreviewMessage() {
	//Format is "amqp://guest:guest@rabbitmq:5672/"
	amqpURI := fmt.Sprintf("amqp://guest:guest@%s:%s", testerApp.RabbitPreviewHost, testerApp.RabbitPreviewPort)
	fmt.Printf("Sending production message on queue %s at %s:%s\n", testerApp.RabbitPreviewQueue, testerApp.RabbitPreviewHost, testerApp.RabbitPreviewPort)

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NopLogger{})
	if err != nil {
		panic(err)
	}

	//Just to distinguish messages from each other show the time that each message was sent
	t := time.Now()
	messageText := fmt.Sprintf("Preview message sent at %s", t.Format("15:04:05"))

	msg := message.NewMessage(watermill.NewUUID(), []byte(messageText))

	if err := publisher.Publish(testerApp.RabbitPreviewQueue, msg); err != nil {
		panic(err)
	}
	testerApp.PreviewMessagesSent++
}
