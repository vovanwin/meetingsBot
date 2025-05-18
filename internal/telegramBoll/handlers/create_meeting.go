package handlers

import (
	"context"
	"fmt"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/keyboards"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/telebot.v4"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
	"github.com/vovanwin/meetingsBot/pkg/helper"
)

// Определяем типы операций и состояний
type OperationType string

const (
	OperationCreate  OperationType = "create"
	OperationEdit    OperationType = "edit"
	OperationConfirm OperationType = "confirm"
)

type UserSession struct {
	Operation    OperationType
	CurrentState string
	Data         interface{}
}

// Структуры данных для разных операций
type MeetingDraft struct {
	Description string
	Limit       int
	PaymentType string
	Amount      int
}

var sessions = make(map[int64]*UserSession)

func (h *Handlers) RegisterAdminHandlers() {
	h.Bot.Handle("/start", h.start)
	h.Bot.Handle("/create", h.handleCreate)
	h.Bot.Handle(telebot.OnText, h.handleText)
	h.Bot.Handle(telebot.OnCallback, h.handleCallback)
}

func (h *Handlers) handleCreate(c telebot.Context) error {
	session := &UserSession{
		Operation:    OperationCreate,
		CurrentState: "description",
		Data:         &MeetingDraft{},
	}
	sessions[c.Sender().ID] = session
	return c.Send("Введите описание встречи:", telebot.ForceReply)
}

// Обработчики сообщений
func (h *Handlers) handleText(c telebot.Context) error {
	h.Lg.Debug("id message", zap.Any("message", c.Message().ID))

	h.Lg.Debug("Обработка события handleText", zap.String("text", c.Text()))
	isExist := isMeetingActive(c.Text())
	h.Lg.Debug("проверка кеш массива кодов", zap.Any("кеш", activeMeetingCodes))
	h.Lg.Debug("проверка кеш массива кодов", zap.Any("кеш", isExist))

	if isExist {
		h.Lg.Debug("проверка пройдена", zap.Any("кеш", isExist))

		return h.showMeeting(c, c.Text())
	}

	// выполняется сценарий
	session, ok := sessions[c.Sender().ID]
	if !ok {
		return c.Send("Начните операцию с помощью команды")
	}

	switch session.Operation {
	case OperationCreate:
		return h.handleCreateText(c, session)
	}
	return nil
}

func (h *Handlers) handleCreateText(c telebot.Context, session *UserSession) error {
	draft := session.Data.(*MeetingDraft)
	switch session.CurrentState {
	case "description":
		draft.Description = c.Text()
		session.CurrentState = "limit"
		return c.Send("Введите лимит участников:", telebot.ForceReply)

	case "limit":
		limit, err := strconv.Atoi(c.Text())
		if err != nil {
			return c.Send("Введите число")
		}
		draft.Limit = limit
		session.CurrentState = "payment"
		return h.showPaymentOptions(c)

	case "payment_amount":
		amount, err := strconv.Atoi(c.Text())
		if err != nil {
			return c.Send("Введите число")
		}
		draft.Amount = amount
		if draft.PaymentType == "SPLIT" {
			draft.Amount = (amount + draft.Limit - 1) / draft.Limit
		}

		return h.finalizeMeeting(c, draft)
	}
	return nil
}

// Обработчики callback
func (h *Handlers) handleCallback(c telebot.Context) error {
	raw := c.Callback().Data
	data := strings.TrimSpace(raw)
	dataParts := strings.Split(data, "|")
	switch data {
	case "create_meeting":
		return h.handleCreate(c)
	}

	switch dataParts[0] {
	case "status_meeting":
		return h.StartMeeting(c)
	case "vote":
		return h.VoteMeeting(c)
	}

	session, ok := sessions[c.Sender().ID]
	if !ok {
		return c.Respond()
	}
	h.Lg.Debug("handleCallback", zap.Any("session", session))

	switch session.Operation {
	case OperationCreate:
		return h.handleCreateCallback(c, session, data)
	}
	return nil
}

func (h *Handlers) handleCreateCallback(c telebot.Context, session *UserSession, data string) error {
	draft := session.Data.(*MeetingDraft)

	switch {
	case data == "FREE":
		draft.PaymentType = "FREE"
		return h.finalizeMeeting(c, draft)
	case data == "SPLIT":
		draft.PaymentType = "SPLIT"
		session.CurrentState = "payment_amount"
		return c.Send("Введите общую сумму (пересчитывается от количества участников):")
	case data == "FIXED":
		draft.PaymentType = "FIXED"
		session.CurrentState = "payment_amount"
		return c.Send("Введите сумму с человека:")
	}
	return nil
}

// Вспомогательные методы
func (h *Handlers) showPaymentOptions(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(
			markup.Data("Бесплатно", "FREE"),
			markup.Data("Поровну", "SPLIT"),
			markup.Data("Фиксировано", "FIXED"),
		),
	)
	return c.Send("Выберите тип оплаты:", markup)
}

func (h *Handlers) finalizeMeeting(c telebot.Context, draft *MeetingDraft) error {
	delete(sessions, c.Sender().ID)
	// Здесь сохранение в БД
	ctx := context.Background()
	code, _ := helper.GenerateCode()
	meet, err := h.rep.CreateMeeting(ctx, dto.CreateMeeting{
		Limit:   int64(draft.Limit),
		Cost:    int64(draft.Amount),
		Msg:     draft.Description,
		OwnerID: c.Sender().ID,
		TypePay: draft.PaymentType,
		Status:  dto.StatusMeetingЧерновик.String(),
		Code:    code,
	})
	if err != nil {
		return err
	}
	h.refreshActiveMeetings(ctx)
	return c.Send(fmt.Sprintf(
		`Создана встреча!
Описание: %s
Лимит: %d
Оплата: %s
Сумма: %d`,
		meet.Message, meet.Max.Int64, meet.TypePay, meet.Cost,
	), keyboards.EventMeetingStartKeyboard(meet.Code))

}

func (h *Handlers) StartMeeting(c telebot.Context) error {
	raw := c.Data()
	data := strings.TrimSpace(raw)
	dataParts := strings.Split(data, "|")

	if dataParts[2] == "START" {
		//находим встречу по code и обновляем статус
		h.rep.UpdateMeetingStatus(context.Background(), dto.UpdateMeetingStatus{
			Code:   dataParts[1],
			Status: dto.StatusMeetingАктивная.String(),
		})
		c.Send(fmt.Sprintf("Уникальный код для вставки в чаты. просто вставьте в чат где присутсвует бот:  %s", dataParts[1]))
		return c.Respond(&telebot.CallbackResponse{Text: "Встреча начата!"})
	}

	if dataParts[2] == "END" {
		//находим встречу по code и обновляем статус
		h.rep.UpdateMeetingStatus(context.Background(), dto.UpdateMeetingStatus{
			Code:   dataParts[1],
			Status: dto.StatusMeetingОтменена.String(),
		})
		return c.Respond(&telebot.CallbackResponse{Text: "Встреча отменена!"})
	}
	if dataParts[3] == "COMPLETED" {
		//находим встречу по code и обновляем статус
		h.rep.UpdateMeetingStatus(context.Background(), dto.UpdateMeetingStatus{
			Code:   dataParts[1],
			Status: dto.StatusMeetingЗакончена.String(),
		})
		return c.Respond(&telebot.CallbackResponse{Text: "Встреча закончена!"})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "???"})
}

func (h *Handlers) showMeeting(c telebot.Context, code string) error {
	meet, err := h.rep.GetMeetingByCode(context.Background(), code)
	if err != nil {
		h.Lg.Error("ошибка получения встречи", zap.Error(err))
		return err
	}

	userVotes, err := h.rep.GetUsersMeetings(context.Background(), meet.ID)
	if err != nil {
		h.Lg.Error("ошибка получения участников", zap.Error(err))
		return err
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
			guests = append(guests, fmt.Sprintf("%s — %d гостей", username, v.Count))
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
		meet.Message,
		meet.Max,
		meet.TypePay,
		meet.Cost,
		nonEmptyList(participants),
		nonEmptyList(guests),
	)

	// Обновляем сообщение, если оно уже было
	chatMeeting, err := h.rep.GetChatMeeting(context.Background(), c.Chat().ID, meet.ID)
	if err == nil {
		msg := &telebot.Message{
			ID:   int(chatMeeting.MessageID),
			Chat: &telebot.Chat{ID: chatMeeting.ChatID},
		}
		_, err := h.Bot.Edit(msg, text, keyboards.EventKeyboard(meet.Code))
		if err != nil {
			h.Lg.Error("ошибка Edit", zap.Error(err))
		}
		return nil
	}

	// Иначе создаём новое
	send, err := h.Bot.Send(c.Chat(), text, keyboards.EventKeyboard(meet.Code))
	if err != nil {
		h.Lg.Error("ошибка отправки сообщения", zap.Error(err))
		return err
	}

	err = h.rep.CreateChat(context.Background(), dto.CreateChat{
		ChatID:    send.Chat.ID,
		ChatTitle: c.Chat().Title,
		MeetID:    meet.ID,
		MessageID: int64(send.ID),
	})
	if err != nil {
		h.Lg.Error("ошибка создания чата", zap.Error(err))
		return err
	}

	return nil
}

// вспомогательная функция, чтобы избежать пустых блоков
func nonEmptyList(lines []string) string {
	if len(lines) == 0 {
		return "—"
	}
	return strings.Join(lines, "\n")
}

func (h *Handlers) VoteMeeting(c telebot.Context) error {
	// 0 тип клавиатуры
	// 1 code
	// 2 event
	ctx := context.Background()
	raw := c.Data()
	data := strings.TrimSpace(raw)
	dataParts := strings.Split(data, "|")

	h.rep.CreateUser(context.Background(), dto.CreateUser{
		ID:       c.Sender().ID,
		Username: c.Sender().Username,
	})
	meet, _ := h.rep.GetMeetingByCode(ctx, dataParts[1])

	switch dataParts[2] {
	case keyboards.Yes:
		err := h.rep.VoteYes(ctx, c.Sender().ID, meet.ID)
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Неизвестная ошибка"})
		}

	case keyboards.Cancel:
		err := h.rep.VoteCancel(ctx, dataParts[1], c.Sender().ID, meet.ID)
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Неизвестная ошибка"})
		}
	case keyboards.MinusAnother:
		err := h.rep.VoteMinusAnother(ctx, dataParts[1], c.Sender().ID, meet.ID)
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Неизвестная ошибка"})
		}
	case keyboards.PlusAnother:
		err := h.rep.VotePlusAnother(ctx, dataParts[1], c.Sender().ID, meet.ID)
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Неизвестная ошибка"})
		}
	}

	h.showMeeting(c, dataParts[1])
	return nil
}
