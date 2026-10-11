package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	todo := h.store.Create(req.Title)

	writeJSON(w, http.StatusCreated, todo)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	return strconv.Atoi(idStr)
}

func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

type updateTodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	todo, ok := h.store.Update(id, req.Title, req.Done)
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if !h.store.Delete(id) {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
