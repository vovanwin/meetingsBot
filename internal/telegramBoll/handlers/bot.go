package handlers

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.uber.org/fx"
	"gopkg.in/telebot.v4"

	"github.com/vovanwin/meetingsBot/config"
)

const name = "telegramBoll"

var activeMeetingCodes = make(map[string]struct{})

type TelegramBot struct {
	Bot *telebot.Bot
}

func StartBot(lc fx.Lifecycle, bot *TelegramBot, handler *Handlers) {
	slog.Debug("Старт бота")
	handler.RegisterAdminHandlers()
	go bot.Bot.Start()

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			slog.Debug("Стоп бота")
			bot.Bot.Stop()
			return nil
		},
	})
}

func ProvideBot(cfg *config.Config) (*TelegramBot, error) {
	activeMeetingCodes = make(map[string]struct{})

	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: nil,
		OnError: func(err error, c telebot.Context) {
			log.Println("Bot error:", err)
		},
	}

	if cfg.Telegram.UseWebhook {
		slog.Info("📡 Using webhook mode")
		pref.Poller = &telebot.Webhook{
			Listen: cfg.Listen,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: cfg.Telegram.PublicURL + cfg.Telegram.Webhook,
			},
		}
	} else {
		slog.Info("🕵️ Using polling mode")
		pref.Poller = &telebot.LongPoller{Timeout: 5 * time.Second}
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{Bot: bot}, nil
}
