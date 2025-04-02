-- name: ListFeeds :many
select feeds.name, feeds.url, users.name 
from feeds
inner join users
on feeds.user_id = users.id;
