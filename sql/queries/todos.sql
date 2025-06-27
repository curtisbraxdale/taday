-- name: CreateTodo :one
INSERT INTO todos (id, user_id, created_at, updated_at, date, title, description)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: DeleteTodos :exec
DELETE FROM todos;

-- name: GetTodoByID :one
SELECT * FROM todos WHERE id = $1;

-- name: GetTodosByUserID :many
SELECT * FROM todos WHERE user_id = $1;
