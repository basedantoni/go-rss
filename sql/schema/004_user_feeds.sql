-- +goose Up
CREATE TABLE users_feeds (
    id UUID NOT NULL,
    user_id UUID REFERENCES users(id),
    feed_id UUID REFERENCES feeds(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT users_feeds_pk PRIMARY KEY(user_id, feed_id)
);

-- +goose Down
DROP TABLE IF EXISTS users_feeds;