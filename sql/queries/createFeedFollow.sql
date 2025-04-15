-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    Insert into feed_follows(id,created_at, updated_at, user_id, feed_id)
    Values(
    $1,
    $2,
    $3,
    $4,
    $5
)   
    Returning *
)  
select inserted_feed_follow.*,
feeds.name as feed_name,
users.name as user_name
from inserted_feed_follow
inner join users on users.id = inserted_feed_follow.user_id
inner join feeds on feeds.id  = inserted_feed_follow.feed_id;
