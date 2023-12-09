package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	var (
		user     LoginUser
		response map[string]interface{}
	)

	// Try to decode the JSON request to a LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a JWT
	token, err := usecase.LoginUser(r.Context(), user.Username, user.Password)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	response = map[string]interface{}{
		"data": map[string]interface{}{
			"token": token,
		},
		"server":  "chat",
		"status":  "OK",
		"success": 1,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func HandleRegister(w http.ResponseWriter, r *http.Request) {

	var (
		user     LoginUser
		response map[string]interface{}
	)

	// Try to decode the JSON request to a LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = usecase.RegisterUser(r.Context(), user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response = map[string]interface{}{
		"data": map[string]interface{}{
			"success": true,
		},
		"server":  "chat",
		"status":  "OK",
		"success": 1,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
