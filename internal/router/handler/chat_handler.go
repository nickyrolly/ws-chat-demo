package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		// Exercise 2.1
		// Please complete this block to register user in group chatbox
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
			} else {
				chatBox.Broadcast(userChatboxID, conn, message)
			}
		}
	}

}

func GetChatHistory(w http.ResponseWriter, r *http.Request, chatBox *usecase.ChatBox, groupChatBox *usecase.GroupChatBox) {
	var (
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

		chatHistoryData, err = chatBox.GetChatHistory(r.Context(), repository.ChatHistoryData{
			UserIDA: userIDA,
			UserIDB: userIDB,
		})
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
