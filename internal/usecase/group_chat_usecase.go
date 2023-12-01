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

type groupConnmap map[int]userConnMap

type GroupChatBox struct {
	clients groupConnmap
	mu      sync.Mutex
}

func NewGroupChatBox() *GroupChatBox {
	return &GroupChatBox{
		clients: make(groupConnmap),
	}
}

func (cb *GroupChatBox) AddClient(groupID, userID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[groupID]; !ok {
		cb.clients[groupID] = make(userConnMap)
	}
	if _, ok := cb.clients[groupID][userID]; !ok {
		cb.clients[groupID][userID] = []*websocket.Conn{}
	}
	cb.clients[groupID][userID] = append(cb.clients[groupID][userID], conn)
	log.Printf("Add client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *GroupChatBox) RemoveClient(groupID, userID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[groupID]; ok {
		if _, ok := cb.clients[groupID][userID]; ok {
			// Find conn index
			idx := cb.findUserConn(groupID, userID, conn)
			// Remove conn from slice if conn found
			if idx != -1 {
				cb.clients[groupID][userID] = append(cb.clients[groupID][userID][:idx], cb.clients[groupID][userID][idx+1:]...)
			}
		}
	}
	log.Printf("Remove client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *GroupChatBox) Broadcast(userID, groupID int, curConn *websocket.Conn, message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	log.Println("groupID :", groupID)

	// Send to destination group pools except for current user connection
	for _, userConns := range cb.clients[groupID] {
		for _, conn := range userConns {
			if conn == curConn {
				continue
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error broadcasting message to user :", err)
			}
		}
	}

	log.Printf("Broadcast clients : %+v\n", cb.clients)
}

func (cb *GroupChatBox) findUserConn(groupID, userID int, conn *websocket.Conn) int {
	if _, ok := cb.clients[groupID]; ok {
		for i, c := range cb.clients[groupID][userID] {
			if c == conn {
				return i
			}
		}
	}
	return -1
}

func (cb *GroupChatBox) PublishGroupSaveChatHistory(params repository.GroupChatHistoryData) error {
	// Publish a message
	messageBody, err := json.Marshal(params)
	if err != nil {
		log.Println("Error Marshal:", err)
		return err
	}

	err = chat_nsq.NSQProducer.Publish("save-group-chat-history-topic", messageBody)
	if err != nil {
		log.Println("Error Publish NSQ:", err)
		return err
	}

	return nil
}

func (cb *GroupChatBox) GetGroupChatHistory(ctx context.Context, params repository.GroupChatHistoryData) ([]map[string]interface{}, error) {
	var (
		groupChatHistoryData = []map[string]interface{}{}
		err                  error
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	rows, err := postgre.DBChat.QueryContext(ctx, postgre.QuerySelectGroupChatHistory, params.GroupID)
	if err != nil {
		log.Println("[GetGroupChatHistory] Error QueryContext: ", err)
		return groupChatHistoryData, err
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
			log.Println("[GetGroupChatHistory] Error Scan: ", err)
			return groupChatHistoryData, err
		}

		data = map[string]interface{}{
			"sender_user_id": senderUserID,
			"message":        message,
			"reply_time":     replyTime,
		}

		groupChatHistoryData = append(groupChatHistoryData, data)
	}

	return groupChatHistoryData, nil
}
