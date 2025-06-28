-- name: CreateEvent :one
INSERT INTO events (id, user_id, created_at, updated_at, start_date, end_date, title, description, priority, recur_d, recur_w, recur_m, recur_y)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET
    updated_at = NOW(),
    start_date = @start_date,
    end_date = @end_date,
    title = @title,
    description = @description,
    priority = @priority,
    recur_d = @recur_d,
    recur_w = @recur_w,
    recur_m = @recur_m,
    recur_y = @recur_y
WHERE id = @event_id
RETURNING *;

-- name: DeleteEvents :exec
DELETE FROM events;

-- name: DeleteEventByID :exec
DELETE FROM events WHERE id = $1;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: GetEventsByUserID :many
SELECT * FROM events WHERE user_id = $1;
