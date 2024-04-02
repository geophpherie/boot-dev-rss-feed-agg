-- name: GetPostsByUser :many
SELECT
  *
FROM
  Posts
WHERE
  feed_id IN (
    SELECT
      id
    FROM
      Feeds
    WHERE
      user_id = $1
  )
LIMIT
  $2;

-- name: CreatePost :one
INSERT INTO
  Posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  *;
