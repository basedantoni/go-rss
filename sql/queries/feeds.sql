-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds ORDER BY last_fetched_at DESC LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
  SET (last_fetched_at, updated_at) = ($1, $1)
  WHERE id = $2;