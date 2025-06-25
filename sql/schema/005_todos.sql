-- +goose Up
CREATE TABLE todos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    date TIMESTAMP,
    title TEXT NOT NULL,
    description TEXT,
);

-- +goose Down
DROP TABLE todos;
