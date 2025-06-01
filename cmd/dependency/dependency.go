package dependency

import (
	"fmt"
	"github.com/vovanwin/meetingsBot/config"
	"github.com/vovanwin/meetingsBot/pkg/clients/postgres"
	"github.com/vovanwin/meetingsBot/pkg/fxslog"
	"log/slog"
)

// Специальный тип-заглушка, сигнализирующий об инициализации логгера

func ProvideConfig() (*config.Config, error) {
	return config.NewConfig()
}
func ProvideLogger(config *config.Config) (*slog.Logger, error) {
	opt := fxslog.NewOptions(fxslog.WithEnv(config.Env), fxslog.WithLevel(config.Level))
	return fxslog.NewLogger(opt)
}

func ProvidePgx(config *config.Config, logger *slog.Logger) (*postgres.Postgres, error) {
	opt := postgres.NewOptions(
		logger,
		config.PG.HostPG,
		config.PG.User,
		config.PG.Password,
		config.PG.DB,
		config.PG.Port,
		config.PG.Scheme,
		config.IsProduction(),
	)

	connect, err := postgres.New(opt)
	if err != nil {
		return nil, fmt.Errorf("create pgx connection: %w", err)
	}

	return connect, nil
}
