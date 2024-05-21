package handlers

import (
	"JotunBack/handlers/botH"
	"JotunBack/model"
	"JotunBack/repository"
	"JotunBack/server"
	"JotunBack/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InLineKeyboardHandler(bot *tgbotapi.BotAPI, acState *model.ACState,
	userRepo *repository.UserRepository, update tgbotapi.Update,
	hub *server.Hub) {

	callbackData := update.CallbackQuery.Data

	switch callbackData {
	case "decrease":
		acState.TargetTemp--
	case "increase":
		acState.TargetTemp++
	case "cool":
		acState.Config.Mode = 1
	case "heat":
		acState.Config.Mode = 2
	case "dry":
		acState.Config.Mode = 3
	case "temperature":
		return
	case "start":
		StartAC(acState, hub, userRepo, update)
		return
	case "stop":
		err := botH.TurnOffAc(update.CallbackQuery.From.UserName, userRepo, hub)
		if err != nil {
			return
		}
		acState.Stop = true
	}

	isOnline := hub.GetConnectionByID(update.CallbackQuery.From.UserName) != nil
	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
		ui.ConfigMsg(acState, isOnline))

	keyboard := ui.ConfigKeyboard(acState)
	msg.ReplyMarkup = &keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
