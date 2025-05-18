package handlers

import (
	"context"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/Tdep"
	"log"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"

	"github.com/vovanwin/meetingsBot/cmd/dependency"
	"github.com/vovanwin/meetingsBot/config"
)

const name = "telegramBoll"

var activeMeetingCodes = make(map[string]struct{})

type TelegramBot struct {
	Bot *telebot.Bot
	Lg  *zap.Logger
}

func StartBot(lc fx.Lifecycle, bot *TelegramBot, handler *Handlers) {
	bot.Lg.Debug("Старт бота")
	handler.RegisterAdminHandlers()
	go bot.Bot.Start()

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			bot.Lg.Debug("Стоп бота")
			bot.Bot.Stop()
			return nil
		},
	})
}

func ProvideBot(cfg *config.Config, _ dependency.LoggerReady) (*TelegramBot, *Tdep.TelegramLogger, error) {
	lg := zap.L().Named(name)
	lg.Info("Инициализация бота")
	activeMeetingCodes = make(map[string]struct{})

	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: nil,
		OnError: func(err error, c telebot.Context) {
			log.Println("Bot error:", err)
		},
	}

	if cfg.Telegram.UseWebhook {
		lg.Info("📡 Using webhook mode")
		pref.Poller = &telebot.Webhook{
			Listen: cfg.Listen,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: cfg.Telegram.PublicURL + cfg.Telegram.Webhook,
			},
		}
	} else {
		lg.Info("🕵️ Using polling mode")
		pref.Poller = &telebot.LongPoller{Timeout: 5 * time.Second}
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, nil, err
	}

	return &TelegramBot{
			Bot: bot,
			Lg:  lg,
		}, &Tdep.TelegramLogger{
			Lg: lg,
		}, nil
}
