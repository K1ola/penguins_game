package main

import (
	//"game/helpers"
	"fmt"
	"game/helpers"
	"game/metrics"
	"sync"
)

var PingGame *Game
var maxRooms uint

func InitGame(rooms uint) *Game {
	maxRooms = rooms
	game :=  NewGame(maxRooms)
	//for i := 0; i < 10; i++ {
	//	//game.roomsSingle[i].gameState = INITIALIZED
	//	game.roomsMulti[i].gameState = INITIALIZED
	//}
	return game
}

type Game struct {
	MaxRooms    uint
	roomsSingle []*RoomSingle
	roomsMulti  []*RoomMulti
	mu       sync.RWMutex
	register chan *Player
	unregister chan *Player
	Players    map[string]*Player
}

func NewGame(maxRooms uint) *Game {
	return &Game{
		MaxRooms: maxRooms,
		register: make(chan *Player),
		unregister: make(chan *Player),
		Players:    make(map[string]*Player),
	}
}

func (g *Game) Run() {
	defer helpers.RecoverPanic()
//LOOP:
	for {
		select {

		case player, _ := <-g.register:
			switch player.GameMode {
				case SINGLE:
					g.ProcessSingle(player)
				case MULTI:
					g.ProcessMulti(player)
				default:
					fmt.Println("Empty")
				}
		case <-g.unregister:
			//remove from rooms
			//(do mot forget to free pointers - use same logic as with players)
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
	helpers.LogMsg("Player " + player.ID + " queued to add")
	g.mu.Lock()
	g.Players[player.ID] = player
	g.mu.Unlock()
	metrics.PlayersCountInGame.Inc()
	g.register <- player
}

func (g *Game) RoomsCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.roomsSingle) + len(g.roomsMulti)
}

func (g *Game) ProcessSingle(player *Player) {
	for _, room := range g.roomsSingle {
		if room.Player == nil {
			g.mu.Lock()
			room.AddPlayer(player)
			//TODO add game states
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
						Name: string(GUN),
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
			//continue LOOP
			return
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
	player.out <- &OutcomeMessage{
		Type: START,
		Payload:OutPayloadMessage{
			Gun:GunMessage{
				Bullet:BulletMessage{
					Alpha: 0,
					DistanceFromCenter: 0,
				},
				Alpha: 0,
				Name: string(GUN),
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
}

func (g *Game) ProcessMulti(player *Player) {
	for _, room := range g.roomsMulti {
		if room.gameState == PICKINGUP { //len(room.Players) < int(room.MaxPlayers) {
			g.mu.Lock()
			room.AddPlayer(player)
			g.mu.Unlock()

			room.state = CreateInitialState(room)

			//penguin, gun := room.SelectPlayersRoles()
			//message := &OutcomeMessage{
			//	Type: START,
			//	Payload: OutPayloadMessage{
			//		Gun: GunMessage{
			//			Name: gun,
			//		},
			//		Penguin: PenguinMessage{
			//			Name: penguin,
			//		},
			//		PiscesCount: 24,
			//	},
			//}
			//room.SendRoomState(message)
			room.StartNewRound()

			//room.gameState = RUNNING
			//continue LOOP
			return
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
	room.gameState = PICKINGUP
}
