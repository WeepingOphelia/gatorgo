-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT
  *
FROM
  feeds
WHERE
  url = $1;

-- name: GetFeeds :many
SELECT
  *
FROM
  feeds;

-- name: GetNextFeedToFetch :one
SELECT
  *
FROM
  feeds
ORDER BY
  last_fetched_at ASC NULLS FIRST LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;

-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING *
)

SELECT
  new_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM
  new_feed_follow
INNER JOIN feeds ON new_feed_follow.feed_id = feeds.id
INNER JOIN users ON new_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT
  *
FROM
  feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE
  users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM
  feed_follows
WHERE
  feed_id = $1
  AND user_id = $2;

-- name: CreatePost :exec
INSERT INTO posts (
  id,
  created_at,
  updated_at,
  title,
  url,
  description,
  published_at,
  feed_id
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
);

-- name: GetPostsForUser :many
SELECT
  *
FROM
  posts
INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE user_id = $1
ORDER BY published_at ASC
LIMIT $2;