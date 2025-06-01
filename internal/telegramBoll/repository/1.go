package repository

import (
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dbsqlc"
	"github.com/vovanwin/meetingsBot/pkg/clients/postgres"
	"go.uber.org/zap"
	"log"

	"go.uber.org/fx"
)

type Options struct {
	fx.In
	Db *postgres.Postgres `validate:"required"`
}

type Repo struct {
	Db     *dbsqlc.Queries
	logger *zap.Logger
}

func New(opts Options) (*Repo, error) {
	if opts.Db == nil {
		log.Fatalf("")
	}

	return &Repo{Db: dbsqlc.New(opts.Db.Pool)}, nil
}
