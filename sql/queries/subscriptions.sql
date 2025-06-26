-- name: CreateSubscription :one
INSERT INTO subscriptions (id, user_id, stripe_customer_id, stripe_subscription_id, plan, status, current_period_start, current_period_end, cancel_at_period_end, canceled_at, trial_start, trial_end, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
)
RETURNING *;

-- name: DeleteSubscriptions :exec
DELETE FROM subscriptions;

-- name: GetSubscriptionByUserID :one
SELECT * FROM subscriptions WHERE user_id = $1;
