package main

import (
	//"game/helpers"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RoomMulti struct {
	ID         string
	MaxPlayers uint
	Players    map[string]*Player
	mu         sync.Mutex
	register   chan *Player
	unregister chan *Player
	ticker     *time.Ticker
	state      *RoomState

	broadcast chan *OutcomeMessage
	finish chan *Player
}

func NewRoomMulti(MaxPlayers uint) *RoomMulti {
	return &RoomMulti{
		MaxPlayers: MaxPlayers,
		Players:    make(map[string]*Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(100 * time.Millisecond),
		state: &RoomState{
			Penguin: new(PenguinState),
			Gun: new(GunState),
			Fishes: make(map[int]*FishState, 24),
		},
		broadcast: make(chan *OutcomeMessage),
		finish: make(chan *Player),
	}
}

func (r *RoomMulti) Run() {
	LogMsg("Room Multi loop started")
	//r.state.Gun.Bullet = CreateBullet(r)
	//GameInit(r)
	for {
		select {
		case player := <-r.unregister:
			delete(r.Players, player.ID)
			LogMsg("Player " + player.ID + " was removed from room")
		case player := <-r.register:
			r.mu.Lock()
			r.Players[player.ID] = player
			r.mu.Unlock()
			LogMsg("Player " + player.ID + " joined")
			//r.broadcast <- &OutcomeMessage{Type:START}
		case message := <- r.broadcast:
			fmt.Println("IN BROADCAST")
			for _, player := range r.Players {
				select {
				case player.out <- message:
				default:
					close(player.out)
				}
			}
			//HandleCommand(r, message)
		case <-r.ticker.C:
			//ProcessGameMulti(r)
		case player := <- r.finish:
			LogMsg("Player " + player.ID + " finished game")
			player.out <- &OutcomeMessage{Type:FINISH}
			r.state.Penguin = nil
			r.state.Gun = nil
			//FinishGame(r)
		}
	}
}

func (r *RoomMulti) AddPlayer(player *Player) {
	ps := &PenguinState{
		ID:                 player.ID,
		Alpha:              0,
		ClockwiseDirection: true,
		Score:				0,
	}
	r.mu.Lock()
	r.state.Penguin = ps
	r.mu.Unlock()
	player.roomMulti = r
	r.register <- player
}

func (r *RoomMulti) RemovePlayer(player *Player) {
	r.unregister <- player
}

func (r *RoomMulti) SelectPlayersRoles() {
	//TODO make it random
	for _, player := range r.Players {
		digit := rand.Intn(1)
		if digit == 0 {
			player.Type = "penguin"
		}

		if digit != 0 {
			player.Type = "gun"
		}
	}
}