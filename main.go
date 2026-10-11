package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewStore()
	handler := NewTodoHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", handler.GetAll)
	mux.HandleFunc("POST /todos", handler.Create)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
