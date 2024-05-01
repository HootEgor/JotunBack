package model

import "time"

type Temp struct {
	Temperature float32 `json:"temperature"`
}

type TempDB struct {
	Temperature float32
	TimeStamp   time.Time
}
