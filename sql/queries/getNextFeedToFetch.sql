-- name: GetNextFeedToFetch :one
select * from feeds
where user_id = $1
order by last_fetched_at nulls first
limit 1;