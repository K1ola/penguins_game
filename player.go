package main

import (
	//"game/helpers"
	"fmt"
	"game/metrics"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	conn *websocket.Conn
	ID   string
	in   chan *IncomeMessage
	out  chan *OutcomeMessage
	roomSingle *RoomSingle
	roomMulti *RoomMulti
	GameMode string
	Type string
}

func NewPlayer(conn *websocket.Conn, id string) *Player {
	return &Player{
		conn: conn,
		ID:   id,
		in:   make(chan *IncomeMessage),
		out:  make(chan *OutcomeMessage, 1),
		roomMulti: nil,
		roomSingle: nil,
		Type: PENGUIN,
	}
}

func (p *Player) Listen() {
	//defer helpers.RecoverPanic()
	go func() {
		//defer helpers.RecoverPanic()
		for {
			//слушаем фронт
			message := &IncomeMessage{}
			err := p.conn.ReadJSON(message)
			fmt.Println("ReadJSON error: ", err)
			if websocket.IsUnexpectedCloseError(err) {
				p.RemovePlayerFromRoom()
				LogMsg("Player " + p.ID +" disconnected")
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
				case NEWPLAYER:
					//стартовая инициализация, производится строго вначале один раз
					if message.Payload.Mode != "" {
						p.GameMode = message.Payload.Mode
						//p.ID = message.Payload.Name
						PingGame.AddPlayer(p)
					}
				case NEWCOMMAND:
					//get name, do rotate
					//TODO select game mode
					p.roomMulti.ProcessCommand(message)

				default:
					fmt.Println("Default in Player.Listen() - in")
			}

		case message := <-p.out:
			fmt.Printf("Back says: %#v", message)
			fmt.Println("")
			//шлем всем фронтам текущее состояние
			switch message.Type {
				case START:
					fmt.Println("Process START")
				case WAIT:
					fmt.Println("Process WAIT")
				case FINISH:
					fmt.Println("Process FINISH")
				case STATE:
					fmt.Println("Process STATE")
				default:
					fmt.Println("Default in Player.Listen() - out")
			}
			_ = p.conn.WriteJSON(message)
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

func (p *Player) Finish() {
	if p.roomSingle != nil {
		//TODO finish single
		//p.roomSingle.(p)
	}
	if p.roomMulti != nil {
		p.roomMulti.FinishGame(p)
	}
}

