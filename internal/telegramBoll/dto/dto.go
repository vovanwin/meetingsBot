package dto

import "github.com/jackc/pgx/v5/pgtype"

//go:generate go-enum --names

// StatusMeeting описывает состояние встречи.
// ENUM(Активная, Отменена, Закончена, Черновик)
type StatusMeeting string

// TypePay описывает способ оплаты участия.
// ENUM(Фиксированная, Поровну, Бесплатно)
type TypePay string

// VoteStatus Учавствует в мероприятии или нет.
// ENUM(Учавствует, Нет, Думает)
type VoteStatus string

type CreateMeeting struct {
	Limit   int64
	Cost    int64
	Msg     string
	OwnerID int64
	TypePay string
	Status  string
	Code    string
}

type CreateUser struct {
	ID       int64
	Username string
	Nickname string
}

type UpdateMeetingStatus struct {
	Code   string
	Status string
}

type CreateChat struct {
	ChatID    int64
	ChatTitle string
	MeetID    int64
	MessageID int64
}
type UserRow struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	IsOwner  bool        `json:"is_owner"`
	Nickname pgtype.Text `json:"nickname"`
}
