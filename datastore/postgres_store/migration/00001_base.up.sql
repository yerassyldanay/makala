create schema if not exists feed;
create table feed.posts (
    id bigserial primary key,
    title varchar not null,
    author varchar not null,
    link varchar,
    submakala varchar not null ,
    content varchar,
    score float8 not null default 0,
    promoted bool not null default false,
    nsfw bool not null default false
);
