package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nickyrolly/ws-chat-demo/internal/router/handler"
	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

func Init(usercb *usecase.ChatBox, groupcb *usecase.GroupChatBox) {
	r := mux.NewRouter()
	r.HandleFunc("/check", handler.CheckServices).Methods("GET")
	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleChat(w, r, usercb, groupcb)
	})

	//Exercise 3.3.1
	// Please complete this block to add new route for Get Historical Chat From Database
	//--

	http.Handle("/", r)

	port := "8080"
	fmt.Printf("Listening on port %s...\n", port)
	go usecase.NsqHealthCheck()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
