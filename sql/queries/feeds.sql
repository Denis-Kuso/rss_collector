-- name: CreateFeed :many
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6)
RETURNING *;
INSERT INTO feed_follows(
    user_id,
    feed_id,
    ID_FF)
VALUES (
    $6,
    $1,
    $7)
RETURNING *;
-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
INSERT INTO feed_follows(
    user_id,
    feed_id,
    ID_FF)
VALUES (
    $1,
    $2,
    $3)
RETURNING *;
