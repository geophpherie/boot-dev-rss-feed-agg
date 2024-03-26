-- name: CreateFeed :one
INSERT INTO Feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetFeeds :many
SELECT *
FROM Feeds;
-- name: GetNextFeedsToFetch :many
SELECT *
FROM Feeds
ORDER BY last_fetched_at DESC
LIMIT $1;
-- name: MarkFeedFetched :execrows
UPDATE Feeds
SET created_at = $1,
    last_fetched_at = $1
WHERE id = $2;