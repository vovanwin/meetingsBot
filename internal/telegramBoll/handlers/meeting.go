package handlers

import (
	"context"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll/dto"
	"github.com/vovanwin/meetingsBot/pkg/fxslog/sl"
	"log/slog"
	"sync"
	"time"
)

var mu sync.RWMutex

// isMeetingActive проверяет наличие кода встречи в мапе
func isMeetingActive(code string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, exists := activeMeetingCodes[code]
	return exists
}

// обновляет мапу активных кодов встреч
func updateActiveMeetingCodes(codes []string) {
	newMap := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		newMap[code] = struct{}{}
	}

	mu.Lock()
	activeMeetingCodes = newMap
	mu.Unlock()
}

func (h *Handlers) StartActiveMeetingsUpdater() {
	ticker := time.NewTicker(time.Second * 5)
	ctx := context.Background()
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			h.refreshActiveMeetings(ctx)
		}
	}()
}

// refreshActiveMeetings обновляет мапу активных встреч (внутренний метод)
func (h *Handlers) refreshActiveMeetings(ctx context.Context) {
	// Здесь ты пишешь свою логику получения кодов из базы
	codes, err := h.rep.Db.GetMeetingsWithStatus(ctx, dto.StatusMeetingАктивная.String())
	if err != nil {
		slog.Error("Не удалось получить коды активных встреч", sl.Err(err))
		return
	}
	updateActiveMeetingCodes(codes)
}
