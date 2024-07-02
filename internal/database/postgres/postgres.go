// Package postgres provides structures and methods to operate databases
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tasktracker/internal/logger"

	"tasktracker/configs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type Database struct {
	DB     *sqlx.DB
	Logger *logger.Logger
}

const (
	driverName = "pgx"

	// migrationsPath is used during migrations
	migrationsPath = "../../internal/database/migrations"
)

// OpenDB returns an object representing an abstraction over a DB conn and an errors if any
func OpenDB(logger *logger.Logger) (*Database, error) {

	db, err := sqlx.Connect(driverName, configs.GenPostgresURI())
	if err != nil {
		return nil, fmt.Errorf("opening and verifying a new conn: %w", err)
	}

	if err := migrateUp(db.DB); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("applying migrations: %w; closing db: %w", err, closeErr)
		}

		return nil, fmt.Errorf("applying migrations: %w", err)
	}

	

	return &Database{DB: db, Logger: logger}, nil
}

// Close downs the migrations and closes the DB instance and returns an error if any
func (db *Database) Close() error {
	var (
		resultingErr error
	)

	if err := migrateDown(db.DB.DB); err != nil {
		resultingErr = fmt.Errorf("applying migration: %w", err)
	}

	dbClosingErr := db.DB.Close()
	if dbClosingErr == nil {
		return resultingErr
	}

	if resultingErr != nil {
		resultingErr = fmt.Errorf("closing db: %w; %w", dbClosingErr, resultingErr)
	}
	resultingErr = fmt.Errorf("closing db: %w", dbClosingErr)

	return resultingErr
}

// Begin begins a transaction and returns the object to work with it and an error if any
func (db *Database) Begin(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}

// Exec execurtes the query given and return the result and an error if any
func (db *Database) Exec(query string) (sql.Result, error) {
	return db.DB.Exec(query)
}

// migrateUp applies the up migrations referencing to the global const migrationsPath
func migrateUp(db *sql.DB) error {
	return goose.Up(db, migrationsPath)
}

// migrateUp applies the up migrations referencing to the global const migrationsPath
func migrateDown(db *sql.DB) error {
	return goose.Down(db, migrationsPath)
}
