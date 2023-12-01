package usecase

import (
	"fmt"
	"log"
	"net"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/chat_nsq"
	"github.com/nsqio/go-nsq"
)

func ServiceHealthCheck(result *map[string]string, serviceName, address string) {
	fmt.Printf("Checking %s...\n", serviceName)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("%s is unreachable: %v\n", serviceName, err)
		(*result)[serviceName] = "Failed"
		return
	}
	conn.Close()
	fmt.Printf("%s is reachable\n", serviceName)
	(*result)[serviceName] = "Success"
}

func PublishHealthCheck(healthCheckResult *map[string]string) {
	// Publish a test message
	(*healthCheckResult)["PublishNSQ"] = "Success"
	messageBody := []byte("Hello, NSQ!")
	if err := chat_nsq.NSQProducer.Publish("test-topic", messageBody); err != nil {
		(*healthCheckResult)["PublishNSQ"] = "Failed"
	}
}

func NsqHealthCheck() {
	// Create an NSQ consumer to consume the test message
	consumer, err := nsq.NewConsumer("test-topic", "test-channel", nsq.NewConfig())
	if err != nil {
		log.Fatalf("Failed to create NSQ consumer: %v", err)
	}

	// Configure the message handler for the consumer
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		fmt.Printf("Received message: %s\n", message.Body)
		message.Finish()
		return nil
	}))

	// Connect the consumer to the NSQD server
	if err := consumer.ConnectToNSQD(domain.HostNSQd); err != nil {
		log.Fatalf("Failed to connect to NSQD: %v", err)
	}

	// Block indefinitely to keep the health check running
	select {}
}
