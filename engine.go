package main

import (
	"math/rand"
)

func CreatePenguin(id string) *PenguinState {
	return &PenguinState{
		ID:                 id,
		Result:             "",
		Alpha:              rand.Intn(360),
		Score:              0,
		ClockwiseDirection: true,
	}
}

func CreateGun(id string) *GunState {
	return &GunState{
		ID:                 id,
		Result:             "",
		Alpha:              rand.Intn(360),
		Score:              0,
		ClockwiseDirection: true,
		Bullet:             CreateBullet(),
	}
}

func CreateBullet() *BulletState {
	return &BulletState{
		Alpha:              rand.Intn(360),
		DistanceFromCenter: 0,
	}
}

func CreateFishes() map[int]*FishState {
	fishes := make(map[int]*FishState, 24)
	for i := 0; i < 24; i++ {
		fishes[i] = &FishState{Eaten: false, Alpha: 360 / 24 * i}
	}
	return fishes
}

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
	if rs.Gun.Alpha >= 360 {
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

func (rs *RoomState) RecalcBullet() *OutcomeMessage {
	if rs.Gun.Bullet.DistanceFromCenter > 100*0.8/2 {
		if rs.Gun.Bullet.Alpha%360 >= rs.Penguin.Alpha-7 && rs.Gun.Bullet.Alpha%360 <= rs.Penguin.Alpha+7 {
			//lost
			return &OutcomeMessage{
				Type: FINISH,
				Payload: OutPayloadMessage{
					Penguin: PenguinMessage{
						Name:   rs.Penguin.ID,
						Score:  uint(rs.Penguin.Score),
						Result: LOST,
					},
					Gun: GunMessage{
						Name:   rs.Gun.ID,
						Score:  uint(rs.Gun.Score),
						Result: WIN,
					},
				}}
		}

		rs.Gun.Bullet.Alpha = rs.Gun.Alpha

		rs.Gun.Bullet.DistanceFromCenter = 0
	}
	rs.Gun.Bullet.DistanceFromCenter += 5
	return nil
}

func (rs *RoomState) RecalcPenguin() *OutcomeMessage {
	if rs.Penguin.Alpha == 360 {
		rs.Penguin.Alpha = 0
	}

	if rs.Penguin.Alpha == -1 {
		rs.Penguin.Alpha = 359
	}

	for i := 0; i < len(rs.Fishes); i++ {
		if rs.Penguin.Alpha == rs.Fishes[i].Alpha {
			rs.Penguin.Score++

			rs.Fishes[i].Eaten = true
			break
		}
	}

	count := 0
	for i := 0; i < len(rs.Fishes); i++ {
		if rs.Fishes[i].Eaten == false {
			count++
		}
	}

	if count == 0 {
		// win penguin
		return &OutcomeMessage{
			Type: FINISH,
			Payload: OutPayloadMessage{
				Penguin: PenguinMessage{
					Name:   rs.Penguin.ID,
					Score:  uint(rs.Penguin.Score),
					Result: WIN,
				},
				Gun: GunMessage{
					Name:   rs.Gun.ID,
					Score:  uint(rs.Gun.Score),
					Result: LOST,
				},
			}}
	}

	if rs.Penguin.ClockwiseDirection {
		rs.Penguin.Alpha++
	} else {
		rs.Penguin.Alpha--
	}
	return nil
}

func (rs *RoomState) GetState() *OutcomeMessage {
	return &OutcomeMessage{
		Type: STATE,
		Payload: OutPayloadMessage{
			Penguin: PenguinMessage{
				Alpha:     rs.Penguin.Alpha,
				Score:     uint(rs.Penguin.Score),
				Result:    rs.Penguin.Result,
				Name:      rs.Penguin.ID,
				Clockwise: rs.Penguin.ClockwiseDirection,
			},
			Gun: GunMessage{
				Name:   rs.Gun.ID,
				Result: rs.Gun.Result,
				Score:  uint(rs.Gun.Score),
				Alpha:  rs.Gun.Alpha,
				Bullet: BulletMessage{
					Alpha:              rs.Gun.Bullet.Alpha,
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
