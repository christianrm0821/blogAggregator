-- name: GetFeedNameFromID :one
select name from feeds
where id = $1;