package handlers

import (
	"JotunBack/model"
	"JotunBack/repository"
	"JotunBack/server"
	"JotunBack/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func StartAC(acState *model.ACState, hub *server.Hub, userRepo *repository.UserRepository,
	bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	acConfig := acState.Config
	acConfig.Degrees = int(acState.TargetTemp)

	err := userRepo.UpdateACState(acConfig)
	if err != nil {
		return
	}

	acConfig.Power = true
	err = hub.SendACConfig(acConfig)
	if err != nil {
		return
	}

	go HandleTemperature(acState, hub, userRepo, bot, update)
}

func HandleTemperature(acState *model.ACState, hub *server.Hub, userRepo *repository.UserRepository,
	bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	CoolMode := 1
	HeatMode := 2
	newTemp := acState.TargetTemp
	nextCheck := time.Now().Add(1 * time.Minute)
	nextUpdate := time.Now()
	tempData, err := userRepo.GetTemp(acState.Username, 2)
	if err != nil {
		return
	}
	currentTemp := acState.TargetTemp
	prevTemp := acState.TargetTemp

	currentTemp, prevTemp = 0, 0

	chatID := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID
	acState.Stop = false

	for !acState.Stop {
		if time.Now().After(nextCheck) || time.Now().After(nextUpdate) {
			tempData, err = userRepo.GetTemp(acState.Username, 2)
			if err != nil {
				continue
			}

			if len(tempData) > 1 {
				currentTemp = tempData[0].Temperature
				prevTemp = tempData[1].Temperature
			}

			log.Println(tempData)
		}

		if len(tempData) > 1 && time.Now().After(nextCheck) {
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
			acConfig := acState.Config
			acConfig.Degrees = int(newTemp)
			err = hub.SendACConfig(acConfig)
			if err != nil {
				continue
			}

			log.Println("Temperature:", prevTemp, "->", currentTemp, "->", acState.TargetTemp)
			log.Println("Mode:", acState.Config.Mode)
			nextCheck = time.Now().Add(1 * time.Minute)
		}

		if time.Now().After(nextUpdate) {
			isOnline := hub.GetConnectionByID(acState.Username) != nil
			msg := tgbotapi.NewEditMessageText(chatID, msgID,
				ui.StateMsg(acState, isOnline))

			keyboard := ui.StateKeyboard(acState)
			msg.ReplyMarkup = &keyboard

			editedMsg, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			msgID = editedMsg.MessageID
			nextUpdate = time.Now().Add(5 * time.Second)
			log.Println("Temperature:", prevTemp, "->", currentTemp, "->", acState.TargetTemp)
		}

	}
}
