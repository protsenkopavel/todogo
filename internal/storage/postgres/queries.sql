-- name: CreateTodo :one
INSERT INTO todos (title) VALUES ($1)
    RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY id;

-- name: UpdateTodo :one
UPDATE todos SET title = $1, done = $2
WHERE id = $3
    RETURNING *;

-- name: DeleteTodo :execrows
DELETE FROM todos
WHERE id = $1;