package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"todogo/internal/todo"
)

type Store struct {
	mu     sync.RWMutex
	todos  map[int]todo.Todo
	nextID int
}

func New() *Store {
	return &Store{
		todos:  make(map[int]todo.Todo),
		nextID: 1,
	}
}

func (s *Store) Create(ctx context.Context, title string) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item := todo.Todo{
		ID:        s.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	s.todos[item.ID] = item
	s.nextID++
	return item, nil
}

func (s *Store) Get(ctx context.Context, id int) (todo.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.todos[id]
	if !ok {
		return todo.Todo{}, fmt.Errorf("get todo %d: %w", id, todo.ErrNotFound)
	}
	return item, nil
}

func (s *Store) GetAll(ctx context.Context) ([]todo.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]todo.Todo, 0, len(s.todos))
	for _, item := range s.todos {
		result = append(result, item)
	}
	return result, nil
}

func (s *Store) Update(ctx context.Context, id int, title string, done bool) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.todos[id]
	if !ok {
		return todo.Todo{}, fmt.Errorf("update todo %d: %w", id, todo.ErrNotFound)
	}

	item.Title = title
	item.Done = done
	s.todos[id] = item

	return item, nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("delete todo %d: %w", id, todo.ErrNotFound)
	}
	delete(s.todos, id)
	return nil
}
