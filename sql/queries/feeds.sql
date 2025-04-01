-- name: CreateFeed :one
Insert into feeds(id,created_at, updated_at, name, url,user_id )
Values(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
Returning *;