// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package poststore

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createStmt, err = db.PrepareContext(ctx, create); err != nil {
		return nil, fmt.Errorf("error preparing query Create: %w", err)
	}
	if q.getAllStmt, err = db.PrepareContext(ctx, getAll); err != nil {
		return nil, fmt.Errorf("error preparing query GetAll: %w", err)
	}
	if q.getByIdsStmt, err = db.PrepareContext(ctx, getByIds); err != nil {
		return nil, fmt.Errorf("error preparing query GetByIds: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createStmt != nil {
		if cerr := q.createStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createStmt: %w", cerr)
		}
	}
	if q.getAllStmt != nil {
		if cerr := q.getAllStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllStmt: %w", cerr)
		}
	}
	if q.getByIdsStmt != nil {
		if cerr := q.getByIdsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getByIdsStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db           DBTX
	tx           *sql.Tx
	createStmt   *sql.Stmt
	getAllStmt   *sql.Stmt
	getByIdsStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:           tx,
		tx:           tx,
		createStmt:   q.createStmt,
		getAllStmt:   q.getAllStmt,
		getByIdsStmt: q.getByIdsStmt,
	}
}
