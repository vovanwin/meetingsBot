package dependency

import (
	"context"
	"github.com/vovanwin/meetingsBot/pkg/clients/sqlite"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vovanwin/meetingsBot/config"
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

func ProvideStoreClient(lifecycle fx.Lifecycle, cfg *config.Config) (*sqlite.SQLiteClient, error) {
	log := zap.L().Named("logger")
	// собираем DSN с оптимальными параметрами
	// создаём клиента, без ping’а — это сделаем в OnStart
	lite, err := sqlite.ConnectSQLite(context.Background(), cfg.SQLite)
	if err != nil {
		return nil, err
	}

	// регистрируем жизненный цикл
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Проверка соединения с SQLite")
			// перепинговать с таймаутом
			return lite.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Закрытие соединения с SQLite")
			return lite.Close()
		},
	})

	return lite, nil
}
