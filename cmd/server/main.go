package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"todogo/internal/httpapi"
	"todogo/internal/storage/postgres"
)

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

	store := postgres.New(pool)
	handler := httpapi.NewHandler(store)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpapi.NewRouter(handler)))
}
