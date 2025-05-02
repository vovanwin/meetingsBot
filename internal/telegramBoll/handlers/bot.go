package handlers

import (
	"context"
	"log"
	"time"

	"github.com/vovanwin/meetingsBot/cmd/dependency"
	"github.com/vovanwin/meetingsBot/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

const name = "telegramBoll"

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

func ProvideBot(cfg *config.Config, _ dependency.LoggerReady) (*TelegramBot, error) {
	lg := zap.L().Named(name)
	lg.Info("Инициализация бота")
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
		return nil, err
	}

	return &TelegramBot{
		Bot: bot,
		Lg:  lg,
	}, nil
}
