package main

import (
	"fmt"
	"game/easyjson"
	"game/helpers"
	"game/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"time"
)

type RoomMulti struct {
	ID         int
	MaxPlayers uint
	Players    map[string]*Player
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

func NewRoomMulti(MaxPlayers uint, id int) *RoomMulti {
	return &RoomMulti{
		ID: id,
		MaxPlayers: MaxPlayers,
		Players:    make(map[string]*Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(100 * time.Millisecond),
		state: &easyjson.RoomState{
			Penguin: new(easyjson.PenguinState),
			Gun: new(easyjson.GunState),
			Fishes: make(map[int]*easyjson.FishState, 24),
			Round: -1,
		},
		round: -1,
		broadcast: make(chan *easyjson.OutcomeMessage),
		finish: make(chan *Player),
	}
}

func (r *RoomMulti) Run() {
	//defer helpers.RecoverPanic()
	helpers.LogMsg("Room Multi loop started")
	for {
		select {
		case player := <-r.unregister:
			delete(r.Players, player.ID)
			helpers.LogMsg("Player " + player.ID + " was removed from room")
		case player := <-r.register:
			r.mu.Lock()
			r.Players[player.ID] = player
			r.mu.Unlock()
			helpers.LogMsg("Player " + player.ID + " joined")
		case message := <- r.broadcast:
			r.SendRoomState(message)
		case <-r.ticker.C:
			if r.gameState == easyjson.RUNNING {
				  message := RunMulti(r)
				  if message.Type != easyjson.STATE {
					  switch message.Type {
					  case easyjson.FINISHROUND:
					  		fmt.Println(easyjson.FINISHROUND)
					  		fmt.Println(r.gameState)
					  case easyjson.FINISHGAME:
							message = r.FinishGame()
					  }
				  }
				  r.SendRoomState(message)
				if r.round == 2 && r.gameState == easyjson.FINISHED {
					message := r.FinishGame()
					r.SendRoomState(message)
					//r.gameState = FINISHED
					r.SaveResult()
					return
				}
			}
		}
	}
}

func (r *RoomMulti) AddPlayer(player *Player) {
	ps := &easyjson.PenguinState{
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
	//r.unregister <- player
	delete(r.Players, player.ID)
	helpers.LogMsg("Player " + player.ID + " was removed from room")
}

func (r *RoomMulti) SelectPlayersRoles() (string, string) {
	var penguin, gun string
	digit := rand.Intn(2)
	time.Sleep(500* time.Millisecond)
	for _, player := range r.Players {
		if digit == 0 {
			player.Type = easyjson.PENGUIN
			penguin = player.ID
			digit = 1
		} else {
			player.Type = easyjson.GUN
			gun = player.ID
			digit = 0
		}
	}
	return penguin, gun
}

func (r *RoomMulti) ProcessCommand(message *easyjson.IncomeMessage) {
	login := message.Payload.Name
	for _, player := range r.Players {
		if player.ID != login {
			continue
		}
		fmt.Println(r.state)
		switch player.Type {
		case easyjson.PENGUIN:
			r.state.RotatePenguin()
		case easyjson.GUN:
			r.state.RotateGun()
		default:
			fmt.Println("Incorrect player type!")
		}
		break
	}
}

func (r *RoomMulti) FinishGame() *easyjson.OutcomeMessage {
	for _, player := range r.Players {
		helpers.LogMsg("Player " + player.ID + " finished game")
	}
	//r.gameState = FINISHED
	if r.state.Penguin.Score > r.state.Gun.Score {
		message := &easyjson.OutcomeMessage{
			Type: easyjson.FINISHGAME,
			Payload: easyjson.OutPayloadMessage{
				Penguin: easyjson.PenguinMessage{
					Name:   r.state.Penguin.ID,
					Score:  uint(r.state.Penguin.Score),
					Result: easyjson.WIN,
				},
				Gun: easyjson.GunMessage{
					Name:   r.state.Gun.ID,
					Score:  uint(r.state.Gun.Score),
					Result: easyjson.LOST,
				},
			},
		}

		r.gameState = easyjson.FINISHED
		return message
	} else {
		message := &easyjson.OutcomeMessage{
			Type: easyjson.FINISHGAME,
			Payload: easyjson.OutPayloadMessage{
				Penguin: easyjson.PenguinMessage{
					Name:   r.state.Penguin.ID,
					Score:  uint(r.state.Penguin.Score),
					Result: easyjson.LOST,
				},
				Gun: easyjson.GunMessage{
					Name:   r.state.Gun.ID,
					Score:  uint(r.state.Gun.Score),
					Result: easyjson.WIN,
				},
			},
		}
		r.gameState = easyjson.FINISHED
		return message
	}
	//r.state.Penguin = nil
	//r.state.Gun = nil
}

func (r *RoomMulti) FinishRound() {
	for _, player := range r.Players {
		helpers.LogMsg("Player " + player.ID + " finished round")
	}
	r.gameState = easyjson.WAITING
	if r.round == 2 {
			r.gameState = easyjson.FINISHED
	}
}

func (r *RoomMulti) SendRoomState(message *easyjson.OutcomeMessage) {
	for _, player := range r.Players {
		select {
		case player.out <- message:
		default:
			close(player.out)
		}
	}
}

func (r *RoomMulti) StartNewRound() {
	time.Sleep(1000 * time.Millisecond)
	if r.state != nil && r.round < 2 {
		r.round += 1
		r.state.Round = r.round
		//penguin, gun := r.SelectPlayersRoles()
		message := &easyjson.OutcomeMessage{
			Type: easyjson.START,
			Payload: easyjson.OutPayloadMessage{
				Gun: easyjson.GunMessage{
					//Name: gun,
					Name: r.state.Gun.ID,
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
		r.SendRoomState(message)
		r.state = CreateInitialState(r)
		r.gameState = easyjson.RUNNING
	} else {
		if r.round == 2 {
			message := r.FinishGame()
			r.SendRoomState(message)
			//r.gameState = FINISHED
		}
	}
}

func (r *RoomMulti) SaveResult() {
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
	for _, player := range r.Players {
		if player.Type == easyjson.PENGUIN {
			player.instance.Score = uint64(player.roomMulti.state.Penguin.Score)
		}
		if player.Type == easyjson.GUN {
			player.instance.Score = uint64(player.roomMulti.state.Gun.Score)
		}
		ctx := context.Background()
		_, err := AuthManager.SaveUserGame(ctx, player.instance)
		fmt.Println(err)
	}
}
