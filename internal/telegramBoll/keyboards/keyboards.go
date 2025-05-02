package keyboards

import "gopkg.in/telebot.v4"

// Клавиатура для записи на сбор
func EventKeyboard(eventID string) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{}

	// Кнопки "+1" до "+5" и "Отмена"
	btnPlus1 := kb.Data("✅ +1", "vote", eventID, "1")
	btnPlus2 := kb.Data("+2", "vote", eventID, "2")
	btnCancel := kb.Data("❌ Отмена", "vote_cancel", eventID)

	kb.Inline(
		kb.Row(btnPlus1, btnPlus2),
		kb.Row(btnCancel),
	)

	return kb
}
