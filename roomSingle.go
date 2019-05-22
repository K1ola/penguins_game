package main

import (
	"fmt"
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
		ticker:     time.NewTicker(10 * time.Millisecond),
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
					case FINISHROUND:
						fmt.Println(FINISHROUND)
						fmt.Println(r.gameState)
						r.StartNewRound()
					case FINISHGAME:
						//message = r.FinishGame()

						r.gameState = FINISHED
						r.SaveResult()
					}
				}
				r.Player.out <- message
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

func (r *RoomSingle) FinishRound() {
	r.round++
	helpers.LogMsg("Player " + r.Player.ID + " finished round")
	r.gameState = WAITING
}

func (r *RoomSingle) FinishGame() {
	helpers.LogMsg("Player " + r.Player.ID + " finished round")
	r.gameState = FINISHED
}

func (r *RoomSingle) StartNewRound() {
	//time.Sleep(500 * time.Millisecond)
		message := &OutcomeMessage{
			Type: START,
			Payload: OutPayloadMessage{
				Gun: GunMessage{
					Name: string(GUN),
					Score: uint(r.state.Gun.Score),
				},
				Penguin: PenguinMessage{
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
		r.gameState = RUNNING
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

