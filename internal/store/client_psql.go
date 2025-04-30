package store

import (
	"database/sql"
	"fmt"
	"net/url"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	storegen "github.com/vovanwin/meetingsBot/internal/store/gen"
)

//go:generate options-gen -out-filename=client_psql_options.gen.go -from-struct=PSQLOptions
type PSQLOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
	debug    bool
}

func NewPSQLClient(opts PSQLOptions) (*storegen.Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	db, err := NewPgxDB(NewPgxOptions(opts.address, opts.username, opts.password, opts.database))
	if err != nil {
		return nil, fmt.Errorf("init db driver: %v", err)
	}

	var clientOpts []storegen.Option
	clientOpts = append(clientOpts, storegen.Driver(entsql.OpenDB(dialect.Postgres, db)))
	if opts.debug {
		clientOpts = append(clientOpts, storegen.Debug())
	}

	return storegen.NewClient(clientOpts...), nil
}

//go:generate options-gen -out-filename=client_psql_pgx_options.gen.go -from-struct=PgxOptions
type PgxOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
}

func NewPgxDB(opts PgxOptions) (*sql.DB, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("pgx options: %v", err)
	}

	connStr := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(opts.username, opts.password),
		Host:   opts.address,
		Path:   opts.database,
	}).String()
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v", err)
	}

	return db, nil
}
