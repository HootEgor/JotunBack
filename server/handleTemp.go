package server

import (
	"JotunBack/model"
	"JotunBack/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func HandleTemperature(acState *model.ACState, hub *Hub, currentTemp float32) {
	CoolMode := 1
	HeatMode := 2
	newTemp := float32(acState.Config.Degrees)

	prevTemp := acState.CurrentTemp

	if time.Now().After(acState.NextCheck) {
		if acState.Config.Mode == CoolMode {
			if currentTemp > prevTemp {
				newTemp -= 1
			}
		} else if acState.Config.Mode == HeatMode {
			if currentTemp < prevTemp {
				newTemp += 1
			}
		}

		acState.CurrentTemp = currentTemp
		acState.Config.Degrees = int(newTemp)
		err := hub.SendACConfig(acState.Config)
		if err != nil {
			return
		}

		acState.NextCheck = time.Now().Add(1 * time.Minute)
	}

	isOnline := hub.GetConnectionByID(acState.Username) != nil
	msg := tgbotapi.NewEditMessageText(acState.ChatID, acState.MsgID,
		ui.StateMsg(acState, isOnline))

	keyboard := ui.StateKeyboard(acState)
	msg.ReplyMarkup = &keyboard

	editedMsg, err := acState.Bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	acState.MsgID = editedMsg.MessageID
	log.Println("Temperature:", prevTemp, "->", currentTemp, "->", acState.TargetTemp)
	acState.EmojiNum += 1
	if acState.EmojiNum > 3 {
		acState.EmojiNum = 0
	}
}
