package handlers

import (
	"context"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/keyboards"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/repository"
)

type Handlers struct {
	*TelegramBot
	rep *repository.Repo
}

var (
	payMarkup     = telebot.ReplyMarkup{}
	btnFree       = payMarkup.Data("Бесплатно", "TYPE_PAY", "FREE")
	btnSplitEqual = payMarkup.Data("Разделить на всем", "TYPE_PAY", "SPLIT")
	btnFixedPer   = payMarkup.Data("Фиксирована", "TYPE_PAY", "FIXED")
)

func NewHandlers(bot *TelegramBot, rep *repository.Repo) *Handlers {
	payMarkup.Inline(
		payMarkup.Row(btnFree, btnSplitEqual, btnFixedPer),
	)

	handlers := &Handlers{
		TelegramBot: bot,
		rep:         rep,
	}
	go handlers.StartActiveMeetingsUpdater()

	return handlers
}

func (h *Handlers) start(c telebot.Context) error {
	h.Lg.Debug("Обработка события start")
	ctx := context.Background()
	h.rep.CreateUser(ctx, dto.CreateUser{
		ID:       c.Sender().ID,
		Username: c.Sender().Username,
	})

	c.Send("Привет! Я бот 🤖 для создания встреч и отслеживания участников")
	rules := `📌 Правила использования бота:

			 1. Для создания встречи:
			    - Бот должен быть добавлен в чат с правами на создание сообщений
			    - Создатель встречи должен быть администратором чата

			 2. Настройка встречи:
			    - Укажите описание встречи
			    - Выберите тип оплаты (бесплатно/платно)
			    - Для платных встреч укажите стоимость:
			      • Фиксированная - одинаковая для всех
			      • Поделенная - стоимость делится между участниками

			 3. Отправка в чат:
			    - Встречу можно отправить только в чаты, где вы администратор`

	err := h.Bot.SetCommands([]telebot.Command{
		{Text: "start", Description: "Начать работу"},
		{Text: "create", Description: "Создать встречу"},
		{Text: "edit", Description: "Редактировать встречу"},
		{Text: "admin", Description: "Админстраторские команды"},
	}, &telebot.CommandScope{Type: telebot.CommandScopeAllPrivateChats})
	if err != nil {
		h.Lg.Error("Не удалось установить команды", zap.Error(err))
	}

	return c.Send(rules, keyboards.EventStartKeyboard())
}
