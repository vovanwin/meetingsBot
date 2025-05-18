package repository

import (
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/Tdep"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dbsqlc"
	"github.com/vovanwin/meetingsBot/pkg/clients/sqlite"
	"go.uber.org/zap"
	"log"

	"go.uber.org/fx"
)

type Options struct {
	fx.In
	Db     *sqlite.SQLiteClient `validate:"required"`
	Logger *Tdep.TelegramLogger
}

type Repo struct {
	Db     *dbsqlc.Queries
	logger *zap.Logger
}

func New(opts Options) (*Repo, error) {
	if opts.Db == nil {
		log.Fatalf("")
	}

	return &Repo{Db: dbsqlc.New(opts.Db), logger: opts.Logger.Lg}, nil
}
