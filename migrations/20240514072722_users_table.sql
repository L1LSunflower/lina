-- +goose Up
-- +goose StatementBegin
create table users (
    id varchar(100) unique not null,
    username varchar(100) not null,
    first_name varchar(100) null,
    blocked boolean not null,
    created_at timestamp not null default now(),
    update_at timestamp not null default  now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
