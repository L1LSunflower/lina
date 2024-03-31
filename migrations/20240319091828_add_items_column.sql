-- +goose Up
-- +goose StatementBegin
create table items(
    id uuid primary key not null,
    url varchar(250) not null,
    name varchar(250) null,
    article varchar(20) null,
    price float not null,
    price_on_sale float not null,
    currency varchar(10) not null,
    colors text[] null,
    sizes text[] null,
    image_links text[] null,
    hash varchar(255) not null,
    status varchar(20) not null
);
create index hash_index on items (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists hash_index;
drop table if exists items;
-- +goose StatementEnd
