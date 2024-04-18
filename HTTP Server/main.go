package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

var users = make(map[string]User)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	users[newUser.ID] = newUser
	w.WriteHeader(http.StatusCreated)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	_, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	users[id] = updatedUser
	w.WriteHeader(http.StatusOK)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	_, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusOK)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/user/create", CreateUserHandler)
	http.HandleFunc("/user/get", GetUserHandler)
	http.HandleFunc("/user/update", UpdateUserHandler)
	http.HandleFunc("/user/delete", DeleteUserHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
