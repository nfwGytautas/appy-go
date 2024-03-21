package driver

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TransactionFn func(tx *Tx) error

type Tx struct {
	internal pgx.Tx
}

type ExecResult struct {
	res pgconn.CommandTag
}

type RowResult struct {
	isRead bool
	row    pgx.Row
}

type RowsResult struct {
	rows pgx.Rows
}

type Scannable interface {
	Scan(dest ...any) error
	Close()
	Next() bool
	Err() error
}

func StartTransaction() (*Tx, error) {
	tx, err := gDatabaseConnection.Begin(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Tx{
		internal: tx,
	}, nil
}

func RunTransaction(fn TransactionFn) error {
	tx, err := StartTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)

	if err == nil {
		return tx.Commit()
	}

	return err
}

func (tx *Tx) Exec(query string, args ...any) (ExecResult, error) {
	res, err := tx.internal.Exec(context.TODO(), query, args...)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{
		res: res,
	}, nil
}

func (tx *Tx) QueryRow(query string, args ...any) RowResult {
	res := tx.internal.QueryRow(context.TODO(), query, args...)
	return RowResult{
		isRead: false,
		row:    res,
	}
}

func (tx *Tx) Query(query string, args ...any) (RowsResult, error) {
	res, err := tx.internal.Query(context.TODO(), query, args...)
	return RowsResult{
		rows: res,
	}, err
}

func (tx *Tx) Commit() error {
	return tx.internal.Commit(context.TODO())
}

func (tx *Tx) Rollback() error {
	return tx.internal.Rollback(context.TODO())
}

// Commit the transaction or rollback on failure
func (tx *Tx) CommitOrRollback() {
	if err := tx.Commit(); err != nil {
		tx.Rollback()
	}
}

// Compatibility only
func (rr RowResult) Err() error {
	return nil
}

func (rr RowResult) Scan(dest ...any) error {
	rr.isRead = true
	return rr.row.Scan(dest...)
}

func (rr RowResult) Close() {
	// Do nothing
}

func (rr RowResult) Next() bool {
	return rr.isRead
}

func (rr RowsResult) Scan(dest ...any) error {
	return rr.rows.Scan(dest...)
}

func (rr RowsResult) Err() error {
	return rr.rows.Err()
}

func (rr RowsResult) Close() {
	rr.rows.Close()
}

func (rr RowsResult) Next() bool {
	return rr.rows.Next()
}

func (er ExecResult) RowsAffected() int64 {
	return er.res.RowsAffected()
}
