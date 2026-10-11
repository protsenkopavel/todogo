package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	store := NewPostgresStore(pool)
	handler := NewTodoHandler(store)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", newRouter(handler)))
}