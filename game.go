package main

import (
	//"game/helpers"
	"fmt"
	"game/metrics"
	"sync"
)

var PingGame *Game

const (
	SINGLE = "SINGLE"
	MULTI  = "MULTI"

	WAIT   = "SIGNAL_TO_WAIT_OPPONENT"
	START  = "SIGNAL_START_THE_GAME"
	FINISH = "SIGNAL_FINISH_GAME"
	STATE  = "SIGNAL_NEW_GAME_STATE"

	NEWPLAYER  = "newPlayer"
	NEWCOMMAND = "newCommand"

	ROTATE = "ROTATE"
	SHOT = "SHOT"

	PENGUIN = "PENGUIN"
	GUN = "GUN"
)

var maxRooms uint

func InitGame(rooms uint) *Game {
	maxRooms = rooms
	return NewGame(maxRooms)
}

type Game struct {
	MaxRooms    uint
	roomsSingle []*RoomSingle
	roomsMulti  []*RoomMulti
	//mu *sync.Mutex
	mu       sync.RWMutex
	register chan *Player
	unregister chan *Player
}

func NewGame(maxRooms uint) *Game {
	return &Game{
		MaxRooms: maxRooms,
		register: make(chan *Player),
		unregister: make(chan *Player),
	}
}

func (g *Game) Run() {
LOOP:
	for {
		select {
		case player, _ := <-g.register:
			//fmt.Println("register ch is ", ok)
			//fmt.Println("State is "+ player.GameMode)

			switch player.GameMode {
				case SINGLE:
					//start roomSingle
					for _, room := range g.roomsSingle {
						if room.Player == nil {
							g.mu.Lock()
							room.AddPlayer(player)
							g.mu.Unlock()
							player.out <- &OutcomeMessage{
								Type: START,
								Payload:OutPayloadMessage{
									Gun:GunMessage{
										Bullet:BulletMessage{
											Alpha: 0,
											DistanceFromCenter: 0,
										},
										Alpha: 0,
										Name: GUN,
										Result: "",
									},
									Penguin:PenguinMessage{
										Clockwise:false,
										Alpha: 0,
										Name: player.ID,
										Result: "",
									},
								},
							}
							continue LOOP
						}
					}

					//если все комнаты заняты - делой новую
					room := NewRoomSingle(1)
					g.mu.Lock()
					g.AddToRoomSingle(room)
					g.mu.Unlock()

					go room.Run()

					g.mu.Lock()
					room.AddPlayer(player)
					g.mu.Unlock()
					//room.broadcast <- &OutcomeMessage{Type: START}
					player.out <- &OutcomeMessage{
						Type: START,
						Payload:OutPayloadMessage{
							Gun:GunMessage{
								Bullet:BulletMessage{
									Alpha: 0,
									DistanceFromCenter: 0,
								},
								Alpha: 0,
								Name: GUN,
								Result: "",
							},
							Penguin:PenguinMessage{
								Clockwise:false,
								Alpha: 0,
								Name: player.ID,
								Result: "",
							},
						},
					}
				case MULTI:
					//start roomMulty
					var penguin, gun string

					for _, room := range g.roomsMulti {
						if len(room.Players) < int(room.MaxPlayers) {
							g.mu.Lock()
							room.AddPlayer(player)
							///////////////////
							room.SelectPlayersRoles()
							for _, player := range room.Players {
								if player.Type == PENGUIN {
									penguin = player.ID
								} else {
									gun = player.ID
								}
							}
							//player.out <- &OutcomeMessage{Type:START}
							room.broadcast <- &OutcomeMessage{
								Type: START,
								Payload:OutPayloadMessage{
									Gun:GunMessage{
										Bullet:BulletMessage{
											Alpha: 0,
											DistanceFromCenter: 0,
										},
										Alpha: 0,
										Name: gun,
										Result: "",
									},
									Penguin:PenguinMessage{
										Clockwise:false,
										Alpha: 0,
										Name: penguin,
										Result: "",
									},
								},
							}
							g.mu.Unlock()
							continue LOOP
						}
					}

					//если все комнаты заняты - делой новую
					room := NewRoomMulti(2)
					g.mu.Lock()
					g.AddToRoomMulti(room)
					g.mu.Unlock()

					go room.Run()

					g.mu.Lock()
					room.AddPlayer(player)
					player.out <- &OutcomeMessage{Type: WAIT}
					g.mu.Unlock()
				default:
					fmt.Println("Empty")
				}
		case <-g.unregister:
			//remove from rooms
		}

	}
}

func (g *Game) AddToRoomSingle(room *RoomSingle) {
	metrics.ActiveRooms.Inc()
	g.roomsSingle = append(g.roomsSingle, room)
}

func (g *Game) AddToRoomMulti(room *RoomMulti) {
	metrics.ActiveRooms.Inc()
	g.roomsMulti = append(g.roomsMulti, room)
}

func (g *Game) AddPlayer(player *Player) {
	LogMsg("Player " + player.ID + " queued to add")
	metrics.PlayersCountInGame.Inc()
	g.register <- player
}

func (g *Game) RoomsCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.roomsSingle) + len(g.roomsMulti)
}
