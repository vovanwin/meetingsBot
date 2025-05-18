package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dbsqlc"
	"go.uber.org/zap"

	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
)

func (r *Repo) CreateMeeting(ctx context.Context, dto dto.CreateMeeting) (dbsqlc.CreateMeetingRow, error) {
	p, err := r.Db.CreateMeeting(ctx, dbsqlc.CreateMeetingParams{
		Max: sql.NullInt64{
			Int64: dto.Limit,
			Valid: true,
		},
		Cost: sql.NullInt64{
			Int64: dto.Cost,
			Valid: true,
		},
		Message: sql.NullString{
			String: dto.Msg,
			Valid:  true,
		},
		OwnerID: dto.OwnerID,
		TypePay: dto.TypePay,
		Status:  dto.Status,
		Code:    dto.Code,
	})

	if err != nil {
		r.logger.Error("ошибка создания встречи", zap.Error(err))
		return p, fmt.Errorf("query problem: %v", err)
	}

	return p, nil
}

func (r *Repo) GetMeeting(ctx context.Context, id int64) (dbsqlc.GetMeetingRow, error) {
	p, err := r.Db.GetMeeting(ctx, id)
	if err != nil {
		r.logger.Error("ошибка получения встречи", zap.Error(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) GetMeetingByCode(ctx context.Context, code string) (dbsqlc.GetMeetingByCodeRow, error) {
	p, err := r.Db.GetMeetingByCode(ctx, code)
	if err != nil {
		r.logger.Error("ошибка получения встречи", zap.Error(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) CreateUser(ctx context.Context, dto dto.CreateUser) (dbsqlc.User, error) {
	user, err := r.Db.GetUser(ctx, dto.ID)

	if err != sql.ErrNoRows {
		if user.Username != dto.Username {
			err := r.Db.UpdateUsername(ctx, dbsqlc.UpdateUsernameParams{
				Username: dto.Username,
				ID:       dto.ID,
			})
			if err != nil {
				r.logger.Error("ошибка обновления пользователя", zap.Error(err))
			}
		}
		return user, nil
	}

	isOwner := false
	if 984891975 == dto.ID { // Аккаунт админа, да id захардкожен
		isOwner = true
	}

	p, err := r.Db.CreateUser(ctx, dbsqlc.CreateUserParams{
		ID:       dto.ID,
		Username: dto.Username,
		IsOwner:  isOwner,
	})

	if err != nil {
		r.logger.Debug("ошибка создания пользователя", zap.Error(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil

}

func (r *Repo) GetUsers(ctx context.Context) ([]dbsqlc.User, error) {
	p, err := r.Db.GetUsers(ctx)
	if err != nil {
		r.logger.Error("ошибка получения пользователей", zap.Error(err))
		return p, fmt.Errorf("query problem: %v", err)
	}
	return p, nil
}

func (r *Repo) GetUser(ctx context.Context, id int64) (dbsqlc.User, error) {
	p, err := r.Db.GetUser(ctx, id)
	if err != nil {
		r.logger.Error("ошибка получения пользователя", zap.Error(err))
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
		r.logger.Error("ошибка обновления статуса встречи", zap.Error(err))
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
				Count:     sql.NullInt64{},
			})
			if err != nil {
				r.logger.Error("VoteYes", zap.Error(err))
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
				Count:     sql.NullInt64{Valid: false},
			})
			if err != nil {
				r.logger.Error("VoteCancel: create", zap.Error(err))
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
		r.logger.Error("VoteCancel: update", zap.Error(err))
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
				Count:     sql.NullInt64{Int64: 1, Valid: true},
			})
			if err != nil {
				r.logger.Error("VotePlusAnother: create", zap.Error(err))
			}
			return err
		}
		return fmt.Errorf("query problem: %w", err)
	}
	// обновляем существующую запись
	newCount := um.Count.Int64 + 1
	err = r.Db.UpdateUserMeetingCount(ctx, dbsqlc.UpdateUserMeetingCountParams{
		Count:     sql.NullInt64{Int64: newCount, Valid: true},
		UserID:    userID,
		MeetingID: meetID,
	})
	if err != nil {
		r.logger.Error("VotePlusAnother: update", zap.Error(err))
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
			Count:     sql.NullInt64{Int64: newCount, Valid: true},
			UserID:    userID,
			MeetingID: meetID,
		})
		if err != nil {
			r.logger.Error("VoteMinusAnother: update", zap.Error(err))
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
			r.logger.Error("CreateChat: create", zap.Error(err))
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
			r.logger.Error("CreateChat: link", zap.Error(err))
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
		r.logger.Error("GetChatMeetingAllChatWithMeeting:", zap.Error(err))
		return nil, err
	}
	return meeting, nil
}
