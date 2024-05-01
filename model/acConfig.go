package model

type AirConditionerConfig struct {
	Username string `json:"username"`
	Config   bool   `json:"config"`
	Protocol int    `json:"protocol"`
	Model    int    `json:"model"`
	Mode     int    `json:"mode"`
	Celsius  bool   `json:"celsius"`
	Degrees  int    `json:"degrees"`
	FanSpeed int    `json:"fanspeed"`
	SwingV   int    `json:"swingv"`
	SwingH   int    `json:"swingh"`
	Light    bool   `json:"light"`
	Beep     bool   `json:"beep"`
	Econo    bool   `json:"econo"`
	Filter   bool   `json:"filter"`
	Turbo    bool   `json:"turbo"`
	Quiet    bool   `json:"quiet"`
	Sleep    int    `json:"sleep"`
	Clean    bool   `json:"clean"`
	Clock    int    `json:"clock"`
	Power    bool   `json:"power"`
}
