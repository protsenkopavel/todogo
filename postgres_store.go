package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

var _ TodoStore = (*PostgresStore)(nil)

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Create(ctx context.Context, title string) (Todo, error) {
	title, err := validateTitle(title)
	if err != nil {
		return Todo{}, err
	}

	var todo Todo
	err = s.pool.QueryRow(ctx,
		`INSERT INTO todos (title) VALUES ($1)
         RETURNING id, title, done, created_at`,
		title,
	).Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt)
	if err != nil {
		return Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	return todo, nil
}

func (s *PostgresStore) Get(ctx context.Context, id int) (Todo, error) {
	var todo Todo
	err := s.pool.QueryRow(ctx,
		`SELECT id, title, done, created_at FROM todos WHERE id = $1`,
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Todo{}, fmt.Errorf("get todo %d: %w", id, ErrNotFound)
	}
	if err != nil {
		return Todo{}, fmt.Errorf("get todo %d: %w", id, err)
	}

	return todo, nil
}

func (s *PostgresStore) GetAll(ctx context.Context) ([]Todo, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, done, created_at FROM todos ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	return todos, nil
}

func (s *PostgresStore) Update(ctx context.Context, id int, title string, done bool) (Todo, error) {
	title, err := validateTitle(title)
	if err != nil {
		return Todo{}, err
	}

	var todo Todo
	err = s.pool.QueryRow(ctx,
		`UPDATE todos SET title = $1, done = $2 WHERE id = $3
         RETURNING id, title, done, created_at`,
		title, done, id,
	).Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Todo{}, fmt.Errorf("update todo %d: %w", id, ErrNotFound)
	}
	if err != nil {
		return Todo{}, fmt.Errorf("update todo %d: %w", id, err)
	}

	return todo, nil
}

func (s *PostgresStore) Delete(ctx context.Context, id int) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete todo %d: %w", id, ErrNotFound)
	}

	return nil
}
