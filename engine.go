package main

import (
	"game/easyjson"
	"math/rand"
)

func CreatePenguin(id string) *easyjson.PenguinState {
	return &easyjson.PenguinState{
		ID: id,
		Result: "",
		Alpha: rand.Intn(360),
		//Score: 0,
		ClockwiseDirection: true,
	}
}

func CreateGun(id string) *easyjson.GunState {
	return &easyjson.GunState{
		ID: id,
		Result: "",
		Alpha: rand.Intn(360),
		//Score: 0,
		ClockwiseDirection: true,
		Bullet: CreateBullet(),
	}
}

func CreateBullet() *easyjson.BulletState {
	return &easyjson.BulletState{
		Alpha: rand.Intn(360),
		DistanceFromCenter: 0,
	}
}

func CreateFishes() map[int]*easyjson.FishState {
	fishes := make(map[int]*easyjson.FishState, 24)
	for i := 0; i < 24; i++ {
		fishes[i] = &easyjson.FishState{Eaten: false, Alpha: 360/24*i}
	}
	return fishes
}

func RunMulti(room *RoomMulti) *easyjson.OutcomeMessage {
	msg := room.state.RecalcPenguin()
	if msg != nil {
		room.FinishRound()
		return msg
	}
	room.state.RecalcGun()
	msg = room.state.RecalcBullet()
	if msg != nil {
		room.FinishRound()
		return msg
	}
	return room.state.GetState()
}

func RunSingle(room *RoomSingle) *easyjson.OutcomeMessage {
	msg := room.state.RecalcPenguin()
	if msg != nil {
		room.FinishRound()
		return msg
	}
	room.state.RecalcGun()
	msg = room.state.RecalcBullet()
	if msg != nil {
		//room.FinishGame()
		return msg
	}
	return room.state.GetState()
}

//TODO remove repeat
func CreateInitialStateSingle(room *RoomSingle) *easyjson.RoomState {
	state := new(easyjson.RoomState)
	var penguin, gun string

	state.Penguin = CreatePenguin(penguin)
	state.Gun = CreateGun(gun)
	state.Fishes = CreateFishes()
	state.Round = room.round
	var penguinScore, gunScore int
	if room.state != nil {
		penguinScore = room.state.Penguin.Score
		gunScore = room.state.Gun.Score
	}
	room.state = state
	room.state.Penguin.Score = penguinScore
	room.state.Gun.Score = gunScore
	room.state.Gun.ID = string(easyjson.GUN)
	return state
}

func CreateInitialState(room *RoomMulti) *easyjson.RoomState {
	state := new(easyjson.RoomState)
	var penguin, gun string
	for _, player := range room.Players {
		if player.Type == easyjson.PENGUIN {
			penguin = player.ID
		} else {
			gun = player.ID
		}
	}
	state.Penguin = CreatePenguin(penguin)
	state.Gun = CreateGun(gun)
	state.Fishes = CreateFishes()
	state.Round = room.round
	var penguinScore, gunScore int
	if room.state != nil {
		penguinScore = room.state.Penguin.Score
		gunScore = room.state.Gun.Score
	}
	room.state = state
	room.state.Penguin.Score = penguinScore
	room.state.Gun.Score = gunScore
	return state
}

