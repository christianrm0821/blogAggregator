-- name: CreatPost :one
Insert into posts(id,created_at,updated_at,title,url,description,published_at,feed_id)
Values(
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6,
    $7,
    $8
)
Returning *;