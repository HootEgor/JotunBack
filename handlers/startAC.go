package handlers

import (
	"JotunBack/model"
	"JotunBack/repository"
	"JotunBack/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func StartAC(acState *model.ACState, hub *server.Hub, userRepo *repository.UserRepository,
	update tgbotapi.Update) {
	acState.Update = update
	acState.Config.Degrees = int(acState.TargetTemp)
	acState.Config.Power = true

	err := userRepo.UpdateACState(acState.Config)
	if err != nil {
		return
	}

	err = hub.SendACConfig(acState.Config)
	if err != nil {
		return
	}
	acState.Stop = false
	acState.NextCheck = acState.NextCheck.Add(1 * time.Minute)
}
