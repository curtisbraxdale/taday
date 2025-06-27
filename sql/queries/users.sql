-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, hashed_password, phone_number)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = NOW(),
    username = @username,
    email = @email,
    hashed_password = @hashed_password,
    phone_number = @phone_number
WHERE id = @userID
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
