package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func HandleChat(w http.ResponseWriter, r *http.Request, chatBox *usecase.ChatBox, groupChatBox *usecase.GroupChatBox) {
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
		groupChatBox.AddClient(groupID, conn)
		defer groupChatBox.RemoveClient(groupID, conn)

	} else {
		// Create user chatbox id
		chatboxUsers := []string{userIDStr, recipientIDStr}
		if userID > recipientID {
			chatboxUsers = []string{recipientIDStr, userIDStr}
		}
		userChatboxID = strings.Join(chatboxUsers, "-")

		chatBox.AddClient(userChatboxID, conn)
		defer chatBox.RemoveClient(userChatboxID, conn)
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
				groupChatBox.Broadcast(groupID, conn, message)

				// Exercise 3.2.1
				// Please complete this block to publish group chat message data via NSQ

				groupChatBox.PublishGroupSaveChatHistory(repository.GroupChatHistoryData{
					GroupID:      groupID,
					SenderUserID: userID,
					Message:      message,
					ReplyTime:    time.Now(),
				})

			} else {
				chatBox.Broadcast(userChatboxID, conn, message)
				chatBox.Broadcast(userChatboxID, conn, message)

				// Sort user id for A and B to be from lowest to highest
				userIDA := userID
				userIDB := recipientID

				if userIDA > userIDB {
					userIDA = recipientID
					userIDB = userID
				}

				// Exercise 3.1.1
				// Please complete this block to publish personal message data via NSQ
			}
		}
	}

}

func GetChatHistory(w http.ResponseWriter, r *http.Request, chatBox *usecase.ChatBox, groupChatBox *usecase.GroupChatBox) {
	var (
		//Uncomment this section (3.3.2)
		userID          int
		recipientID     int
		groupID         int
		chatHistoryData []map[string]interface{}
		response        map[string]interface{}
		err             error
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

	if groupID > 0 {
		// Exercise 3.4.1
		// Please complete this block to call get group chat history functionality

		chatHistoryData, err = groupChatBox.GetGroupChatHistory(r.Context(), repository.GroupChatHistoryData{
			GroupID: groupID,
		})
	} else {
		// Sort user id for A and B to be from lowest to highest
		userIDA := userID
		userIDB := recipientID

		if userIDA > userIDB {
			userIDA = recipientID
			userIDB = userID
		}

		// Exercise 3.3.2
		// Please complete this block to call chat history functionality
	}

	if err != nil {
		log.Printf("Error get chat history: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response = map[string]interface{}{
		"data":    chatHistoryData,
		"server":  "chat",
		"status":  "OK",
		"success": 1,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
