-- name: GetNextFeedToFetch :one
select * from feed_follows
where user_id = $1
order by updated_at nulls first
limit 1;