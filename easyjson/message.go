package easyjson


//from back to front
//(это я генерю и шлю)

//easyjson:json
type OutcomeMessage struct {
	Type    ServerMessage `json:"type"`
	Payload OutPayloadMessage  `json:"payload"`
}

//easyjson:json
type OutPayloadMessage struct {
	Penguin PenguinMessage `json:"penguin"`
	Gun GunMessage `json:"gun"`
	PiscesCount uint `json:"PiscesCount"`
	Round uint `json:"round"`
}

//easyjson:json
type PenguinMessage struct {
	Name      string          `json:"name"`
	Clockwise bool            `json:"clockwise"`
	Alpha     int             `json:"alpha"`
	Result    GameResult `json:"result"`
	Score     uint            `json:"score"`
}

//easyjson:json
type GunMessage struct {
	Name      string          `json:"name"`
	Clockwise bool            `json:"clockwise"`
	Alpha     int             `json:"alpha"`
	Result    GameResult `json:"result"`
	Score     uint            `json:"score"`
	Bullet    BulletMessage   `json:"bullet"`
}

//easyjson:json
type BulletMessage struct {
	DistanceFromCenter int `json:"distance_from_center"`
	Alpha int `json:"alpha"`
}

//from front to back
//(это я ТОЛЬКО парсю и никогда не шлю)

//easyjson:json
type IncomeMessage struct {
	Type    ClientCommand `json:"type"`
	Payload InPayloadMessage   `json:"payload"`
}

//easyjson:json
type InPayloadMessage struct {
	Name string        `json:"name"`
	Mode GameMode `json:"mode"`
}
