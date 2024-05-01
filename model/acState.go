package model

import "time"

type ACState struct {
	Username        string
	TargetTemp      float32
	CurrentTemp     float32
	Temp2minAgo     float32
	Temp2minAgoTime time.Time
	Config          AirConditionerConfig
}
