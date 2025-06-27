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

-- name: DeleteEvents :exec
DELETE FROM events;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: GetEventsByUserID :many
SELECT * FROM events WHERE user_id = $1;

-- name: GetFilteredEvents :many
SELECT *
FROM events
WHERE events.user_id = @user_id
  AND (@start_date::timestamptz IS NULL OR start_date >= @start_date)
  AND (@end_date::timestamptz IS NULL OR start_date < @end_date)
  AND (
      @tag::text IS NULL OR
      EXISTS (
          SELECT 1 FROM event_tags et
          JOIN tags t ON t.id = et.tag_id
          WHERE et.event_id = events.id AND t.name = @tag
      )
  );
