package handlers

import (
	"github.com/vovanwin/meetingsBot/internal/store/gen"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/keyboards"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type Handlers struct {
	*TelegramBot
	db *gen.Database
}

func NewHandlers(bot *TelegramBot, db *gen.Database) *Handlers {
	return &Handlers{
		TelegramBot: bot,
		db:          db,
	}
}
func (h *Handlers) RegisterAdminHandlers() {
	// Создание нового сбора
	h.Bot.Handle("/new_event", func(c telebot.Context) error {
		h.Lg.Debug("Обработка события new_event")
		// Отправка сообщения с кнопками
		return c.Send(
			"🏐 Новый сбор создан!\n"+
				"Нажмите кнопку ниже, чтобы записаться:",
			keyboards.EventKeyboard("1"),
		)
	})

	// Пример простого хендлера
	h.Bot.Handle("/start", func(c telebot.Context) error {
		h.Lg.Debug("Обработка события start")
		return c.Send("Привет! Я бот 🤖")
	})

	h.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		user := c.Sender()
		text := c.Text()
		chat := c.Chat()

		// Пример логирования и проверки
		zap.L().Info("Новое сообщение",
			zap.String("user", user.Username),
			zap.Int64("user_id", user.ID),
			zap.Int64("chat_id", chat.ID),
			zap.String("text", text),
		)

		// Пример простой фильтрации
		if text == "запрещенное слово" {
			// удалить сообщение
			_ = c.Delete()
			// предупредить
			return c.Send("Нельзя писать запрещённые слова!")
		}

		// можно ничего не отправлять, если не нужно
		return nil
	})
}
