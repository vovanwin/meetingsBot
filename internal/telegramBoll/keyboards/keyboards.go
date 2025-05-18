package keyboards

import "gopkg.in/telebot.v4"

const (
	Yes          = "yes"         // голос за себя
	PlusAnother  = "plusAnother" // плюс люди со стороны
	MinusAnother = "minus"       // минус люди со стороны
	Cancel       = "cancel"      // отмена голоса за себя
)

// Клавиатура для записи на сбор
func EventKeyboard(eventID string) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{}

	// Кнопки "+1" до "+5" и "Отмена"
	btnPlus1 := kb.Data("✅ Иду", "vote", eventID, Yes)
	btnPlus2 := kb.Data("+1 - люди со стороны", "vote", eventID, PlusAnother)
	btnMinus3 := kb.Data("-1 - люди со стороны", "vote", eventID, MinusAnother)
	btnCancel := kb.Data("❌ Отмена личного голоса", "vote", eventID, Cancel)

	kb.Inline(
		kb.Row(btnPlus1),
		kb.Row(btnPlus2),
		kb.Row(btnMinus3),
		kb.Row(btnCancel),
	)

	return kb
}

func EventStartKeyboard() *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{}

	// Кнопки "+1" до "+5" и "Отмена"
	btnPlus1 := kb.Data("Создать Встречу", "create_meeting")
	btnPlus2 := kb.Data("Закрыть встречу", "close_meeting")
	btnCancel := kb.Data("Редактировать встречу", "edit_meeting")

	kb.Inline(
		kb.Row(btnPlus1, btnPlus2),
		kb.Row(btnCancel),
	)

	return kb
}

func EventMeetingStartKeyboard(code string) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{}

	btnPlus1 := kb.Data("Начать встречу", "status_meeting", code, "START")
	btnPlus2 := kb.Data("Закрыть встречу", "status_meeting", code, "END")
	btnPlus3 := kb.Data("Успешно завершена", "status_meeting", code, "COMPLETED")

	kb.Inline(
		kb.Row(btnPlus1, btnPlus2),
		kb.Row(btnPlus3),
	)

	return kb
}
