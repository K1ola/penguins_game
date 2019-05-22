package main

import (
	//"game/helpers"
	"fmt"
	"game/easyjson"
	"game/helpers"
	"game/metrics"
	"game/models"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	instance   *models.User
	conn       *websocket.Conn
	ID         string
	game       *Game
	in         chan *easyjson.IncomeMessage
	out        chan *easyjson.OutcomeMessage
	roomSingle *RoomSingle
	roomMulti  *RoomMulti
	GameMode   easyjson.GameMode
	Type       easyjson.ClientRole
	Playing    bool
}

func NewPlayer(conn *websocket.Conn, id string, instance *models.User) *Player {
	return &Player{
		instance:   instance,
		conn:       conn,
		ID:         id,
		game:       PingGame,
		in:         make(chan *easyjson.IncomeMessage),
		out:        make(chan *easyjson.OutcomeMessage, 1),
		roomMulti:  nil,
		roomSingle: nil,
		Type:       easyjson.PENGUIN,
		Playing:    false,
	}
}

func (p *Player) Listen() {
	//defer helpers.RecoverPanic()
	go func() {
		//defer helpers.RecoverPanic()
		for {
			//слушаем фронт
			message := &easyjson.IncomeMessage{}
			err := p.conn.ReadJSON(message)
			fmt.Println("ReadJSON error: ", err)
			if websocket.IsUnexpectedCloseError(err) {
				p.RemovePlayerFromRoom()
				p.RemovePlayerFromGame()
				helpers.LogMsg("Player " + p.ID +" disconnected")
				metrics.PlayersCountInGame.Dec()
				return
			}
			if err != nil {
				log.Printf("Cannot read json")
				continue
			}
			p.in <- message
		}
	}()

	for {
		select {
		//получаем команды от фронтов
		case message := <-p.in:
			fmt.Printf("Front says: %#v", message)
			fmt.Println("")
			switch message.Type {
				case easyjson.NEWPLAYER:
					//стартовая инициализация, производится строго вначале один раз
					if message.Payload.Mode != "" {
						p.GameMode = message.Payload.Mode
						//p.ID = message.Payload.Name
						PingGame.AddPlayer(p)
					}
				case easyjson.NEWCOMMAND:
					//get name, do rotate
					//TODO select game mode
					if message.Payload.Mode == easyjson.MULTI {
						p.roomMulti.ProcessCommand(message)
					}
					if message.Payload.Mode == easyjson.SINGLE {
						p.roomSingle.ProcessCommand(message)
					}

				case easyjson.NEWROUND:
					if p.roomMulti.gameState == easyjson.WAITING {
						fmt.Println(p.roomMulti)
						p.roomMulti.SendRoomState(&easyjson.OutcomeMessage{Type: easyjson.WAIT})
						p.roomMulti.gameState = easyjson.INITIALIZED
						continue
					}
					p.roomMulti.StartNewRound()
				default:
					fmt.Println("Default in Player.Listen() - in")
			}

		case message := <-p.out:
			fmt.Printf("Back says: %#v", message)
			fmt.Println("")
			//шлем всем фронтам текущее состояние
			if message != nil {
				switch message.Type {
				case easyjson.START:
					fmt.Println("Process START")
				case easyjson.WAIT:
					fmt.Println("Process WAIT")
				case easyjson.FINISHROUND:
					fmt.Println("Process FINISH ROUND")
				case easyjson.FINISHGAME:
					fmt.Println("Process FINISH GAME")
				case easyjson.STATE:
					fmt.Println("Process STATE")
				default:
					fmt.Println("Default in Player.Listen() - out")
				}
				_ = p.conn.WriteJSON(message)
			}
		}
	}
}

func (p *Player) RemovePlayerFromRoom() {
	if p.roomSingle != nil {
		p.roomSingle.RemovePlayer(p)
	}
	if p.roomMulti != nil {
		p.roomMulti.RemovePlayer(p)
	}
}

func (p *Player) RemovePlayerFromGame() {
	p.game.unregister <- p
}

func (p *Player) FinishGame() {
	if p.roomSingle != nil {
		//TODO finish single
		//p.roomSingle.(p)
	}
	if p.roomMulti != nil {
		p.roomMulti.FinishGame()
	}
}

func (p *Player) FinishRound() {
	if p.roomSingle != nil {
		//TODO finish single
		//p.roomSingle.(p)
	}
	if p.roomMulti != nil {
		p.roomMulti.FinishRound()
	}
}

