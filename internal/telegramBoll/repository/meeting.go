package repository

import (
	"context"
	"fmt"

	store "github.com/vovanwin/meetingsBot/internal/store/gen"
)

func (r *Repo) First(ctx context.Context) (*store.User, error) {
	p, err := r.Db.User(ctx).Query().First(ctx)
	if store.IsNotFound(err) {
		r.logger.Debug(" Не найден пользователь")
		return p, fmt.Errorf("query problem: %v", err)
	}

	return p, nil
}
