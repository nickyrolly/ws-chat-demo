package chat_nsq

import (
	"fmt"
	"log"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
	"github.com/nsqio/go-nsq"
)

var NSQProducer *nsq.Producer

type ConsumerStruct struct {
	Topic    string
	Channel  string
	Function func(message *nsq.Message) error
}

var consumerList = []ConsumerStruct{
	//Exercise 3.1.3
	//Register Topic name and Consumer for Chat Personal History
	//--

	//Exercise 3.2.3
	//Register Topic name and Consumer for Group Personal History
	//--
}

func InitNSQProducer() error {
	// Create an NSQ producer to publish message
	producer, err := nsq.NewProducer(domain.HostNSQd, nsq.NewConfig())
	if err != nil {
		return fmt.Errorf("Error creating producer aliyun on %s", domain.HostNSQd)
	}

	NSQProducer = producer

	return nil
}

func InitNSQConsumer() error {
	for _, consumerData := range consumerList {
		// Create an NSQ consumer to consume the message
		consumer, err := nsq.NewConsumer(consumerData.Topic, consumerData.Channel, nsq.NewConfig())
		if err != nil {
			log.Fatalf("Failed to create NSQ consumer: %v", err)
		}

		// Configure the message handler for the consumer
		consumer.AddHandler(nsq.HandlerFunc(consumerData.Function))

		// Connect the consumer to the NSQD server
		if err := consumer.ConnectToNSQD(domain.HostNSQd); err != nil {
			log.Fatalf("Failed to connect to NSQD: %v", err)
		}
	}

	return nil
}

func consumerSaveChatHistory(message *nsq.Message) error {
	//Exercise 3.1.4
	//Create consumer handler for Inserting Chat Personal History Data

	return nil
}

func consumerGroupSaveChatHistory(message *nsq.Message) error {
	//Exercise 3.2.4
	//Create consumer handler for Inserting Chat Group History Data

	return nil
}
