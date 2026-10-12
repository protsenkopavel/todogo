package postgres

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"todogo/internal/todo"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	flag.Parse()
	if testing.Short() {
		return m.Run()
	}

	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:17",
		tcpostgres.WithDatabase("todo"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.WithInitScripts("../../../schema.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Printf("start postgres container: %v", err)
		return 1
	}
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("connection string: %v", err)
		return 1
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Printf("create pool: %v", err)
		return 1
	}
	defer testPool.Close()

	return m.Run()
}

func newTestStore(t *testing.T) *Store {
	t.Helper()

	if testing.Short() {
		t.Skip("postgres tests are skipped in short mode")
	}

	_, err := testPool.Exec(context.Background(), "TRUNCATE todos RESTART IDENTITY")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}

	return New(testPool)
}

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	created, err := store.Create(ctx, "купить молоко")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("ID = %d, want 1", created.ID)
	}
	if created.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be set by database")
	}

	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "купить молоко" || got.Done {
		t.Errorf("unexpected todo: %+v", got)
	}
}

func TestCreate_ValidatesTitle(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	_, err := store.Create(ctx, "   ")
	if !errors.Is(err, todo.ErrEmptyTitle) {
		t.Fatalf("err = %v, want ErrEmptyTitle", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	_, err := store.Get(ctx, 999)
	if !errors.Is(err, todo.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestGetAll_OrderedByID(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	titles := []string{"первая", "вторая", "третья"}
	for _, title := range titles {
		if _, err := store.Create(ctx, title); err != nil {
			t.Fatalf("create %q: %v", title, err)
		}
	}

	items, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(items) != len(titles) {
		t.Fatalf("got %d items, want %d", len(items), len(titles))
	}
	for i, item := range items {
		if item.Title != titles[i] {
			t.Errorf("items[%d].Title = %q, want %q", i, item.Title, titles[i])
		}
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	created, err := store.Create(ctx, "купить молоко")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	updated, err := store.Update(ctx, created.ID, "купить кефир", true)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "купить кефир" || !updated.Done {
		t.Errorf("unexpected todo after update: %+v", updated)
	}

	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "купить кефир" || !got.Done {
		t.Errorf("update not persisted: %+v", got)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	_, err := store.Update(ctx, 999, "что-то", false)
	if !errors.Is(err, todo.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)

	created, err := store.Create(ctx, "удалить меня")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	err = store.Delete(ctx, created.ID)
	if !errors.Is(err, todo.ErrNotFound) {
		t.Fatalf("second delete: err = %v, want ErrNotFound", err)
	}
}
