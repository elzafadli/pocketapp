package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"monitorapp/config"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Sqlx interface {
	Ping() error
	GetDB() *sqlx.DB

	// Prepare statement

	// SelectContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// GetContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Executor

	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// QueryRowxContext queries the database and returns an *sqlx.Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// Unit of work

	// Tx starts a transaction and execute function
	Tx(ctx context.Context, fn func(ctx context.Context, err chan error)) error
	// TxWithValue starts a transaction and execute function
	TxWithValue(ctx context.Context, fn func(ctx context.Context, err chan error, res chan interface{})) (interface{}, error)

	// IsErrorDuplicate check if error is duplicate error
	IsErrorDuplicate(err error) bool
}

type SqlxDB struct {
	*sqlx.DB
	Conf *config.Config `inject:"config"`
}

func (s *SqlxDB) Startup() error {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&application_name=%s",
		url.QueryEscape(s.Conf.Database.Username),
		url.QueryEscape(s.Conf.Database.Password),
		url.QueryEscape(s.Conf.Database.Host),
		url.QueryEscape(s.Conf.Database.Port),
		url.QueryEscape(s.Conf.Database.DBName),
		url.QueryEscape(s.Conf.Database.SSLMode),
		url.QueryEscape(s.Conf.GetAppName()),
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(s.Conf.Database.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(s.Conf.Database.ConnMaxLifetime) * time.Hour)
	db.SetMaxOpenConns(s.Conf.Database.MaxOpenConn)

	s.DB = db

	// Run migrations automatically on startup
	m, err := migrate.New(
		"file://internal/adapter/repository/migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	return nil
}

func (s *SqlxDB) Shutdown() error {
	db := s.DB

	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *SqlxDB) Ping() error {
	db := s.DB

	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// GetDB get db
func (r *SqlxDB) GetDB() *sqlx.DB {
	return r.DB
}

func (r *SqlxDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.DB.SelectContext(ctx, dest, query, args...)
}

func (r *SqlxDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.DB.GetContext(ctx, dest, query, args...)
}

func (r *SqlxDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, ok := ctx.Value(TX_KEY).(*sqlx.Tx)
	if !ok {
		return r.DB.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (r *SqlxDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	tx, ok := ctx.Value(TX_KEY).(*sqlx.Tx)
	if !ok {
		return r.DB.QueryRowxContext(ctx, query, args...)
	}
	return tx.QueryRowxContext(ctx, query, args...)
}

func (r *SqlxDB) Tx(ctx context.Context, fun func(ctx context.Context, err chan error)) error {
	e := make(chan error, 1)

	if _, ok := ctx.Value(TX_KEY).(*sqlx.Tx); !ok {
		tx, err := r.DB.Beginx()
		if err != nil {
			return err
		}
		ctx = context.WithValue(ctx, TX_KEY, tx)

		go fun(ctx, e)

		if err := <-e; err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				err = fmt.Errorf("%q: %w", errTx, err)
			}
			return err
		}
		return tx.Commit()
	}

	go fun(ctx, e)

	return <-e
}

func (r *SqlxDB) TxWithValue(ctx context.Context, fun func(ctx context.Context, err chan error, res chan interface{})) (interface{}, error) {
	e, result := make(chan error, 1), make(chan interface{}, 1)

	if _, ok := ctx.Value(TX_KEY).(*sqlx.Tx); !ok {
		tx, err := r.DB.Beginx()
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, TX_KEY, tx)

		go fun(ctx, e, result)

		if err := <-e; err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				err = fmt.Errorf("%q: %w", errTx, err)
			}
			return nil, err
		}
		return <-result, tx.Commit()
	}

	go fun(ctx, e, result)

	return <-result, <-e
}

func (r *SqlxDB) IsErrorDuplicate(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
