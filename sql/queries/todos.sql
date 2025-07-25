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

-- name: UpdateToDo :one
UPDATE todos
SET
    updated_at = NOW(),
    date = @date,
    title = @title,
    description = @description
WHERE id = @todo_id
RETURNING *;

-- name: DeleteTodos :exec
DELETE FROM todos;

-- name: DeleteTodoByID :exec
DELETE FROM todos WHERE id = $1;

-- name: GetTodoByID :one
SELECT * FROM todos WHERE id = $1;

-- name: GetTodosByUserID :many
SELECT * FROM todos WHERE user_id = $1;
