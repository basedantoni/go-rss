-- +goose Up
CREATE TABLE feeds (
    id UUID,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    PRIMARY KEY(id),
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;