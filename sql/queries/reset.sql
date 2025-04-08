-- name: ResetUsers :exec
delete from users;

-- name: ResetFeed :exec
delete from feeds;

-- name: ResetFeedFollows :exec
delete from feed_follows;