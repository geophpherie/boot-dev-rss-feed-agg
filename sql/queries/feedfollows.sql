-- name: CreateFeedFollow :one
INSERT INTO FeedFollows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteFeedFollow :exec
DELETE FROM FeedFollows
WHERE id = $1;
-- name: GetUserFeedFollows :many
SELECT *
FROM FeedFollows
WHERE user_id = $1;