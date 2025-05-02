package dependency

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vovanwin/meetingsBot/config"
	"github.com/vovanwin/meetingsBot/internal/store"
	storegen "github.com/vovanwin/meetingsBot/internal/store/gen"
	"github.com/vovanwin/meetingsBot/pkg/logger"
)

// Специальный тип-заглушка, сигнализирующий об инициализации логгера
type LoggerReady struct{}

func ProvideConfig() (*config.Config, error) {
	return config.NewConfig()
}

//func ProvideInitGlobalLogger(cfg *config.Config) {
//	logger.MustInit(logger.NewOptions(
//		cfg.Log.Level,
//		cfg.Server.Env,
//		logger.WithProductionMode(cfg.IsProduction()),
//	))
//	defer logger.Sync()
//}

func ProvideInitGlobalLogger(lifecycle fx.Lifecycle, cfg *config.Config) (LoggerReady, error) {
	logger.MustInit(logger.NewOptions(
		cfg.Log.Level,
		cfg.Server.Env,
		logger.WithProductionMode(cfg.IsProduction()),
	))

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			zap.L().Named("logger").Info("остановка логгера")
			logger.Sync()
			return nil
		},
	})
	return LoggerReady{}, nil
}

//func ProvideStoreClient(lifecycle fx.Lifecycle, cfg *config.Config) (*storegen.Database, error) {
//
//	var (
//		client   *storegen.Client
//		database *storegen.Database
//	)
//	// Регистрация OnStart/OnStop после создания клиента
//	lifecycle.Append(fx.Hook{
//		OnStart: func(ctx context.Context) error {
//			client, err := store.NewSQLiteClient(store.NewSQLiteOptions(
//				"./bot.db?cache=shared&_fk=1",
//				store.WithIsDebug(cfg.IsDebug()),
//			))
//			if err != nil {
//				return fmt.Errorf("create sqlite client: %v", err)
//			}
//
//			database = storegen.NewDatabase(client)
//			return nil
//		},
//		OnStop: func(ctx context.Context) error {
//			if err := client.Close(); err != nil {
//				return fmt.Errorf("close sqlite: %v", err)
//			}
//			return nil
//		},
//	})
//
//	return database, nil
//}

func ProvideStoreClient(lifecycle fx.Lifecycle, cfg *config.Config, _ LoggerReady) (*storegen.Database, error) {
	client, err := store.NewSQLiteClient(store.NewSQLiteOptions(
		"./bot.db?cache=shared&_fk=1",
		store.WithIsDebug(cfg.IsDebug()),
	))
	if err != nil {
		return nil, fmt.Errorf("create sqlite client: %v", err)
	}

	database := storegen.NewDatabase(client)
	// Регистрация OnStart/OnStop после создания клиента
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := client.Close(); err != nil {
				return fmt.Errorf("close sqlite: %v", err)
			}
			return nil
		},
	})

	return database, nil
}
