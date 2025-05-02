package repository

import (
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vovanwin/meetingsBot/internal/store/gen"
)

type Options struct {
	fx.In
	Db *gen.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Db     *gen.Database
	logger *zap.Logger
}

func New(opts Options) (*Repo, error) {
	if opts.Db == nil {
		log.Fatalf("")
	}
	lg := zap.L().Named("Repo")
	return &Repo{Db: opts.Db, logger: lg}, nil
}
