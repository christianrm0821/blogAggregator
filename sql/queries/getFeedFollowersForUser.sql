-- name: GetFeedFollowsForUser :many
select feed_follows.* ,feeds.name as feed_name, users.name as user_name
from feed_follows
inner join users on users.id = feed_follows.user_id
inner join feeds on feeds.id = feed_follows.feed_id
where users.id = $1;