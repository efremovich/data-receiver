package postgresdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/efremovich/data-receiver/pkg/alogger"
)

var ErrorMultipleRows = fmt.Errorf("pg_error_multiple_rows")

type DBConnection struct {
	writeConn PostgresConnection
	readConn  PostgresConnection
}

func (p *DBConnection) GetWriteConnection() PostgresConnection {
	return p.writeConn
}

func (p *DBConnection) GetReadConnection() PostgresConnection {
	if p.readConn == nil {
		return p.writeConn
	}
	return p.readConn
}

func New(ctx context.Context, masterConnString string, slaveConnString string) (*DBConnection, error) {
	writeConn, err := newConnect(ctx, masterConnString)
	if err != nil {
		return nil, err
	}
	res := &DBConnection{writeConn: &postgresConnectionImpl{writeConn}}

	alogger.InfoFromCtx(ctx, "инициировано подключение к master DB %s", masterConnString)

	if slaveConnString != "" {
		readConn, err := newConnect(ctx, slaveConnString)
		if err != nil {
			return nil, err
		}

		alogger.InfoFromCtx(ctx, "инициировано подключение к slave DB %s", slaveConnString)
		res.readConn = &postgresConnectionImpl{readConn}
	}

	return &DBConnection{writeConn: &postgresConnectionImpl{writeConn}}, nil
}

func newConnect(ctx context.Context, connString string) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	connConfig.DescriptionCacheCapacity = 0
	connConfig.StatementCacheCapacity = 0

	db, err := sqlx.ConnectContext(ctx, "pgx", stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = db.Close()
	}()

	return db, nil
}

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecOne(query string, args ...interface{}) (sql.Result, error) // возвращает ошибку, если затронут не ровно 1 ряд
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	QueryAndScan(dest interface{}, query string, args ...interface{}) error
	Ping() error
}

type PostgresConnection interface {
	QueryExecutor
	BeginTX(ctx context.Context) (Transaction, error)
}

type postgresConnectionImpl struct {
	*sqlx.DB
}

func (r *postgresConnectionImpl) ExecOne(query string, args ...interface{}) (sql.Result, error) {
	res, err := r.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		alogger.WarnFromCtx(context.Background(), "драйвер не поддерживает RowsAffected(): "+err.Error(), nil, nil, false)
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

func (r *postgresConnectionImpl) QueryAndScan(dest interface{}, query string, args ...interface{}) error {
	rows, err := r.Queryx(query, args...)
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

func (r *postgresConnectionImpl) BeginTX(ctx context.Context) (Transaction, error) {
	tx, err := r.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	newTx := transactionImpl{tx}

	return &newTx, nil
}
