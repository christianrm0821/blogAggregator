-- name: CreateFeedFollow :one
Insert into feed_follows(id,created_at, updated_at, user_id, feed_id)
Values(
    $1,
    $2,
    $3,
    $4,
    $5
)    