package main

//from back to front
//(это я генерю и шлю)
type OutcomeMessage struct {
	Type string `json:"type"`
	Payload OutPayloadMessage `json:"payload"`
}

type OutPayloadMessage struct {
	Penguin PenguinMessage `json:"penguin"`
	Gun GunMessage `json:"gun"`
	PiscesCount uint `json:"PiscesCount"`
}

type PenguinMessage struct {
	Name string `json:"name"`
	Clockwise bool `json:"clockwise"`
	Alpha int `json:"alpha"`
	Result string `json:"result"`
	Score uint `json:"score"`
}

type GunMessage struct {
	Name string `json:"name"`
	Bullet BulletMessage `json:"bullet"`
	Alpha int `json:"alpha"`
	Result string `json:"result"`
	Score uint `json:"score"`
}

type BulletMessage struct {
	DistanceFromCenter int `json:"distance_from_center"`
	Alpha int `json:"alpha"`
}

//from front to back
//(это я ТОЛЬКО парсю и никогда не шлю)
type IncomeMessage struct {
	Type string `json:"type"`
	Payload InPayloadMessage `json:"payload"`
}

type InPayloadMessage struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}
