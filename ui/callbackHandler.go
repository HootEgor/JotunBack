package ui

import (
	"JotunBack/model"
	"JotunBack/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InLineKeyboardHandler(bot *tgbotapi.BotAPI, acState *model.ACState,
	userRepo *repository.UserRepository, update tgbotapi.Update) {

	callbackData := update.CallbackQuery.Data

	switch callbackData {
	case "decrease":
		acState.Config.Degrees--
	case "increase":
		acState.Config.Degrees++
	case "cool":
		acState.Config.Mode = 1
	case "heat":
		acState.Config.Mode = 2
	case "temperature":
		return
	}

	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
		ConfigMsg(acState))

	keyboard := ConfigKeyboard(acState)
	msg.ReplyMarkup = &keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
