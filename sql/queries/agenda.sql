-- name: GetUserEventsToday :many
SELECT title, description
FROM events
WHERE user_id = @user_id
  AND start_date < @date_plus_1_day
  AND end_date >= @date
ORDER BY start_date ASC;

-- name: GetUserEventsWeek :many
SELECT title, description
FROM events
WHERE user_id = @user_id
  AND start_date < @date_plus_7_days
  AND end_date >= @date
ORDER BY start_date ASC;

-- name: GetAllUsers :many
SELECT id, username, phone_number FROM users;
