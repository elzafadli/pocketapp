package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"seedapp/config"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=mocks/sqlx.go -package=mocks -source=sqlx.go
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
	SqlxDBIsObjectMigrationAlreadyExists(err error) bool
}

type SqlxDB struct {
	*sqlx.DB
	Conf *config.Config `inject:"config"`
}

func (s *SqlxDB) Startup() error {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&application_name=%s",
		url.QueryEscape(s.Conf.Database.User),
		url.QueryEscape(s.Conf.Database.Password),
		url.QueryEscape(s.Conf.Database.Host),
		url.QueryEscape(s.Conf.Database.Port),
		url.QueryEscape(s.Conf.Database.Name),
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
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" WITH SCHEMA public")

	_, err = db.ExecContext(context.Background(), `DO $$
	BEGIN
		CREATE SCHEMA IF NOT EXISTS main;
		CREATE TABLE IF NOT EXISTS main.migrations (
			id SERIAL PRIMARY KEY,
			schema VARCHAR(255) NOT NULL,
			version VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			error TEXT,
			started_at TIMESTAMP NOT NULL,
			finished_at TIMESTAMP NULL DEFAULT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL DEFAULT NULL
		);
	END $$;`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(context.Background(), `DO $$
	BEGIN
		CREATE SCHEMA IF NOT EXISTS main;
		CREATE TABLE IF NOT EXISTS main.seeds (
			id SERIAL PRIMARY KEY,
			schema VARCHAR(255) NOT NULL,
			version VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			entity_processed TEXT[] NULL DEFAULT NULL,
			error TEXT,
			started_at TIMESTAMP NOT NULL,
			finished_at TIMESTAMP NULL DEFAULT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL DEFAULT NULL
		);
	END $$;`)
	if err != nil {
		return err
	}

	s.DB = db
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
	var pgErr *pgx.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *SqlxDB) SqlxDBIsObjectMigrationAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	// Check for pgx.PgError with specific error codes
	var pgErr *pgx.PgError
	if errors.As(err, &pgErr) {
		// 42P07: duplicate_table (table already exists)
		// 42710: duplicate_object (constraint, index, etc. already exists)
		if pgErr.Code == "42P07" || pgErr.Code == "42710" {
			return true
		}
	}

	// Also check error message for common "already exists" patterns
	// This is a fallback for cases where the error might be wrapped
	errMsg := strings.ToLower(err.Error())

	// Check for specific "already exists" patterns
	alreadyExistsKeywords := []string{
		"already exists",
		"duplicate_table",
		"duplicate_object",
		"relation.*already exists",
		"constraint.*already exists",
	}

	for _, keyword := range alreadyExistsKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	// Check for "duplicate" combined with object types
	if strings.Contains(errMsg, "duplicate") {
		objectTypes := []string{"table", "constraint", "index", "relation", "object"}
		for _, objType := range objectTypes {
			if strings.Contains(errMsg, objType) {
				return true
			}
		}
	}

	return false
}
