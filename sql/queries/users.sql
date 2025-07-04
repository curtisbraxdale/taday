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

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByStripeID :one
SELECT * FROM users WHERE stripe_customer_id = $1;

-- name: UpdateStripeCustomerID :exec
UPDATE users
SET stripe_customer_id = $2
WHERE id = $1;

-- name: GetStripeID :one
SELECT stripe_customer_id FROM users WHERE id = $1;

-- name: GetEmail :one
SELECT email FROM users WHERE id = $1;
