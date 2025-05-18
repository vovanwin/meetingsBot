package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/sqlite" // чисто-Go драйвер
)

type SQLiteClient struct {
	*sql.DB
}

// ConnectSQLite открывает SQLite с хорошими настройками.
// dsn — например: "file:data.db" или "file::memory:?cache=shared".
func ConnectSQLite(ctx context.Context, dsn string) (*SQLiteClient, error) {
	// Добавляем оптимальные PRAGMA через DSN
	fullDSN := fmt.Sprintf("%s?_foreign_keys=on&journal_mode=WAL&busy_timeout=5000", dsn)

	db, err := sql.Open("sqlite", fullDSN)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	// Пул соединений
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Ping с таймаутом
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return &SQLiteClient{db}, nil
}
