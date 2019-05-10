package main

import "time"

type PenguinState struct {
	ID                 string
	ClockwiseDirection bool
	Alpha              int
	Result 			   string
	Score              int
}

type GunState struct {
	ID                 string
	ClockwiseDirection bool
	Alpha              int
	Bullet 			   *BulletState
	Result 			   string
	Score              int
}

type BulletState struct {
	Alpha int
	DistanceFromCenter int
}

type FishState struct {
	Alpha int
	Eaten bool
}

type RoomState struct {
	Penguin *PenguinState
	Gun  *GunState
	Fishes 	map[int]*FishState
	CurrentTime time.Time
}
