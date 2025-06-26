-- +goose Up
CREATE TABLE tags (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

CREATE TABLE event_tags (
    event_id UUID REFERENCES events (id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);

-- +goose Down
DROP TABLE event_tags;

DROP TABLE tags;
