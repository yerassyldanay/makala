-- name: Create :one
-- stores post in database (without any further operation on the post)
insert into feed.posts (title, author, link, submakala, content, score, promoted, nsfw)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetAll :many
-- fetches all posts ordered by their id
select * from feed.posts where promoted = false order by id offset $1 limit $2;

-- name: GetByIds :many
-- fetches posts by ids
select * from feed.posts where id = ANY($1::bigint[]);
