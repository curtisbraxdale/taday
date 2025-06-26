-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    priority BOOLEAN NOT NULL DEFAULT FALSE,
    recur_d BOOLEAN NOT NULL DEFAULT FALSE,
    recur_w BOOLEAN NOT NULL DEFAULT FALSE,
    recur_m BOOLEAN NOT NULL DEFAULT FALSE,
    recur_y BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE events;
