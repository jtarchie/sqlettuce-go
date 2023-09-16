// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package writers

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
	if q.addFloatStmt, err = db.PrepareContext(ctx, addFloat); err != nil {
		return nil, fmt.Errorf("error preparing query AddFloat: %w", err)
	}
	if q.addIntStmt, err = db.PrepareContext(ctx, addInt); err != nil {
		return nil, fmt.Errorf("error preparing query AddInt: %w", err)
	}
	if q.appendStmt, err = db.PrepareContext(ctx, append); err != nil {
		return nil, fmt.Errorf("error preparing query Append: %w", err)
	}
	if q.flushAllStmt, err = db.PrepareContext(ctx, flushAll); err != nil {
		return nil, fmt.Errorf("error preparing query FlushAll: %w", err)
	}
	if q.setStmt, err = db.PrepareContext(ctx, set); err != nil {
		return nil, fmt.Errorf("error preparing query Set: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addFloatStmt != nil {
		if cerr := q.addFloatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addFloatStmt: %w", cerr)
		}
	}
	if q.addIntStmt != nil {
		if cerr := q.addIntStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addIntStmt: %w", cerr)
		}
	}
	if q.appendStmt != nil {
		if cerr := q.appendStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing appendStmt: %w", cerr)
		}
	}
	if q.flushAllStmt != nil {
		if cerr := q.flushAllStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing flushAllStmt: %w", cerr)
		}
	}
	if q.setStmt != nil {
		if cerr := q.setStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setStmt: %w", cerr)
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
	addFloatStmt *sql.Stmt
	addIntStmt   *sql.Stmt
	appendStmt   *sql.Stmt
	flushAllStmt *sql.Stmt
	setStmt      *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:           tx,
		tx:           tx,
		addFloatStmt: q.addFloatStmt,
		addIntStmt:   q.addIntStmt,
		appendStmt:   q.appendStmt,
		flushAllStmt: q.flushAllStmt,
		setStmt:      q.setStmt,
	}
}
