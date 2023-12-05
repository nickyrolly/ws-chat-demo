package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var WSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for this example. You might want to restrict this in production.
		return true
	},
}

func HandleChat(w http.ResponseWriter, r *http.Request) {

}

func GetChatHistory(w http.ResponseWriter, r *http.Request) {

}
