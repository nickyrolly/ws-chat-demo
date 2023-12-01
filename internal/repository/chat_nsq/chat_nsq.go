package chat_nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/postgre"
	"github.com/nsqio/go-nsq"
)

var NSQProducer *nsq.Producer

type ConsumerStruct struct {
	Topic    string
	Channel  string
	Function func(message *nsq.Message) error
}

var consumerList = []ConsumerStruct{
	{
		Topic:    "save-chat-history-topic",
		Channel:  "chat-channel",
		Function: consumerSaveChatHistory,
	},
	{
		Topic:    "save-group-chat-history-topic",
		Channel:  "chat-channel",
		Function: consumerGroupSaveChatHistory,
	},
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
	var (
		data repository.ChatHistoryData
		err  error
	)
	err = json.Unmarshal(message.Body, &data)
	if err != nil {
		log.Println("Error Unmarshal: ", err.Error())
		message.Finish()
		return err
	}

	err = postgre.InsertChatHistory(context.Background(), data)
	if err != nil {
		log.Println("Error InsertChatHistory: ", err.Error())
		message.Requeue(1)
		return err
	}

	message.Finish()
	return nil
}

func consumerGroupSaveChatHistory(message *nsq.Message) error {
	var (
		data repository.GroupChatHistoryData
		err  error
	)
	err = json.Unmarshal(message.Body, &data)
	if err != nil {
		log.Println("Error Unmarshal: ", err.Error())
		message.Finish()
		return err
	}

	err = postgre.InsertGroupChatHistory(context.Background(), data)
	if err != nil {
		log.Println("Error InsertChatHistory: ", err.Error())
		message.Requeue(1)
		return err
	}

	message.Finish()
	return nil
}
