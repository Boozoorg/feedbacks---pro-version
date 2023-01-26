package handlers

import (
	"feedbacks/models"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//_______________M____________SendToSupportGroup_______________M____________
func SenderBot(client *Chat) {
	var MessageToSupport = models.MessageToBot{
		ProductID: 1,
		RoomID:    client.user.Room_id,
		Fio:       client.user.Name,
		Message:   "made in heaven",
		Time:      time.Now().Format("15:04"),
	}
	//Tmken
	bot, err := tgbotapi.NewBotAPI("5738565649:AAHnDyp3XFkViRvCa-skF3zN-FSp7n7diEU")
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	escapedMessage := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, MessageToSupport.Message)
	escapedUserName := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, MessageToSupport.Fio)

	// updates := bot.GetUpdatesChan(updateConfig)
	msg := tgbotapi.NewMessage(-534115585,
		fmt.Sprintf("__ID клиента №: %d__\n\n"+
			"*Имя клиента:* \n_%v_\n\n"+
			"*Текст:*\n_%v_\n\n"+
			"*Время отправки:* \n_%v_\n\n"+
			"[Перейти к чату](%v)",

			MessageToSupport.RoomID,
			escapedUserName,
			escapedMessage,
			MessageToSupport.Time,
			"https://www.youtube.com/watch?v=dQw4w9WgXcQ"))

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
