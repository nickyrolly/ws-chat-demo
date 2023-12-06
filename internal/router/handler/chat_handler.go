package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

var WSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for this example. You might want to restrict this in production.
		return true
	},
}

func HandleChat(w http.ResponseWriter, r *http.Request, chatBox *usecase.ChatBox) {
	var (
		userID      int
		recipientID int
		groupID     int
		err         error
	)

	userIDStr := r.URL.Query().Get("user_id")
	userID, err = strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user_id : \"%s\"\n", userIDStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipientIDStr := r.URL.Query().Get("recipient_id")
	if recipientIDStr != "" {
		recipientID, err = strconv.Atoi(recipientIDStr)
		if err != nil {
			log.Printf("Invalid recipient_id : \"%s\"\n", recipientIDStr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr != "" {
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			log.Printf("Invalid group_id : \"%s\"\n", groupIDStr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if recipientIDStr == "" && groupIDStr == "" {
		log.Println("Empty resipient_id or group_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var userChatboxID string

	if groupID > 0 {
		// Exercise 2
		// Please complete this block
	} else {
		// Excercise 1
		// please complete this block
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.TextMessage {
			message := string(p)
			fmt.Println("Received message:", message)
			if groupID > 0 {
				// Exercise 2
				// Please complete this block
			} else {
				chatBox.Broadcast(userChatboxID, conn, message)

				// Sort user id for A and B to be from lowest to highest
				userIDA := userID
				userIDB := recipientID

				if userIDA > userIDB {
					userIDA = recipientID
					userIDB = userID
				}

				// Save message to database via NSQ
				chatBox.PublishSaveChatHistory(repository.ChatHistoryData{
					UserIDA:      userIDA,
					UserIDB:      userIDB,
					SenderUserID: userID,
					Message:      message,
					ReplyTime:    time.Now(),
				})
			}
		}
	}

}

func GetChatHistory(w http.ResponseWriter, r *http.Request, chatBox *usecase.ChatBox, groupChatBox *usecase.GroupChatBox) {
	// chapter 3
}
