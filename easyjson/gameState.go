package easyjson
import (
	"math/rand"
	"time"
)

//easyjson:json
type PenguinState struct {
	ID                 string
	ClockwiseDirection bool
	Alpha              int
	Result             GameResult
	Score              int
}

//easyjson:json
type GunState struct {
	ID                 string
	ClockwiseDirection bool
	Alpha              int
	Bullet             *BulletState
	Result             GameResult
	Score              int
}

//easyjson:json
type BulletState struct {
	Alpha int
	DistanceFromCenter int
}

//easyjson:json
type FishState struct {
	Alpha int
	Eaten bool
}

//easyjson:json
type RoomState struct {
	Penguin *PenguinState
	Gun  *GunState
	Fishes 	map[int]*FishState
	CurrentTime time.Time
	Round int
}

func (rs *RoomState) RecalcGun() {
	//rs.Gun.Alpha = 1000
	if rs.Gun.Alpha >= 360 {
		rs.Gun.Alpha = 0
	}

	if rs.Gun.Alpha <= -1 {
		rs.Gun.Alpha = 359
	}

	var delta int
	if rs.Gun.ID == string(GUN) {
		delta = 1
	} else {
		delta = 3
	}
	if rs.Gun.ClockwiseDirection {
		rs.Gun.Alpha += delta //3
	} else {
		rs.Gun.Alpha -= delta //3
	}
}

func (rs *RoomState) RecalcBullet() *OutcomeMessage {
	if rs.Gun.Bullet.DistanceFromCenter > 100*0.8/2 {
		if rs.Gun.Bullet.Alpha % 360 >= rs.Penguin.Alpha - 7 && rs.Gun.Bullet.Alpha % 360 <= rs.Penguin.Alpha + 7 {

			//TODO it is single mode logic
			if rs.Gun.ID == string(GUN) {
				return &OutcomeMessage{
					Type: FINISHGAME,
					Payload: OutPayloadMessage{
						Penguin: PenguinMessage{
							Name: rs.Penguin.ID,
							Score: uint(rs.Penguin.Score),
						},
						Gun: GunMessage{
							Name: rs.Gun.ID,
						},
						Round: uint(rs.Round),
					}}
			} else {
				//lost
				scoreGun := rs.Gun.Score + 1
				rs.Gun.Score = scoreGun
				return &OutcomeMessage{
					Type: FINISHROUND,
					Payload: OutPayloadMessage{
						Penguin: PenguinMessage{
							Name: rs.Penguin.ID,
							Score: uint(rs.Penguin.Score),
						},
						Gun: GunMessage{
							Name: rs.Gun.ID,
							Score: uint(scoreGun),
						},
						Round: uint(rs.Round),
					}}
			}
		}

		rs.Gun.Bullet.Alpha = rs.Gun.Alpha
		//TODO it is single mode logic
		if rs.Gun.ID == string(GUN) {
			if rs.Penguin.ClockwiseDirection {
				alpha := rs.Penguin.Alpha + rand.Intn(101)
				if alpha >= 360 {
					rs.Gun.Bullet.Alpha = alpha - 360
				} else {
					rs.Gun.Bullet.Alpha = alpha
				}
				rs.Gun.Bullet.Alpha = rs.Penguin.Alpha + rand.Intn(101)
			} else {
				alpha := rs.Penguin.Alpha - rand.Intn(101)
				if alpha < 0 {
					rs.Gun.Bullet.Alpha = 360 + alpha
				} else {
					rs.Gun.Bullet.Alpha = alpha
				}
			}
		}

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
			//rs.Penguin.Score ++

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
		//if rs.Gun.ID != string(GUN) {
		// win penguin
		scorePenguin := rs.Penguin.Score + 1
		rs.Penguin.Score = scorePenguin
		return &OutcomeMessage{
			Type: FINISHROUND,
			Payload: OutPayloadMessage{
				Penguin: PenguinMessage{
					Name: rs.Penguin.ID,
					Score: uint(scorePenguin),
				},
				Gun: GunMessage{
					Name: rs.Gun.ID,
					Score: uint(rs.Gun.Score),
				},
				Round: uint(rs.Round),
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
		Type: STATE,
		Payload: OutPayloadMessage{
			Penguin: PenguinMessage{
				Alpha: rs.Penguin.Alpha,
				Score: uint(rs.Penguin.Score),
				Result: rs.Penguin.Result,
				Name: rs.Penguin.ID,
				Clockwise: rs.Penguin.ClockwiseDirection,
			},
			Gun: GunMessage{
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
			Round: uint(rs.Round),
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

