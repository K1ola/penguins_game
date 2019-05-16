package main

import (
	"fmt"
	"game/helpers"
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
	state      *RoomState
	gameState GameCurrentState
	round int

	broadcast chan *OutcomeMessage
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
		state: &RoomState{
			Penguin: new(PenguinState),
			Gun: new(GunState),
			Fishes: make(map[int]*FishState, 24),
			Round: -1,
		},
		round: -1,
		broadcast: make(chan *OutcomeMessage),
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
			//if len(r.Players) == 2 {
				//if r.state != nil && r.round < 2 {
				//	penguin, gun := r.SelectPlayersRoles()
				//	message := &OutcomeMessage{
				//		Type: START,
				//		Payload: OutPayloadMessage{
				//			Gun: GunMessage{
				//				Name: gun,
				//			},
				//			Penguin: PenguinMessage{
				//				Name: penguin,
				//			},
				//			PiscesCount: 24,
				//		},
				//	}
				//	r.gameState = START
				//	r.SendRoomState(message)
				//	r.state = CreateInitialState(r)
					//r.StartNewRound()
				//} else {
				//	if r.round >= 2 {
				//		r.gameState = FINISHGAME
				//		if r.state.Penguin.Score > r.state.Gun.Score {
				//			message := &OutcomeMessage{
				//				Type:FINISHGAME,
				//				Payload:OutPayloadMessage{
				//					Penguin:PenguinMessage{
				//						Name: r.state.Penguin.ID,
				//						Score: uint(r.state.Penguin.Score),
				//						Result:WIN,
				//					},
				//					Gun:GunMessage{
				//						Name: r.state.Gun.ID,
				//						Score: uint(r.state.Gun.Score),
				//						Result:LOST,
				//					},
				//				},
				//			}
				//
				//			r.SendRoomState(message)
				//			r.gameState = FINISHGAME
				//		} else {
				//			message := &OutcomeMessage{
				//				Type:FINISHGAME,
				//				Payload:OutPayloadMessage{
				//					Penguin:PenguinMessage{
				//						Name: r.state.Penguin.ID,
				//						Score: uint(r.state.Penguin.Score),
				//						Result:LOST,
				//					},
				//					Gun:GunMessage{
				//						Name: r.state.Gun.ID,
				//						Score: uint(r.state.Gun.Score),
				//						Result:WIN,
				//					},
				//				},
				//			}
				//
				//			r.SendRoomState(message)
				//			r.gameState = FINISHGAME
				//		}
						//send FINISH_GAME
					//}
				//}
			//}
		case message := <- r.broadcast:
			r.SendRoomState(message)
		case <-r.ticker.C:
			if r.gameState == RUNNING {
				  message := RunMulti(r)
				  if message.Type != STATE {
					  switch message.Type {
					  case FINISHROUND:
					  		fmt.Println(FINISHROUND)
					  		fmt.Println(r.gameState)
					  case FINISHGAME:
							message = r.FinishGame()
					  }
				  }
				  r.SendRoomState(message)
				if r.round == 2 && r.gameState == FINISHED {
					message := r.FinishGame()
					r.SendRoomState(message)
					//r.gameState = FINISHED
					return
				}
			}
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
	//r.unregister <- player
	delete(r.Players, player.ID)
	helpers.LogMsg("Player " + player.ID + " was removed from room")
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

func (r *RoomMulti) FinishGame() *OutcomeMessage {
	for _, player := range r.Players {
		helpers.LogMsg("Player " + player.ID + " finished game")
	}
	//r.gameState = FINISHED
	if r.state.Penguin.Score > r.state.Gun.Score {
		message := &OutcomeMessage{
			Type: FINISHGAME,
			Payload: OutPayloadMessage{
				Penguin: PenguinMessage{
					Name:   r.state.Penguin.ID,
					Score:  uint(r.state.Penguin.Score),
					Result: WIN,
				},
				Gun: GunMessage{
					Name:   r.state.Gun.ID,
					Score:  uint(r.state.Gun.Score),
					Result: LOST,
				},
			},
		}

		r.gameState = FINISHED
		return message
	} else {
		message := &OutcomeMessage{
			Type: FINISHGAME,
			Payload: OutPayloadMessage{
				Penguin: PenguinMessage{
					Name:   r.state.Penguin.ID,
					Score:  uint(r.state.Penguin.Score),
					Result: LOST,
				},
				Gun: GunMessage{
					Name:   r.state.Gun.ID,
					Score:  uint(r.state.Gun.Score),
					Result: WIN,
				},
			},
		}
		r.gameState = FINISHED
		return message
	}
	//r.state.Penguin = nil
	//r.state.Gun = nil
}

func (r *RoomMulti) FinishRound() {
	for _, player := range r.Players {
		helpers.LogMsg("Player " + player.ID + " finished round")
	}
	r.gameState = WAITING
	if r.round == 2 {
			r.gameState = FINISHED
	}
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

func (r *RoomMulti) StartNewRound() {
	time.Sleep(500 * time.Millisecond)
	if r.state != nil && r.round < 2 {
		r.round += 1
		r.state.Round = r.round
		//penguin, gun := r.SelectPlayersRoles()
		message := &OutcomeMessage{
			Type: START,
			Payload: OutPayloadMessage{
				Gun: GunMessage{
					//Name: gun,
					Name: r.state.Gun.ID,
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
		r.SendRoomState(message)
		r.state = CreateInitialState(r)
		r.gameState = RUNNING
	} else {
		if r.round == 2 {
			message := r.FinishGame()
			r.SendRoomState(message)
			//r.gameState = FINISHED
		}
	}
}
