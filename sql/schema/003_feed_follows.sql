-- +goose Up
create table feed_follows(
    id primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id UUID not null,
    feed_id UUID not null,
    constraint fk_ff_users
        foreign key (user_id)
        references users(id) On delete cascade,
    constraint fk_ff_feed
        foreign key (feed_id)
        references feeds(id) on delete cascade,
    unique(user_id,feed_id)
);


-- +goose Down
drop table feed_follows;