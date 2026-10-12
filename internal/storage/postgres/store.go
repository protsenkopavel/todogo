package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"todogo/internal/storage/postgres/db"
	"todogo/internal/todo"
)

type Store struct {
	q *db.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{q: db.New(pool)}
}

func toDomain(row db.Todo) todo.Todo {
	return todo.Todo{
		ID:        int(row.ID),
		Title:     row.Title,
		Done:      row.Done,
		CreatedAt: row.CreatedAt,
	}
}

func (s *Store) Create(ctx context.Context, title string) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	row, err := s.q.CreateTodo(ctx, title)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	return toDomain(row), nil
}

func (s *Store) Get(ctx context.Context, id int) (todo.Todo, error) {
	row, err := s.q.GetTodo(ctx, int32(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return todo.Todo{}, fmt.Errorf("get todo %d: %w", id, todo.ErrNotFound)
	}
	if err != nil {
		return todo.Todo{}, fmt.Errorf("get todo %d: %w", id, err)
	}

	return toDomain(row), nil
}

func (s *Store) GetAll(ctx context.Context) ([]todo.Todo, error) {
	rows, err := s.q.ListTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	items := make([]todo.Todo, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomain(row))
	}

	return items, nil
}

func (s *Store) Update(ctx context.Context, id int, title string, done bool) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	row, err := s.q.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:    int32(id),
		Title: title,
		Done:  done,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return todo.Todo{}, fmt.Errorf("update todo %d: %w", id, todo.ErrNotFound)
	}
	if err != nil {
		return todo.Todo{}, fmt.Errorf("update todo %d: %w", id, err)
	}

	return toDomain(row), nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	affected, err := s.q.DeleteTodo(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}

	if affected == 0 {
		return fmt.Errorf("delete todo %d: %w", id, todo.ErrNotFound)
	}

	return nil
}
