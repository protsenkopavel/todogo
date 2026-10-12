package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"todogo/internal/todo"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Create(ctx context.Context, title string) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	var item todo.Todo
	err = s.pool.QueryRow(ctx,
		`INSERT INTO todos (title) VALUES ($1)
         RETURNING id, title, done, created_at`,
		title,
	).Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	return item, nil
}

func (s *Store) Get(ctx context.Context, id int) (todo.Todo, error) {
	var item todo.Todo
	err := s.pool.QueryRow(ctx,
		`SELECT id, title, done, created_at FROM todos WHERE id = $1`,
		id,
	).Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return todo.Todo{}, fmt.Errorf("get todo %d: %w", id, todo.ErrNotFound)
	}
	if err != nil {
		return todo.Todo{}, fmt.Errorf("get todo %d: %w", id, err)
	}

	return item, nil
}

func (s *Store) GetAll(ctx context.Context) ([]todo.Todo, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, done, created_at FROM todos ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	items := []todo.Todo{}
	for rows.Next() {
		var item todo.Todo
		if err := rows.Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	return items, nil
}

func (s *Store) Update(ctx context.Context, id int, title string, done bool) (todo.Todo, error) {
	title, err := todo.ValidateTitle(title)
	if err != nil {
		return todo.Todo{}, err
	}

	var item todo.Todo
	err = s.pool.QueryRow(ctx,
		`UPDATE todos SET title = $1, done = $2 WHERE id = $3
         RETURNING id, title, done, created_at`,
		title, done, id,
	).Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return todo.Todo{}, fmt.Errorf("update todo %d: %w", id, todo.ErrNotFound)
	}
	if err != nil {
		return todo.Todo{}, fmt.Errorf("update todo %d: %w", id, err)
	}

	return item, nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete todo %d: %w", id, todo.ErrNotFound)
	}

	return nil
}
