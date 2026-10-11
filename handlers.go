package main

import (
	"encoding/json"
	"net/http"
)

type TodoHandler struct {
	store *Store
}

func NewTodoHandler(store *Store) *TodoHandler {
	return &TodoHandler{store: store}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todos := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	todo := h.store.Create(req.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
