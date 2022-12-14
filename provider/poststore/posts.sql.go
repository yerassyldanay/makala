// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: posts.sql

package poststore

import (
	"context"

	"github.com/lib/pq"
)

const create = `-- name: Create :one
insert into feed.posts (title, author, link, submakala, content, score, promoted, nsfw)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning id, title, author, link, submakala, content, score, promoted, nsfw
`

type CreateParams struct {
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Link      *string `json:"link"`
	Submakala string  `json:"submakala"`
	Content   *string `json:"content"`
	Score     float64 `json:"score"`
	Promoted  bool    `json:"promoted"`
	Nsfw      bool    `json:"nsfw"`
}

// stores post in database (without any further operation on the post)
func (q *Queries) Create(ctx context.Context, arg CreateParams) (FeedPost, error) {
	row := q.queryRow(ctx, q.createStmt, create,
		arg.Title,
		arg.Author,
		arg.Link,
		arg.Submakala,
		arg.Content,
		arg.Score,
		arg.Promoted,
		arg.Nsfw,
	)
	var i FeedPost
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Link,
		&i.Submakala,
		&i.Content,
		&i.Score,
		&i.Promoted,
		&i.Nsfw,
	)
	return i, err
}

const getAll = `-- name: GetAll :many
select id, title, author, link, submakala, content, score, promoted, nsfw from feed.posts where promoted = false order by id offset $1 limit $2
`

type GetAllParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// fetches all posts ordered by their id
func (q *Queries) GetAll(ctx context.Context, arg GetAllParams) ([]FeedPost, error) {
	rows, err := q.query(ctx, q.getAllStmt, getAll, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedPost
	for rows.Next() {
		var i FeedPost
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Link,
			&i.Submakala,
			&i.Content,
			&i.Score,
			&i.Promoted,
			&i.Nsfw,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getByIds = `-- name: GetByIds :many
select id, title, author, link, submakala, content, score, promoted, nsfw from feed.posts where id = ANY($1::bigint[])
`

// fetches posts by ids
func (q *Queries) GetByIds(ctx context.Context, dollar_1 []int64) ([]FeedPost, error) {
	rows, err := q.query(ctx, q.getByIdsStmt, getByIds, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedPost
	for rows.Next() {
		var i FeedPost
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Link,
			&i.Submakala,
			&i.Content,
			&i.Score,
			&i.Promoted,
			&i.Nsfw,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
