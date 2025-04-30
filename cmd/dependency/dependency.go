package dependency

import (
	"context"
	"fmt"

	"github.com/vovanwin/meetingsBot/config"
	"github.com/vovanwin/meetingsBot/internal/store"
	storegen "github.com/vovanwin/meetingsBot/internal/store/gen"
	"github.com/vovanwin/meetingsBot/pkg/logger"
	"go.uber.org/fx"
)

func ProvideConfig() (*config.Config, error) {
	return config.NewConfig()
}

func ProvideInitGlobalLogger(lifecycle fx.Lifecycle, cfg *config.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.MustInit(logger.NewOptions(
				cfg.Log.Level,
				cfg.Server.Env,
				logger.WithProductionMode(cfg.IsDebug()),
			))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sync()
			return nil
		},
	})
}

func ProvideStoreClient(lifecycle fx.Lifecycle, ctx context.Context, cfg config.Config) (client *storegen.Client, err error) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			client, err = store.NewSQLiteClient(store.NewSQLiteOptions(
				"./bot.db?&cache=shared&_fk=1",

				store.WithIsDebug(cfg.IsDebug()),
			))
			if err != nil {
				return fmt.Errorf("create psql client: %v", err)
			}

			if err = client.Schema.Create(ctx); err != nil {
				return fmt.Errorf("create schema: %v", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := client.Close(); err != nil {
				return fmt.Errorf("close pgsql: %v", err)
			}
			return nil
		},
	})

	return client, nil
}
