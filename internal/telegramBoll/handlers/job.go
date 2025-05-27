package handlers

import (
	"context"
	"fmt"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dbsqlc"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/keyboards"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"time"
)

// Запускает процесс удаления сообщений в чатах и перенос их в начало чата. Событие должно быть не закрытым.
func (h *Handlers) StartUpdateMessageInChat() {
	ticker := time.NewTicker(time.Second * 10)
	h.Lg.Info("запустился job StartUpdateMessageInChat")
	ctx := context.Background()
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			h.updateMessageInChat(ctx)
		}
	}()
}

// refreshActiveMeetings обновляет мапу активных встреч (внутренний метод)
func (h *Handlers) updateMessageInChat(ctx context.Context) {
	h.Lg.Info("Отработал job StartUpdateMessageInChat")
	// Здесь ты пишешь свою логику получения кодов из базы
	messages, err := h.rep.Db.GetMeetingsForUpdateTime(ctx)
	if err != nil {
		h.Lg.Error("Не удалось встречи", zap.Error(err))
		return
	}
	for _, v := range messages {

		// создаем новые, без уведомления
		text, _ := h.textMessage(v.MeetingID, v.Message.String, v.Max.Int64, v.TypePay, v.Cost.Int64)

		chat := &telebot.Chat{
			ID: v.ChatID,
		}

		send, err := h.Bot.Send(chat, text, keyboards.EventKeyboard(v.Code), telebot.Silent)
		if err != nil {
			h.Lg.Error("ошибка отправки сообщения", zap.Error(err))
		}

		//обновляем таблицу chat_meetings хранящую ссылки на сообщения
		h.rep.Db.UpdateChatMeeting(ctx, dbsqlc.UpdateChatMeetingParams{
			MessageID:      int64(send.ID),
			WhereMeetingID: v.MeetingID,
			WhereChatID:    v.ChatID,
		})
		h.rep.Db.UpdateMeetingUpdate(ctx, dbsqlc.UpdateMeetingUpdateParams{
			UpdatedAt:      time.Now(),
			WhereMeetingID: v.MeetingID,
		})

		// удаляем старые сообщения
		msg := &telebot.Message{
			ID:   int(v.MessageID),
			Chat: &telebot.Chat{ID: v.ChatID},
		}
		h.Bot.Delete(msg)
	}
}

func (h *Handlers) textMessage(meetID int64, message string, max int64, typePay string, cost int64) (string, error) {

	userVotes, err := h.rep.GetUsersMeetings(context.Background(), meetID)
	if err != nil {
		h.Lg.Error("ошибка получения участников", zap.Error(err))
		return "", err
	}

	var participants []string
	var guests []string

	for _, v := range userVotes {
		username := "@" + v.Username
		if v.Username == "" {
			username = fmt.Sprintf("id:%d", v.UserID)
		}

		if v.Status == dto.VoteStatusУчавствует.String() {
			participants = append(participants, username)
		}
		if v.Count.Int64 > 0 {
			guests = append(guests, fmt.Sprintf("%s — %d гостей", username, v.Count.Int64))
		}
	}

	text := fmt.Sprintf(`%s

Лимит участников: %d
Оплата: %s
Сумма: %d

Участники:
%s

Гости:
%s`,
		message,
		max,
		typePay,
		cost,
		nonEmptyList(participants),
		nonEmptyList(guests),
	)

	return text, nil
}
