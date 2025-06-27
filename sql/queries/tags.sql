-- name: CreateTag :one
INSERT INTO tags (id, user_id, name, color)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: CreateEventTag :one
INSERT INTO event_tags (event_id, tag_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetTagsByUserID :many
SELECT * FROM tags WHERE user_id = $1;

-- name: GetEventTagsByEventID :many
SELECT * FROM event_tags WHERE event_id = $1;

-- name: GetTagsByEventID :many
SELECT tags.id, tags.name, tags.color
FROM event_tags
JOIN tags ON tags.id = event_tags.tag_id
WHERE event_tags.event_id = @event_id;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: DeleteEventTag :exec
DELETE FROM event_tags WHERE event_id = $1 AND tag_id = $2;
