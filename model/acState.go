package model

import (
	"strconv"
	"time"
)

type ACState struct {
	Username    string
	ChatID      int64
	Until       time.Time
	TargetTemp  float32
	CurrentTemp float32
	Stop        bool
	Config      AirConditionerConfig
}

func (acState *ACState) GetTargetTemp() string {
	return strconv.FormatFloat(float64(acState.TargetTemp), 'f', 1, 32)
}

func (acState *ACState) GetTemp() string {
	return strconv.FormatFloat(float64(acState.CurrentTemp), 'f', 1, 32)
}
