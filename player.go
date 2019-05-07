package main

import (
	//"game/helpers"
	"fmt"
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

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		//ID:   id,
		in:   make(chan *IncomeMessage),
		out:  make(chan *OutcomeMessage, 1),
		roomMulti: nil,
		roomSingle: nil,
		Type: "penguin",
	}
}

func (p *Player) Listen() {
	go func() {
		for {
			//слушаем фронт
			message := &IncomeMessage{}
			err := p.conn.ReadJSON(message)
			fmt.Println("ReadJSON error: ", err)
			if websocket.IsUnexpectedCloseError(err) {
				p.RemovePlayerFromRoom()
				LogMsg("Player " + p.ID +" disconnected")
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
		//получаем состояние игры от фронтов
		case message := <-p.in:
			//оработать, посчитать
			fmt.Printf("Front says: %#v", message)
			fmt.Println("")
			switch message.Type {
				case NEWPLAYER:
					//стартовая инициализация, производится строго вначале один раз
					if message.Payload.Mode != "" {
						p.GameMode = message.Payload.Mode
						p.ID = message.Payload.Name
						PingGame.AddPlayer(p)
					}
				case NEWCOMMAND:
						//process command
				default:
					fmt.Println("Default in Player.Listen() - in")
			}

		case message := <-p.out:
			fmt.Printf("Back says: %#v", message)
			//шлем всем фронтам текущее состояние
			//switch message.Type {
			//	case START:
			//	case WAIT:
			//	case FINISH:
			//	case STATE:
			//	default:
			//		fmt.Println("Default in Player.Listen() - out")
			//}
			_ = p.conn.WriteJSON(message)
			//if p.GameMode != "" {
			//	if p.roomSingle != nil {
			//		p.roomSingle.broadcast <- message
			//	} else {
			//		p.roomMulti.broadcast <- message
			//	}
			//}
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

//func (p *Player) SendState(state *RoomState) {
//	//TODO: send to front
//	if state != nil {
//		//TODO create norm state
//		p.out <- &Message{"SINGLE", PayloadMessage{"STATE", "SOME-STATE"}}
//	}
//}
//
//func (p *Player) SendMessageSingle(message *OutcomeMessage) {
//	p.roomSingle.broadcast <- message
//}
