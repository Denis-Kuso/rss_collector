-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows WHERE user_id = $1;
--

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
--

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
--

-- name: GetFeedFollow :one
SELECT * from feed_follows where id = $1;
