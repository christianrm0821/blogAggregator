-- name: GetFeedByUrl :one
select * from feeds
where url = $1;