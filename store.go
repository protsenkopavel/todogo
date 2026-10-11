package main

import (
	"sync"
	"time"
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

func (s *Store) Create(title string) Todo {
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
	return todo
}

func (s *Store) Get(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	todo, ok := s.todos[id]
	return todo, ok
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

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.todos[id]
	if !ok {
		return false
	}
	delete(s.todos, id)
	return true
}

func (s *Store) Update(id int, title string, done bool) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, false
	}

	todo.Title = title
	todo.Done = done
	s.todos[id] = todo

	return todo, true
}
