package handlers

//func RegisterEventHandlers(bot *telebot.Bot, client *ent.Client) {
//	// Обработка нажатия "+1", "+2" и т.д.
//	bot.Handle(&telebot.Btn{Unique: "vote"}, func(c telebot.Context) error {
//		eventID, _ := strconv.Atoi(c.Data())
//		userID := c.Sender().ID
//		count, _ := strconv.Atoi(c.Args()[1])
//
//		// Обновляем голос в БД
//		_, err := client.Vote.
//			Create().
//			SetCount(count).
//			SetUserID(userID).
//			SetEventID(eventID).
//			OnConflict().
//			UpdateCount().
//			Exec(context.Background())
//
//		if err != nil {
//			return c.Respond(&telebot.CallbackResponse{Text: "Ошибка!"})
//		}
//
//		return c.Respond(&telebot.CallbackResponse{Text: "Записал!"})
//	})
//
//	// Отмена голоса
//	bot.Handle(&telebot.Btn{Unique: "vote_cancel"}, func(c telebot.Context) error {
//		// ... (удаление записи из БД)
//	})
//}
