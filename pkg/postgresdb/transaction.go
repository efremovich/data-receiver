package postgresdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/jmoiron/sqlx"
)

var _ Transaction = &transactionImpl{}

type Transaction interface {
	Commit() error
	Rollback() error
	RollbackIfNotCommitted() error

	QueryExecutor
}

type transactionImpl struct {
	*sqlx.Tx
}

func (tx *transactionImpl) ExecOne(query string, args ...interface{}) (sql.Result, error) {
	res, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		alogger.WarnFromCtx(context.Background(), "драйвер не поддерживает RowsAffected(): " + err.Error(), nil, nil, false)
		return res, nil
	}

	if affected == 1 {
		return res, nil
	}

	if affected > 1 {
		return res, ErrorMultipleRows
	}

	return nil, sql.ErrNoRows
}

func (tx *transactionImpl) QueryAndScan(dest interface{}, query string, args ...interface{}) error {
	rows, err := tx.Queryx(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(dest); err != nil {
			return err
		}
	}

	return nil
}

func (tx *transactionImpl) Ping() error {
	_, err := tx.Tx.Exec("SELECT 1")

	return err
}

func (tx *transactionImpl) RollbackIfNotCommitted() error {
	err := tx.Ping()
	if errors.Is(err, sql.ErrTxDone) {
		return nil
	}

	err = tx.Rollback()
	if err != nil {
		return fmt.Errorf("ошибка роллбека транзакции: %s", err.Error())
	}

	return nil
}

type TransactionMock struct{}

func (TransactionMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (TransactionMock) ExecOne(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (TransactionMock) Get(dest interface{}, query string, args ...interface{}) error {
	return nil
}
func (TransactionMock) Select(dest interface{}, query string, args ...interface{}) error {
	return nil
}
func (TransactionMock) QueryAndScan(dest interface{}, query string, args ...interface{}) error {
	return nil
}
func (TransactionMock) Ping() error {
	return nil
}
func (TransactionMock) Commit() error {
	return nil
}
func (TransactionMock) Rollback() error {
	return nil
}
func (TransactionMock) RollbackIfNotCommitted() error {
	return nil
}
