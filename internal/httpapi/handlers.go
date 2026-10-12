package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"todogo/internal/todo"
)

type TodoStore interface {
	Create(ctx context.Context, title string) (todo.Todo, error)
	Get(ctx context.Context, id int) (todo.Todo, error)
	GetAll(ctx context.Context) ([]todo.Todo, error)
	Update(ctx context.Context, id int, title string, done bool) (todo.Todo, error)
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	store TodoStore
}

func NewHandler(store TodoStore) *Handler {
	return &Handler{store: store}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type updateTodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.GetAll(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.store.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	item, err := h.store.Create(r.Context(), req.Title)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	item, err := h.store.Update(r.Context(), id, req.Title, req.Done)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
