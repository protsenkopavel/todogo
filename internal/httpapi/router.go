package httpapi

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", h.getAll)
	mux.HandleFunc("POST /todos", h.create)
	mux.HandleFunc("GET /todos/{id}", h.getByID)
	mux.HandleFunc("PUT /todos/{id}", h.update)
	mux.HandleFunc("DELETE /todos/{id}", h.delete)
	return mux
}
