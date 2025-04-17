-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
where id = $1;