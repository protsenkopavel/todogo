package main

import "time"

type Store struct {
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
	todo, ok := s.todos[id]
	return todo, ok
}

func (s *Store) GetAll() []Todo {
	result := []Todo{}
	for _, todo := range s.todos {
		result = append(result, todo)
	}
	return result
}

func (s *Store) Delete(id int) bool {
	_, ok := s.todos[id]
	if !ok {
		return false
	}
	delete(s.todos, id)
	return true
}
