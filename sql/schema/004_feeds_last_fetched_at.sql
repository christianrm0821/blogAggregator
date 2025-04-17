-- +goose Up
Alter table feeds
ADD last_fetched_at timestamp;


-- +goose Down
Alter table feeds
drop last_fetched_at;