package main

import "time"

type PenguinState struct {
	ID                 string
	ClockwiseDirection bool
	Alpha              int
	Score              int
}

type GunState struct {
	ID                 string
	Alpha              int
	Bullet 			   *BulletState
}

type BulletState struct {
	Alpha int
	DistanceFromCenter int
}

type FishState struct {
	//ID int
	//X, Y int
	Alpha int
	Eaten bool
}

type RoomState struct {
	Penguin *PenguinState
	Gun  *GunState
	Fishes 	map[int]*FishState
	CurrentTime time.Time
}
