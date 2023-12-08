package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
)

type groupConnmap map[int][]*websocket.Conn

type GroupChatBox struct {
	clients groupConnmap
	mu      sync.Mutex
}

func NewGroupChatBox() *GroupChatBox {
	return &GroupChatBox{
		clients: make(groupConnmap),
	}
}

func (cb *GroupChatBox) AddClient(groupID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[groupID]; !ok {
		cb.clients[groupID] = []*websocket.Conn{}
	}
	cb.clients[groupID] = append(cb.clients[groupID], conn)
	log.Printf("Add client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *GroupChatBox) RemoveClient(groupID int, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[groupID]; ok {
		// Find conn index
		idx := cb.findConn(groupID, conn)
		// Remove conn from slice if conn found
		if idx != -1 {
			cb.clients[groupID] = append(cb.clients[groupID][:idx], cb.clients[groupID][idx+1:]...)
		}
	}
	log.Printf("Remove client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *GroupChatBox) Broadcast(groupID int, curConn *websocket.Conn, message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	log.Println("chatboxID :", groupID)

	// Send to chatbox user pools except for current connection
	for _, conn := range cb.clients[groupID] {
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

func (cb *GroupChatBox) findConn(groupID int, conn *websocket.Conn) int {
	for i, c := range cb.clients[groupID] {
		if c == conn {
			return i
		}
	}
	return -1
}

func (cb *GroupChatBox) PublishGroupSaveChatHistory(params repository.GroupChatHistoryData) error {
	// messageBody, err := json.Marshal(params)
	// if err != nil {
	// 	log.Println("Error Marshal:", err)
	// 	return err
	// }

	// Exercise 3.2.2
	// Please complete this block to publish message to NSQ topic

	return nil
}

func (cb *GroupChatBox) GetGroupChatHistory(ctx context.Context, params repository.GroupChatHistoryData) ([]map[string]interface{}, error) {
	var (
		groupChatHistoryData = []map[string]interface{}{}
		//err                  error
	)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	// Exercise 3.4.2
	// Please complete this block to get group chat history data from database

	// Uncomment code below after get chat history from DB has implemented (Exercise 3.4.2)
	// for rows.Next() {
	// 	var (
	// 		data         map[string]interface{}
	// 		senderUserID int
	// 		message      string
	// 		replyTime    time.Time
	// 	)

	// 	err := rows.Scan(&senderUserID, &message, &replyTime)
	// 	if err != nil {
	// 		log.Println("[GetGroupChatHistory] Error Scan: ", err)
	// 		return groupChatHistoryData, err
	// 	}

	// 	data = map[string]interface{}{
	// 		"sender_user_id": senderUserID,
	// 		"message":        message,
	// 		"reply_time":     replyTime,
	// 	}

	// 	groupChatHistoryData = append(groupChatHistoryData, data)
	// }

	return groupChatHistoryData, nil
}
