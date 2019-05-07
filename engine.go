package main
//
//import (
//	"math"
//)
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
//
//func HandleCommand(r *Room, msg *Message) {
//		switch msg.Payload.Command {
//			case "SHOT":
//				ShotPlayer(r.state.Players[msg.Type], r.state.Bullet)
//				for t, player := range r.state.Players {
//					if player.Shoted {
//						r.Players[t].out <- &Message{"SINGLE", PayloadMessage{msg.Type, "KILLED"}}
//						message := &Message{"SINGLE", PayloadMessage{player.ID, "GAME FINISHED"}}
//						r.Players[t].SendMessage(message)
//					}
//				}
//			case "ROTATE":
//				RotatePlayer(r.state.Players[msg.Type])
//				r.Players[msg.Type].out <- &Message{"SINGLE", PayloadMessage{msg.Type, "ROTATED"}}
//		}
//}
//
//func FinishGame(r *Room) {
//	for _, player := range r.Players {
//		message := &Message{"SINGLE", PayloadMessage{player.ID, "GAME FINISHED"}}
//		player.SendMessage(message)
//		}
//}
