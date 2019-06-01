package main

import (
	"fmt"
	"game/helpers"
	"game/models"
	"golang.org/x/net/context"
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
	round int

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
		ticker:     time.NewTicker(time.Duration(SingleGameSpeed) * time.Millisecond),
		state: &RoomState{
			Penguin: new(PenguinState),
			Gun: new(GunState),
			Fishes: make(map[int]*FishState, 24),
			Round: 1,
		},
		round: 1,
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
				message := RunSingle(r)
				if message.Type != STATE {
					switch message.Type {
					case FINISHGAME:
						r.gameState = FINISHED
						r.SaveResult()
						//r.unregister <- r.Player
						r.SendRoomState(message)
						r.Player.game.unregister <- r.Player
						return
					}
				}
				r.SendRoomState(message)
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
	r.Player = nil
	helpers.LogMsg("Player " + player.ID + " was removed from room")
}


func (r *RoomSingle) ProcessCommand(message *IncomeMessage) {
	r.state.RotatePenguin()
}

func (r *RoomSingle) FinishRound() {
	r.round++
	helpers.LogMsg("Player " + r.Player.ID + " finished round")
	r.gameState = WAITING
}

func (r *RoomSingle) FinishGame() {
	helpers.LogMsg("Player " + r.Player.ID + " finished game")
	r.gameState = FINISHED
}

func (r *RoomSingle) StartNewRound() {
	//time.Sleep(50 * time.Millisecond)
		message := &OutcomeMessage{
			Type: START,
			Payload: OutPayloadMessage{
				Gun: GunMessage{
					Name: string(GUN),
					Score: uint(r.state.Gun.Score),
				},
				Penguin: PenguinMessage{
					Name: r.state.Penguin.ID,
					Score: uint(r.state.Penguin.Score),
				},
				PiscesCount: 24,
				Round:       uint(r.round),
			},
		}
		r.SendRoomState(message)
		r.state = CreateInitialStateSingle(r)
		r.gameState = RUNNING
}


func (r *RoomSingle) SaveResult() {
	r.Player.instance.Score = uint64(r.Player.roomSingle.state.Penguin.Score)
	fmt.Println(r.Player.Type)
	ctx := context.Background()
	_, err := models.AuthManager.SaveUserGame(ctx, r.Player.instance)
	_, err = models.AuthManager.DeleteUserFromGame(ctx, r.Player.instance)
	fmt.Println(err)
}

func (r *RoomSingle) SendRoomState(message *OutcomeMessage) {
	r.Player.out <- message
}




