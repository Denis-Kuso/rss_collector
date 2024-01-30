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
-- name: RemoveFeedFollow :one
DELETE FROM feed_follows 
WHERE ID_FF=$1
RETURNING *;

-- name: GetAllFeedFollows :many
SELECT * FROM feed_follows
WHERE user_id=$1;
