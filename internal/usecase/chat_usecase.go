package usecase

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/chat_nsq"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/postgre"
)

type userConnMap map[int][]*websocket.Conn

type ChatBox struct {
	clients userConnMap
	mu      sync.Mutex
}

func NewChatBox() *ChatBox {
	return &ChatBox{
		clients: make(userConnMap),
	}
}

func (cb *ChatBox) AddClient(userID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[userID]; !ok {
		cb.clients[userID] = []*websocket.Conn{}
	}
	cb.clients[userID] = append(cb.clients[userID], conn)
	log.Printf("Add client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *ChatBox) RemoveClient(userID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[userID]; ok {
		// Find conn index
		idx := cb.findConn(userID, conn)
		// Remove conn from slice if conn found
		if idx != -1 {
			cb.clients[userID] = append(cb.clients[userID][:idx], cb.clients[userID][idx+1:]...)
		}
	}
	log.Printf("Remove client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *ChatBox) Broadcast(userID, destID int, curConn *websocket.Conn, message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	log.Println("destID :", destID)

	// Send to destination user pools
	for _, conn := range cb.clients[destID] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error broadcasting message to user :", err)
		}
	}

	// Send to current user pools except for current user connection
	for _, conn := range cb.clients[userID] {
		if conn == curConn {
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error broadcasting message to user :", err)
		}
	}

	log.Printf("Broadcast clients : %+v\n", cb.clients)
}

func (cb *ChatBox) findConn(userID int, conn *websocket.Conn) int {
	for i, c := range cb.clients[userID] {
		if c == conn {
			return i
		}
	}
	return -1
}

func (cb *ChatBox) PublishSaveChatHistory(params repository.ChatHistoryData) error {
	// Publish a message
	messageBody, err := json.Marshal(params)
	if err != nil {
		log.Println("Error Marshal:", err)
		return err
	}

	err = chat_nsq.NSQProducer.Publish("save-chat-history-topic", messageBody)
	if err != nil {
		log.Println("Error Publish NSQ:", err)
		return err
	}

	return nil
}

func (cb *ChatBox) GetChatHistory(ctx context.Context, params repository.ChatHistoryData) ([]map[string]interface{}, error) {
	var (
		chatHistoryData = []map[string]interface{}{}
		err             error
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()
	log.Printf("%+v", params)

	rows, err := postgre.DBChat.QueryContext(ctx, postgre.QuerySelectChatHistory, params.UserIDA, params.UserIDB)
	if err != nil {
		log.Println("[GetChatHistory] Error QueryContext: ", err)
		return chatHistoryData, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			data         map[string]interface{}
			senderUserID int
			message      string
			replyTime    time.Time
		)

		err := rows.Scan(&senderUserID, &message, &replyTime)
		if err != nil {
			log.Println("[GetChatHistory] Error Scan: ", err)
			return chatHistoryData, err
		}

		data = map[string]interface{}{
			"sender_user_id": senderUserID,
			"message":        message,
			"reply_time":     replyTime,
		}

		chatHistoryData = append(chatHistoryData, data)
	}

	return chatHistoryData, nil
}
