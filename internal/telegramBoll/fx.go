package telegramBoll

import (
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/handlers"
	"go.uber.org/fx"
)

var Module = fx.Module("telegramBoll",
	fx.Provide(
		handlers.ProvideBot,  // инициализация бота
		handlers.NewHandlers, //обработчики хуков и кнопок
	),

	fx.Invoke(handlers.StartBot), //старт бота, запускать последним
)
