package usecase

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
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
