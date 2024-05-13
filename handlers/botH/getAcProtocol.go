package botH

import (
	"JotunBack/repository"
	"JotunBack/server"
	"log"
)

func GetAcProtocol(username string, userRepo *repository.UserRepository,
	hub *server.Hub) error {
	currentConfig, err := userRepo.GetACState(username)
	if err != nil {
		return err
	}

	currentConfig.Config = true

	err = hub.SendACConfig(currentConfig)
	if err != nil {
		log.Println("Error sending AC config to hub:", err)
		return err
	}

	return nil
}

func TurnOnAc(username string, userRepo *repository.UserRepository,
	hub *server.Hub) error {
	currentConfig, err := userRepo.GetACState(username)
	if err != nil {
		return err
	}

	currentConfig.Power = true

	err = hub.SendACConfig(currentConfig)
	if err != nil {
		log.Println("Error sending AC config to hub:", err)
		return err
	}

	return nil
}

func TurnOffAc(username string, userRepo *repository.UserRepository,
	hub *server.Hub) error {
	currentConfig, err := userRepo.GetACState(username)
	if err != nil {
		return err
	}

	currentConfig.Power = false
	err = userRepo.UpdateACState(currentConfig)
	if err != nil {
		return err
	}

	err = hub.SendACConfig(currentConfig)
	if err != nil {
		log.Println("Error sending AC config to hub:", err)
		return err
	}

	return nil
}
