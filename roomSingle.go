package main

import (
	"game/helpers"
	//"game/helpers"
	"sync"
	"time"
)

type RoomSingle struct {
	ID         int
	MaxPlayers uint
	Player     *Player
	mu         sync.Mutex
	register   chan *Player
	unregister chan *Player
	ticker     *time.Ticker
	state      *RoomState
	gameState GameCurrentState

	broadcast chan *OutcomeMessage
	finish chan *Player
}

func NewRoomSingle(MaxPlayers uint, id int) *RoomSingle {
	return &RoomSingle{
		ID: id,
		MaxPlayers: MaxPlayers,
		Player:    new(Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(1 * time.Second),
		state: &RoomState{
			Penguin: new(PenguinState),
			Gun: new(GunState),
			Fishes: make(map[int]*FishState, 24),
		},
		broadcast: make(chan *OutcomeMessage, 1),
		finish: make(chan *Player),
	}
}

func (r *RoomSingle) Run() {
	//defer helpers.RecoverPanic()
	helpers.LogMsg("Room Single loop started")
	for {
		select {
		case player := <-r.unregister:
			r.Player = nil
			helpers.LogMsg("Player " + player.ID + " was removed from room")
		case player := <-r.register:
			r.mu.Lock()
			r.Player = player
			r.mu.Unlock()
			helpers.LogMsg("Player " + player.ID + " joined")
			//r.Player.out <- &OutcomeMessage{Type:START}
		case <-r.ticker.C:
			if r.gameState == RUNNING {
				r.Player.out <- RunSingle(r)
			}
		}
	}
}

func (r *RoomSingle) AddPlayer(player *Player) {
	ps := &PenguinState{
		ID:                 player.ID,
		Alpha:              0,
		ClockwiseDirection: true,
		Score:				0,
	}
	r.mu.Lock()
	r.state.Penguin = ps
	r.mu.Unlock()
	player.roomSingle = r
	r.register <- player
}

func (r *RoomSingle) RemovePlayer(player *Player) {
	//r.unregister <- player
	r.Player = nil
	helpers.LogMsg("Player " + player.ID + " was removed from room")
}


func (r *RoomSingle) ProcessCommand(message *IncomeMessage) {
	r.state.RotatePenguin()
}
