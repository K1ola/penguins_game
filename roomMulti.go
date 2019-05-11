package main

import (
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
	gameState string

	broadcast chan *OutcomeMessage
	finish chan *Player
}

func NewRoomMulti(MaxPlayers uint) *RoomMulti {
	return &RoomMulti{
		MaxPlayers: MaxPlayers,
		Players:    make(map[string]*Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(1 * time.Second),
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
			if len(r.Players) == 2 {
				penguin, gun := r.SelectPlayersRoles()
				message := &OutcomeMessage{
					Type: START,
					Payload: OutPayloadMessage{
						Gun: GunMessage{
							Name: gun,
						},
						Penguin: PenguinMessage{
							Name: penguin,
						},
						PiscesCount: 24,
					},
				}
				r.gameState = START
				r.SendRoomState(message)
				r.state = CreateInitialState(r)

			}
		case message := <- r.broadcast:
			r.SendRoomState(message)
		case <-r.ticker.C:
			if r.gameState == START {
				  message := RunMulti(r)
				  if message.Type != STATE {
				  	r.gameState = FINISH
				  }
				  r.SendRoomState(message)
			}
		//case player := <- r.finish:
		//	LogMsg("Player " + player.ID + " finished game")
		//	r.state.Penguin = nil
		//	r.state.Gun = nil
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

func (r *RoomMulti) SelectPlayersRoles() (string, string) {
	var penguin, gun string
	digit := rand.Intn(2)
	for _, player := range r.Players {
		if digit == 0 {
			player.Type = PENGUIN
			penguin = player.ID
			digit = 1
		} else {
			player.Type = GUN
			gun = player.ID
			digit = 0
		}
	}
	return penguin, gun
}

func (r *RoomMulti) ProcessCommand(message *IncomeMessage) {
	login := message.Payload.Name
	for _, player := range r.Players {
		if player.ID != login {
			continue
		}
		fmt.Println(r.state)
		switch player.Type {
		case PENGUIN:
			r.state.RotatePenguin()
		case GUN:
			r.state.RotateGun()
		default:
			fmt.Println("Incorrect player type!")
		}
		break
	}
}

func (r *RoomMulti) FinishGame(player *Player) {
	//r.finish <- player
	LogMsg("Player " + player.ID + " finished game")
	r.state.Penguin = nil
	r.state.Gun = nil
}

func (r *RoomMulti) SendRoomState(message *OutcomeMessage) {
	for _, player := range r.Players {
		select {
		case player.out <- message:
		default:
			close(player.out)
		}
	}
}

