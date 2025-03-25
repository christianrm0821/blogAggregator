-- +goose Up
create table users(
    id integer primary key,
    create_at timestamp not null,
    updated_at timestamp not null,
    name text unique not null
);



-- +goose Down
drop table users;