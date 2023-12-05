package main

import (
	"log"

	"github.com/nickyrolly/ws-chat-demo/internal/repository/chat_nsq"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/postgre"
	"github.com/nickyrolly/ws-chat-demo/internal/router"
)

func main() {
	err := postgre.InitPostgreSqltDB(postgre.DBMaster, "")
	if err != nil {
		log.Fatalf("Failed to init PostgreDB: %v", err)
	}

	err = chat_nsq.InitNSQProducer()
	if err != nil {
		log.Fatalf("Failed to init NSQ Producer: %v", err)
	}

	err = chat_nsq.InitNSQConsumer()
	if err != nil {
		log.Fatalf("Failed to init NSQ Consumer: %v", err)
	}

	router.Init()
}
