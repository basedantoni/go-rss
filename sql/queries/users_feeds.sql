-- name: CreateUserFeed :one
INSERT INTO users_feeds (id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteUserFeed :exec
DELETE FROM users_feeds WHERE id = $1;

-- name: GetUserFeedFollows :many
SELECT * FROM users_feeds WHERE user_id = $1;