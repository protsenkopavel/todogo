package memory

import (
	"context"
	"errors"
	"sync"
	"testing"

	"todogo/internal/todo"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	store := New()

	item, err := store.Create(ctx, "купить молоко")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != 1 {
		t.Errorf("ID = %d, want 1", item.ID)
	}
	if item.Title != "купить молоко" {
		t.Errorf("Title = %q, want %q", item.Title, "купить молоко")
	}
	if item.Done {
		t.Errorf("new todo should not be done")
	}
}

func TestGet_NotFound(t *testing.T) {
	ctx := context.Background()
	store := New()

	_, err := store.Get(ctx, 999)
	if !errors.Is(err, todo.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestCreate_Concurrent(t *testing.T) {
	ctx := context.Background()
	store := New()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Create(ctx, "task")
		}()
	}
	wg.Wait()

	items, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1000 {
		t.Errorf("created %d todos, want 1000", len(items))
	}
}
