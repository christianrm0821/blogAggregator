-- name: GetPostForUser :many
select * from posts
inner join feed_follows
on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1
order by published_at desc
limit 10;
