package main

import (
	"fmt"
	"game/easyjson"
	"game/helpers"
	"game/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

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
	state      *easyjson.RoomState
	gameState  easyjson.GameCurrentState
	round      int

	broadcast chan *easyjson.OutcomeMessage
	finish chan *Player
}

func NewRoomSingle(MaxPlayers uint, id int) *RoomSingle {
	return &RoomSingle{
		ID: id,
		MaxPlayers: MaxPlayers,
		Player:    new(Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(10 * time.Millisecond),
		state: &easyjson.RoomState{
			Penguin: new(easyjson.PenguinState),
			Gun: new(easyjson.GunState),
			Fishes: make(map[int]*easyjson.FishState, 24),
			Round: 0,
		},
		round: 0,
		broadcast: make(chan *easyjson.OutcomeMessage, 1),
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
			if r.gameState == easyjson.RUNNING {
				message := RunSingle(r)
				if message.Type != easyjson.STATE {
					switch message.Type {
					case easyjson.FINISHROUND:
						fmt.Println(easyjson.FINISHROUND)
						fmt.Println(r.gameState)
						r.StartNewRound()
					case easyjson.FINISHGAME:
						//message = r.FinishGame()

						r.gameState = easyjson.FINISHED
						r.SaveResult()
					}
				}
				r.Player.out <- message
			}
		}
	}
}

func (r *RoomSingle) AddPlayer(player *Player) {
	ps := &easyjson.PenguinState{
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


func (r *RoomSingle) ProcessCommand(message *easyjson.IncomeMessage) {
	r.state.RotatePenguin()
}

func (r *RoomSingle) FinishRound() {
	r.round++
	helpers.LogMsg("Player " + r.Player.ID + " finished round")
	r.gameState = easyjson.WAITING
}

func (r *RoomSingle) FinishGame() {
	helpers.LogMsg("Player " + r.Player.ID + " finished round")
	r.gameState = easyjson.FINISHED
}

func (r *RoomSingle) StartNewRound() {
	//time.Sleep(500 * time.Millisecond)
		message := &easyjson.OutcomeMessage{
			Type: easyjson.START,
			Payload: easyjson.OutPayloadMessage{
				Gun: easyjson.GunMessage{
					Name: string(easyjson.GUN),
					Score: uint(r.state.Gun.Score),
				},
				Penguin: easyjson.PenguinMessage{
					//Name: penguin,
					Name: r.state.Penguin.ID,
					Score: uint(r.state.Penguin.Score),
				},
				PiscesCount: 24,
				Round:       uint(r.round),
			},
		}
		r.Player.out <- message
		r.state = CreateInitialStateSingle(r)
		r.gameState = easyjson.RUNNING
}

func (r *RoomSingle) SaveResult() {
	//TODO do it correctly and once
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		helpers.LogMsg("Can`t connect to grpc")
		return
	}
	defer grcpConn.Close()

	AuthManager = models.NewAuthCheckerClient(grcpConn)
	r.Player.instance.Score = uint64(r.Player.roomSingle.state.Penguin.Score)
	fmt.Println(r.Player.Type)
	ctx := context.Background()
	_, err = AuthManager.SaveUserGame(ctx, r.Player.instance)
	fmt.Println(err)

}

