package main

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	todo, err := store.Create(ctx, "купить молоко")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID != 1 {
		t.Errorf("ID = %d, want 1", todo.ID)
	}
	if todo.Title != "купить молоко" {
		t.Errorf("Title = %q, want %q", todo.Title, "купить молоко")
	}
	if todo.Done {
		t.Errorf("new todo should not be done")
	}
}

func TestGet_NotFound(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	_, err := store.Get(ctx, 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"обычное название", "купить молоко", "купить молоко", nil},
		{"пробелы по краям обрезаются", "  молоко  ", "молоко", nil},
		{"пустая строка", "", "", ErrEmptyTitle},
		{"только пробелы", "     ", "", ErrEmptyTitle},
		{"200 кириллических символов", strings.Repeat("я", 200), strings.Repeat("я", 200), nil},
		{"201 символ", strings.Repeat("я", 201), "", ErrTitleTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTitle(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreate_Concurrent(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Create(ctx, "task")
		}()
	}
	wg.Wait()

	todos, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 1000 {
		t.Errorf("created %d todos, want 1000", len(todos))
	}
}
