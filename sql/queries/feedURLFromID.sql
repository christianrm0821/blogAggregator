-- name: GetFeedURLFromFeedID :one
select url from feeds
where id = $1;