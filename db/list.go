package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettuce/db/drivers/sqlite/writers"
	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound = errors.New("record was not found")
	ErrNotArray = errors.New("not an array")
)

func (c *Client) ListInsert(
	ctx context.Context,
	name string,
	offset int64,
	pivot, value string,
) (int64, bool, error) {
	row := c.db.QueryRowContext(ctx, `
	-- name: ListIndex :one
	 UPDATE keys
	 SET value = json_array_insert(
	 		value,
			?4,
			?2,
			?3
		)
	 WHERE name = ?1
	 RETURNING json_array_length(value);
	`, name, value, offset, pivot)
	if row.Err() != nil {
		return 0, false, fmt.Errorf("could not execute ListInsert: %w", row.Err())
	}

	var newOffset int64

	err := row.Scan(&newOffset)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("could not scan ListInsert: %w", err)
	}

	return newOffset, true, nil
}

func (c *Client) ListRange(ctx context.Context, name string, start, end int64) ([]string, error) {
	rows, err := c.db.QueryContext(ctx, `
	-- name: ListRange :many
		SELECT json_each.value
		FROM keys,
			json_each(keys.value)
		WHERE keys.name = ?1
		AND json_each.key >= IIF(?2 >=0, ?2, json_array_length(keys.value) + ?2)
		AND json_each.key <= IIF(?3 >=0, ?3, json_array_length(keys.value) + ?3);
	`, name, start, end)
	if err != nil {
		return nil, fmt.Errorf("could not execute ListRange: %w", err)
	}
	defer rows.Close()

	var values []string

	for rows.Next() {
		var value string

		_ = rows.Scan(&value)
		values = append(values, value)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not execute ListRange: %w", rows.Err())
	}

	return values, nil
}

func (c *Client) ListLength(ctx context.Context, name string) (int64, error) {
	length, err := c.readers.ListLength(ctx, name)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("could not execute ListSet: %w", err)
	}

	if value, ok := length.(int64); ok {
		return value, nil
	}

	return 0, ErrNotFound
}

func (c *Client) ListRightPush(ctx context.Context, name string, values ...string) (int64, error) {
	transaction, err := c.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("could not start ListRightPush: %w", err)
	}
	//nolint:errcheck
	defer transaction.Rollback()

	var result writers.ListRightPushRow

	queries := c.writers.WithTx(transaction)

	for _, value := range values {
		result, err = queries.ListRightPush(ctx, &writers.ListRightPushParams{
			Name:  name,
			Value: value,
		})

		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		//nolint:errorlint
		if err, ok := err.(sqlite3.Error); ok && err.Error() == "malformed JSON" {
			return 0, ErrNotArray
		}

		if err != nil {
			return 0, fmt.Errorf("could not execute ListRightPush: %w", err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return 0, fmt.Errorf("could not ListRightPush: %w", err)
	}

	return result.Column2, nil
}

func (c *Client) ListRightPushUpsert(ctx context.Context, name string, values ...string) (int64, bool, error) {
	transaction, err := c.db.Begin()
	if err != nil {
		return 0, false, fmt.Errorf("could not start ListRightPushUpsert: %w", err)
	}
	//nolint:errcheck
	defer transaction.Rollback()

	var result writers.ListRightPushUpsertRow

	queries := c.writers.WithTx(transaction)

	for _, value := range values {
		result, err = queries.ListRightPushUpsert(ctx, &writers.ListRightPushUpsertParams{
			Name:  name,
			Value: value,
		})

		//nolint:errorlint
		if err, ok := err.(sqlite3.Error); ok && err.Error() == "malformed JSON" {
			return 0, true, nil
		}

		if err != nil {
			return 0, true, fmt.Errorf("could not execute ListRightPushUpsert: %w", err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return 0, false, fmt.Errorf("could not ListRightPushUpsert: %w", err)
	}

	return result.Column2, result.Column1, nil
}

func (c *Client) ListSet(ctx context.Context, name string, index int64, value string) (bool, error) {
	valid, err := c.writers.ListSet(ctx, &writers.ListSetParams{
		Name:  name,
		Index: index,
		Value: value,
	})

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("could not execute ListSet: %w", err)
	}

	if value, ok := valid.(int64); ok {
		return value == 1, nil
	}

	return false, nil
}
