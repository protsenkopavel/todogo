package main

import (
	"log"
	"net/http"
)

func newRouter(h *TodoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", h.GetAll)
	mux.HandleFunc("POST /todos", h.Create)
	mux.HandleFunc("GET /todos/{id}", h.GetByID)
	mux.HandleFunc("PUT /todos/{id}", h.Update)
	mux.HandleFunc("DELETE /todos/{id}", h.Delete)
	return mux
}

func main() {
	store := NewStore()
	handler := NewTodoHandler(store)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", newRouter(handler)))
}
