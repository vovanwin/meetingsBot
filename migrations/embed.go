package db

import "embed"

//go:embed *.sql
var EmbedMigrations embed.FS //
