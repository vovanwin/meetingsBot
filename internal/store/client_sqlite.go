// Package store provides database clients for the bot.
package store

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
	storegen "github.com/vovanwin/meetingsBot/internal/store/gen"
)

// SQLiteOptions holds configuration for an SQLite client.
//
//go:generate options-gen -out-filename=client_sqlite_options.gen.go -from-struct=SQLiteOptions
type SQLiteOptions struct {
	path    string `option:"mandatory" validate:"required"` // Filepath or DSN (e.g. ":memory:" or "./db.sqlite3?_foreign_keys=1")
	isDebug bool   // Enable debug logging
}

// NewSQLiteClient creates an ent client backed by SQLite.
func NewSQLiteClient(opts SQLiteOptions) (*storegen.Client, error) {
	// Validate required options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate sqlite options: %w", err)
	}

	// Open the SQLite database
	db, err := sql.Open("sqlite3", opts.path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// Wrap in ent SQL driver
	driver := entsql.OpenDB(dialect.SQLite, db)

	// Prepare client options
	var clientOpts []storegen.Option
	clientOpts = append(clientOpts, storegen.Driver(driver))
	if opts.isDebug {
		clientOpts = append(clientOpts, storegen.Debug())
	}

	// Create ent client
	client := storegen.NewClient(clientOpts...)
	return client, nil
}
