package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Store struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
}

func NewStore() *Store {
	return &Store{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *Store) Create(title string) (Todo, error) {
	title, err := validateTitle(title)
	if err != nil {
		return Todo{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		ID:        s.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	s.todos[todo.ID] = todo
	s.nextID++
	return todo, nil
}

func (s *Store) Get(id int) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, fmt.Errorf("get todo %d: %w", id, ErrNotFound)
	}
	return todo, nil
}

func (s *Store) GetAll() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []Todo{}
	for _, todo := range s.todos {
		result = append(result, todo)
	}
	return result
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.todos[id]
	if !ok {
		return fmt.Errorf("delete todo %d: %w", id, ErrNotFound)
	}
	delete(s.todos, id)
	return nil
}

func (s *Store) Update(id int, title string, done bool) (Todo, error) {
	title, err := validateTitle(title)
	if err != nil {
		return Todo{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, fmt.Errorf("update todo %d: %w", id, ErrNotFound)
	}

	todo.Title = title
	todo.Done = done
	s.todos[id] = todo

	return todo, nil
}

const maxTitleLength = 200

func validateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) > maxTitleLength {
		return "", ErrTitleTooLong
	}
	return title, nil
}
