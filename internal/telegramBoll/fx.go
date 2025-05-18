package telegramBoll

import (
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/repository"
	"go.uber.org/fx"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/handlers"
)

var Module = fx.Module("telegramBoll",
	fx.Provide(
		handlers.ProvideBot,  // инициализация бота
		handlers.NewHandlers, // обработчики хуков и кнопок
		repository.New,
	),

	fx.Invoke(handlers.StartBot), // старт бота, запускать последним
)
