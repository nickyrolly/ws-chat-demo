package main

import (
	"github.com/nickyrolly/ws-chat-demo/internal/router"
	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

func main() {
	// err := postgre.InitPostgreSqltDB(postgre.DBMaster, "")
	// if err != nil {
	// 	log.Fatalf("Failed to init PostgreDB: %v", err)
	// }

	// err = chat_nsq.InitNSQProducer()
	// if err != nil {
	// 	log.Fatalf("Failed to init NSQ Producer: %v", err)
	// }

	// err = chat_nsq.InitNSQConsumer()
	// if err != nil {
	// 	log.Fatalf("Failed to init NSQ Consumer: %v", err)
	// }

	chatBox := usecase.NewChatBox()
	router.Init(chatBox)
}
