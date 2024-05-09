package handlers

import (
	"JotunBack/model"
	"errors"
	"log"
	"time"
)

func GPTRespToACState(gptResp model.GPTResponse, acState *model.ACState) error {
	if !gptResp.HaveData {
		return errors.New("I can't understand your request. Please try again.")
	}
	acState.Config.Mode = gptResp.Mode
	if gptResp.TargetTemp > 0 {
		acState.TargetTemp = gptResp.TargetTemp
	}
	extendTime, err := time.Parse("15:04:05", gptResp.Extend)
	if err != nil {
		return err
	}
	log.Println(extendTime)
	hours := extendTime.Hour()
	minutes := extendTime.Minute()
	seconds := extendTime.Second()
	acState.Until = time.Now().Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes) + time.Second*time.Duration(seconds))
	return nil
}
