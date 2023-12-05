package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nickyrolly/ws-chat-demo/internal/router/handler"
	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

func Init() {
	r := mux.NewRouter()
	r.HandleFunc("/check", handler.CheckServices).Methods("GET")
	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleChat(w, r)
	})
	r.HandleFunc("/chat/history", func(w http.ResponseWriter, r *http.Request) {
		handler.GetChatHistory(w, r)
	}).Methods("GET")

	http.Handle("/", r)

	port := "8080"
	fmt.Printf("Listening on port %s...\n", port)
	go usecase.NsqHealthCheck()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
