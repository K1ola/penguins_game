package main

import (
	"math/rand"
)

func CreatePenguin(id string) *PenguinState {
	return &PenguinState{
		ID: id,
		Result: "",
		Alpha: rand.Intn(360),
		Score: 0,
		ClockwiseDirection: true,
	}
}

func CreateGun(id string) *GunState {
	return &GunState{
		ID: id,
		Result: "",
		Alpha: rand.Intn(360),
		Score: 0,
		ClockwiseDirection: true,
		Bullet: CreateBullet(),
	}
}

func CreateBullet() *BulletState {
	return &BulletState{
		Alpha: rand.Intn(360),
		DistanceFromCenter: 0,
	}
}

func CreateFishes() map[int]*FishState {
	fishes := make(map[int]*FishState, 24)
	for i := 0; i < 24; i++ {
		fishes[i] = &FishState{Eaten: false, Alpha: 360/24*i}
	}
	return fishes
}

//
//const ownRandom = 0.25
//
//func RotatePlayer(ps *PlayerState) {
//	if ps.ClockwiseDirection {
//		ps.ClockwiseDirection = false
//	} else {
//		ps.ClockwiseDirection = true
//	}
//}
//
//func ShotPlayer(ps *PlayerState, b *BulletState) {
//	if ps.Shoted {
//		return
//	}
//
//	RecountBullet(ps, b)
//
//	if ps.X == b.X && ps.Y == b.Y {
//		//you loose
//		ps.Shoted = true
//		return
//	}
//}
//
//func RecountBullet(ps *PlayerState, b *BulletState) {
//	if ps.ClockwiseDirection {
//		//TODO: create own random
//		b.Alpha = ps.Alpha + ownRandom*100
//	} else {
//		b.Alpha = ps.Alpha - ownRandom*100
//	}
//
//	b.Radious = 20
//
//	b.X = int(math.Floor(250 + math.Sin(b.Alpha*(math.Pi/180))*float64(b.Radious)))
//	b.Y = int(math.Floor(250 - math.Cos(b.Alpha*(math.Pi/180))*float64(b.Radious)))
//	b.Radious += 5
//}
//
//func CreateBullet(r *Room) *BulletState {
//	return &BulletState{
//		ID:    r.ID,
//		X:     0,
//		Y:     0,
//		Alpha: 0,
//	}
//}
//
////recount coordinates (and maybe you win)
//func ProcessGameSingle(r *Room) {
//	var penguin *PlayerState
//	//var gun *PlayerState
//	for _, player := range r.state.Players {
//		if player.Type == "GOOD" {
//			penguin = player
//		} else {
//			//gun = player
//		}
//	}
//
//	if penguin.Alpha == 360 {
//		penguin.Alpha = 0
//	}
//
//	if penguin.Alpha == -1 {
//		penguin.Alpha = 359
//	}
//
//	fishCount := 24
//	for i := 0; i < fishCount; i++ {
//		if penguin.Alpha == r.state.Fishes[i].Alpha {
//			penguin.Score ++
//
//			r.state.Fishes[i].Eaten = true
//			break
//		}
//	}
//
//	count := 0
//	for i := 0; i < fishCount; i++ {
//		if r.state.Fishes[i].Eaten == false {
//			count ++
//		}
//	}
//
//	if count == 0 {
//		for t, _ := range r.state.Players {
//			r.Players[t].out <- &Message{"SINGLE", PayloadMessage{penguin.ID, "WIN"}}
//
//			message := &Message{"SINGLE", PayloadMessage{r.Players[t].ID, "GAME FINISHED"}}
//			r.Players[t].SendMessage(message)
//		}
//
//		r.finish <- &Message{"SINGLE", PayloadMessage{penguin.ID, "WIN"}}
//
//		return
//	}
//
//	if penguin.ClockwiseDirection {
//		penguin.Alpha ++
//	} else {
//		penguin.Alpha --
//	}
//
//	alphaRad := degreesToRadians(penguin.Alpha)
//
//	penguin.X = int(math.Floor(float64(r.state.Radious) + math.Sin(alphaRad)*float64(r.state.Radious)))
//	penguin.Y = int(math.Floor(float64(r.state.Radious) - math.Cos(alphaRad)*float64(r.state.Radious)))
//}
//
//func degreesToRadians(degrees float64) float64 {
//	return degrees * (math.Pi / 180)
//}
//
////prepare room for game
//func GameInit(r *Room) {
//	fishCount := 24
//	for i := 0; i < fishCount; i++ {
//		Alpha := float64((360/24)*i)
//		alphaRad := degreesToRadians(Alpha)
//		X := int(math.Floor(float64(r.state.Radious) + math.Sin(alphaRad)*float64(r.state.Radious)))
//		Y := int(math.Floor(float64(r.state.Radious) - math.Cos(alphaRad)*float64(r.state.Radious)))
//		r.state.Fishes[i] = &FishState{i, X,Y,Alpha, false}
//	}
//
//	for _, player := range r.state.Players {
//		if player.Type == "GOOD" {
//			player.Alpha = math.Floor(ownRandom*360)
//			alphaRad := degreesToRadians(player.Alpha)
//			player.X = int(math.Floor(float64(r.state.Radious) + math.Sin(alphaRad)*float64(r.state.Radious)))
//			player.Y = int(math.Floor(float64(r.state.Radious) - math.Cos(alphaRad)*float64(r.state.Radious)))
//		}
//	}
//}

//func FinishGame(r *Room) {
//}

func RunMulti(room *RoomMulti) *OutcomeMessage {
	msg := room.state.RecalcPenguin()
	if msg != nil {
		for _, player := range room.Players {
			//TODO
			room.FinishGame(player)
		}
		return msg
	}
	room.state.RecalcGun()
	msg = room.state.RecalcBullet()
	if msg != nil {
		for _, player := range room.Players {
			//room.FinishGame(player)
			player.Finish()
		}
		return msg
	}
	return room.state.GetState()
}

func CreateInitialState(room *RoomMulti) *RoomState {
	state := new(RoomState)
	var penguin, gun string
	for _, player := range room.Players {
		if player.Type == PENGUIN {
			penguin = player.ID
		} else {
			gun = player.ID
		}
	}
	state.Penguin = CreatePenguin(penguin)
	state.Gun = CreateGun(gun)
	state.Fishes = CreateFishes()
	room.state = state
	return state
}

func (rs *RoomState) RecalcGun() {
	if rs.Gun.Alpha == 360 {
		rs.Gun.Alpha = 0
	}

	if rs.Gun.Alpha <= -1 {
		rs.Gun.Alpha = 359
	}

	if rs.Gun.ClockwiseDirection {
		rs.Gun.Alpha += 3
	} else {
		rs.Gun.Alpha -= 3
	}
	//return nil
}

func (rs *RoomState) RecalcBullet() *OutcomeMessage{
	if rs.Gun.Bullet.DistanceFromCenter > 100*0.8/2 {
		if rs.Gun.Bullet.Alpha % 360 >= rs.Penguin.Alpha - 7 && rs.Gun.Bullet.Alpha % 360 <= rs.Penguin.Alpha + 7 {
			//lost
			return &OutcomeMessage{
				Type:FINISH,
				Payload:OutPayloadMessage{
					Penguin:PenguinMessage{
						Name: rs.Penguin.ID,
						Score: uint(rs.Penguin.Score),
						Result:LOST,
					},
					Gun:GunMessage{
						Name: rs.Gun.ID,
						Score: uint(rs.Gun.Score),
						Result:WIN,
					},
				}}
		}
		if rs.Penguin.ClockwiseDirection {
			rs.Gun.Bullet.Alpha = rs.Penguin.Alpha + rand.Intn(101)
		} else {
			rs.Gun.Bullet.Alpha = rs.Penguin.Alpha - rand.Intn(101)
		}
		rs.Gun.Bullet.DistanceFromCenter = 0

	}
	rs.Gun.Bullet.DistanceFromCenter += 5
	return nil
}

func (rs *RoomState) RecalcPenguin() *OutcomeMessage{
		if rs.Penguin.Alpha == 360 {
			rs.Penguin.Alpha = 0
		}

		if rs.Penguin.Alpha == -1 {
			rs.Penguin.Alpha = 359
		}

		for i := 0; i < len(rs.Fishes); i++ {
			if rs.Penguin.Alpha == rs.Fishes[i].Alpha {
				rs.Penguin.Score ++

				rs.Fishes[i].Eaten = true
				break
			}
		}

		count := 0
		for i := 0; i <  len(rs.Fishes); i++ {
			if rs.Fishes[i].Eaten == false {
				count ++
			}
		}

		if count == 0 {
			// win penguin
			return &OutcomeMessage{
				Type:FINISH,
				Payload:OutPayloadMessage{
					Penguin:PenguinMessage{
						Name: rs.Penguin.ID,
						Score: uint(rs.Penguin.Score),
						Result:WIN,
					},
					Gun:GunMessage{
						Name: rs.Gun.ID,
						Score: uint(rs.Gun.Score),
						Result:LOST,
					},
				}}
		}

		if rs.Penguin.ClockwiseDirection {
			rs.Penguin.Alpha ++
		} else {
			rs.Penguin.Alpha --
		}
	return nil
}

func (rs *RoomState) GetState() *OutcomeMessage {
	return &OutcomeMessage{
		Type:STATE,
		Payload:OutPayloadMessage{
			Penguin:PenguinMessage{
				Alpha: rs.Penguin.Alpha,
				Score: uint(rs.Penguin.Score),
				Result: rs.Penguin.Result,
				Name: rs.Penguin.ID,
				Clockwise: rs.Penguin.ClockwiseDirection,
			},
			Gun:GunMessage{
				Name: rs.Gun.ID,
				Result: rs.Gun.Result,
				Score: uint(rs.Gun.Score),
				Alpha: rs.Gun.Alpha,
				Bullet: BulletMessage{
					Alpha: rs.Gun.Bullet.Alpha,
					DistanceFromCenter: rs.Gun.Bullet.DistanceFromCenter,
				},
				Clockwise: rs.Gun.ClockwiseDirection,
			},
			PiscesCount: 24,
		},
	}
}

func (rs *RoomState) RotatePenguin() {
		if rs.Penguin.ClockwiseDirection {
			rs.Penguin.ClockwiseDirection = false
		} else {
			rs.Penguin.ClockwiseDirection = true
		}
}

func (rs *RoomState) RotateGun() {
		if rs.Gun.ClockwiseDirection {
			rs.Gun.ClockwiseDirection = false
		} else {
			rs.Gun.ClockwiseDirection = true
		}
}

