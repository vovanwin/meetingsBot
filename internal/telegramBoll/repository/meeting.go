package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dbsqlc"
	"github.com/vovanwin/meetingsBot/pkg/fxslog/sl"
	"go.uber.org/zap"
	"log/slog"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
)

func (r *Repo) CreateMeeting(ctx context.Context, dto dto.CreateMeeting) (dbsqlc.CreateMeetingRow, error) {
	p, err := r.Db.CreateMeeting(ctx, dbsqlc.CreateMeetingParams{
		Max: pgtype.Int8{
			Int64: dto.Limit,
			Valid: true,
		},
		Cost: pgtype.Int8{
			Int64: dto.Cost,
			Valid: true,
		},
		Message: pgtype.Text{
			String: dto.Msg,
			Valid:  true,
		},
		OwnerID: dto.OwnerID,
		TypePay: dto.TypePay,
		Status:  dto.Status,
		Code:    dto.Code,
	})

	if err != nil {
		slog.Error("ошибка создания встречи", sl.Err(err))
		return p, fmt.Errorf("query problem: %v", err)
	}

	return p, nil
}

func (r *Repo) GetMeeting(ctx context.Context, id int64) (dbsqlc.GetMeetingRow, error) {
	p, err := r.Db.GetMeeting(ctx, id)
	if err != nil {
		slog.Error("ошибка получения встречи", sl.Err(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) GetMeetingByCode(ctx context.Context, code string) (dbsqlc.GetMeetingByCodeRow, error) {
	p, err := r.Db.GetMeetingByCode(ctx, code)
	if err != nil {
		slog.Error("ошибка получения встречи", sl.Err(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) CreateUser(ctx context.Context, data dto.CreateUser) (dto.UserRow, error) {
	user, err := r.Db.GetUser(ctx, data.ID)

	if err != sql.ErrNoRows {
		if user.Username != data.Username {
			err := r.Db.UpdateUsername(ctx, dbsqlc.UpdateUsernameParams{
				Username: data.Username,
				UserID:   data.ID,
			})
			if err != nil {
				slog.Error("ошибка обновления пользователя", sl.Err(err))
			}
		}
		return dto.UserRow{
			ID:       user.ID,
			Username: user.Username,
			IsOwner:  user.IsOwner,
			Nickname: user.Nickname,
		}, nil
	}

	isOwner := false
	if 984891975 == data.ID { // Аккаунт админа, да id захардкожен
		isOwner = true
	}

	p, err := r.Db.CreateUser(ctx, dbsqlc.CreateUserParams{
		ID:       data.ID,
		Username: data.Username,
		IsOwner:  isOwner,
		Nickname: pgtype.Text{
			String: data.Nickname,
			Valid:  true,
		},
	})

	if err != nil {
		slog.Error("ошибка создания пользователя", sl.Err(err))
		return dto.UserRow{}, fmt.Errorf("query problem: %v", err)
	}
	return dto.UserRow{
		ID:       p.ID,
		Username: p.Username,
		IsOwner:  p.IsOwner,
		Nickname: p.Nickname,
	}, nil

}

func (r *Repo) GetUsers(ctx context.Context) ([]dbsqlc.GetUsersRow, error) {
	p, err := r.Db.GetUsers(ctx)
	if err != nil {
		slog.Error("ошибка получения пользователей", sl.Err(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) GetUser(ctx context.Context, id int64) (dbsqlc.GetUserRow, error) {
	p, err := r.Db.GetUser(ctx, id)
	if err != nil {
		slog.Error("ошибка получения пользователя", sl.Err(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) UpdateMeetingStatus(ctx context.Context, dto dto.UpdateMeetingStatus) error {
	err := r.Db.UpdateMeetingStatus(ctx, dbsqlc.UpdateMeetingStatusParams{
		Status: dto.Status,
		Code:   dto.Code,
	})
	if err != nil {
		slog.Error("ошибка обновления статуса встречи", sl.Err(err))
		return fmt.Errorf("query problem: %v", err)
	}
	return nil
}

func (r *Repo) VoteYes(ctx context.Context, userID, meetID int64) error {
	_, err := r.Db.GetUserMeeting(ctx, dbsqlc.GetUserMeetingParams{
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		r.logger.Debug("GetUserMeeting", zap.Error(err))
		if err == sql.ErrNoRows {
			r.logger.Debug("GetUserMeeting ErrNoRows", zap.Error(err))
			_, err := r.Db.CreateUserMeeting(ctx, dbsqlc.CreateUserMeetingParams{
				UserID:    userID,
				MeetingID: meetID,
				Status:    dto.VoteStatusУчавствует.String(),
				Count:     pgtype.Int8{},
			})
			if err != nil {
				slog.Error("VoteYes", sl.Err(err))
				return err
			}
			return nil
		}
	}
	err = r.Db.UpdateUserMeetingStatus(ctx, dbsqlc.UpdateUserMeetingStatusParams{
		Status:    dto.VoteStatusУчавствует.String(),
		UserID:    userID,
		MeetingID: meetID,
	})
	return nil
}

// VoteCancel ставит статус "CANCEL" в user_meetings
func (r *Repo) VoteCancel(ctx context.Context, code string, userID, meetID int64) error {
	// пытаемся получить существующую запись
	_, err := r.Db.GetUserMeeting(ctx, dbsqlc.GetUserMeetingParams{
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// если нет — создаём
			_, err = r.Db.CreateUserMeeting(ctx, dbsqlc.CreateUserMeetingParams{
				UserID:    userID,
				MeetingID: meetID,
				Status:    dto.VoteStatusНет.String(),
				Count:     pgtype.Int8{Valid: false},
			})
			if err != nil {
				slog.Error("VoteCancel: create", sl.Err(err))
			}
			return err
		}
		return fmt.Errorf("query problem: %w", err)
	}
	// если есть — обновляем статус
	err = r.Db.UpdateUserMeetingStatus(ctx, dbsqlc.UpdateUserMeetingStatusParams{
		Status:    dto.VoteStatusНет.String(),
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		slog.Error("VoteCancel: update", sl.Err(err))
		return fmt.Errorf("query problem: %w", err)
	}
	return nil
}

// VotePlusAnother увеличивает поле count на 1
func (r *Repo) VotePlusAnother(ctx context.Context, code string, userID, meetID int64) error {
	um, err := r.Db.GetUserMeeting(ctx, dbsqlc.GetUserMeetingParams{
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// создаём с count=1
			_, err = r.Db.CreateUserMeeting(ctx, dbsqlc.CreateUserMeetingParams{
				UserID:    userID,
				MeetingID: meetID,
				Status:    "CANCEL", // как в Ent: стартовый статус
				Count:     pgtype.Int8{Int64: 1, Valid: true},
			})
			if err != nil {
				slog.Error("VotePlusAnother: create", sl.Err(err))
			}
			return err
		}
		return fmt.Errorf("query problem: %w", err)
	}
	// обновляем существующую запись
	newCount := um.Count.Int64 + 1
	err = r.Db.UpdateUserMeetingCount(ctx, dbsqlc.UpdateUserMeetingCountParams{
		Count:     pgtype.Int8{Int64: newCount, Valid: true},
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		slog.Error("VotePlusAnother: update", sl.Err(err))
		return fmt.Errorf("query problem: %w", err)
	}
	return nil
}

// VoteMinusAnother уменьшает поле count на 1, если >0
func (r *Repo) VoteMinusAnother(ctx context.Context, code string, userID, meetID int64) error {
	um, err := r.Db.GetUserMeeting(ctx, dbsqlc.GetUserMeetingParams{
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // нечего уменьшать
		}
		return fmt.Errorf("query problem: %w", err)
	}
	if um.Count.Valid && um.Count.Int64 > 0 {
		newCount := um.Count.Int64 - 1
		err = r.Db.UpdateUserMeetingCount(ctx, dbsqlc.UpdateUserMeetingCountParams{
			Count:     pgtype.Int8{Int64: newCount, Valid: true},
			UserID:    userID,
			MeetingID: meetID,
		})
		if err != nil {
			slog.Error("VoteMinusAnother: update", sl.Err(err))
			return fmt.Errorf("query problem: %w", err)
		}
	}
	return nil
}

// CreateChat создаёт чат, если его нет, и связь chat_meetings
func (r *Repo) CreateChat(ctx context.Context, dto dto.CreateChat) error {
	// проверяем чат

	_, err := r.Db.GetChat(ctx, dto.ChatID)
	if err != nil && err == sql.ErrNoRows {
		_, err = r.Db.CreateChat(ctx, dbsqlc.CreateChatParams{
			ID:        dto.ChatID,
			Title:     dto.ChatTitle,
			IsMeeting: true,
		})
		if err != nil {
			slog.Error("CreateChat: create", sl.Err(err))
			return fmt.Errorf("query problem: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("query problem: %w", err)
	}

	// проверяем связь chat_meetings
	_, err = r.Db.GetChatMeeting(ctx, dbsqlc.GetChatMeetingParams{
		ChatID:    dto.ChatID,
		MeetingID: dto.MeetID,
	})
	if err != nil && err == sql.ErrNoRows {
		_, err = r.Db.CreateChatMeeting(ctx, dbsqlc.CreateChatMeetingParams{
			ChatID:    dto.ChatID,
			MeetingID: dto.MeetID,
			MessageID: dto.MessageID,
		})
		if err != nil {
			slog.Error("CreateChat: link", sl.Err(err))
			return fmt.Errorf("query problem: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("query problem: %w", err)
	}

	return nil
}

// GetChatMeeting возвращает запись из chat_meetings
func (r *Repo) GetChatMeeting(ctx context.Context, chatID, meetID int64) (dbsqlc.ChatMeeting, error) {
	return r.Db.GetChatMeeting(ctx, dbsqlc.GetChatMeetingParams{
		ChatID:    chatID,
		MeetingID: meetID,
	})
}

// GetUsersMeetings возвращает всех голосовавших участников встречи вместе с данными User
func (r *Repo) GetUsersMeetings(ctx context.Context, MeetingID int64) ([]dbsqlc.GetUsersMeetingsRow, error) {
	return r.Db.GetUsersMeetings(ctx, MeetingID)
}

// GetUsersMeetings Обновить связь между чатом и встречей. Привязку сообщений
func (r *Repo) UpdateChatMeeting(ctx context.Context, dto dbsqlc.UpdateChatMeetingParams) (dbsqlc.ChatMeeting, error) {
	return r.Db.UpdateChatMeeting(ctx, dbsqlc.UpdateChatMeetingParams{
		MessageID:      dto.MessageID,
		WhereMeetingID: dto.WhereMeetingID,
		WhereChatID:    dto.WhereChatID,
	})
}

func (r *Repo) GetChatMeetingAllChatWithMeeting(ctx context.Context, meetingID int64) ([]dbsqlc.ChatMeeting, error) {
	meeting, err := r.Db.GetChatMeetingAllChatWithMeeting(ctx, meetingID)
	if err != nil {
		slog.Error("GetChatMeetingAllChatWithMeeting:", sl.Err(err))
		return nil, err
	}
	return meeting, nil
}

// Получать все встречи которые обновлялись больше суток назад, чтобы переместить их в начало чата
func (r *Repo) GetMeetingsForUpdateTime(ctx context.Context) ([]dbsqlc.GetMeetingsForUpdateTimeRow, error) {
	meeting, err := r.Db.GetMeetingsForUpdateTime(ctx)
	if err != nil {
		slog.Error("GetChatMeetingAllChatWithMeeting:", sl.Err(err))
		return nil, err
	}
	return meeting, nil
}
