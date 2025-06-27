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

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: DeleteEventTag :exec
DELETE FROM event_tags WHERE event_id = $1 AND tag_id = $2;
