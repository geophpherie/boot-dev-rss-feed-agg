-- name: DeleteFeedFollow :exec
DELETE FROM feedfollows
WHERE
  id = $1;

-- name: CreateFeedFollow :one
INSERT INTO
  FeedFollows (id, created_at, updated_at, feed_id, user_id)
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  *;

-- name: GetUserFeedFollows :many
SELECT
  *
FROM
  feedfollows
WHERE
  user_id = $1;

INSERT INTO
  FEEDS (id)
VALUES
  (1);
