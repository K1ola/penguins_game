package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConfig(t *testing.T) {

	go func() {
		main()
	}()
}

func TestGame(t *testing.T) {
	InitGame(10)
	game := NewGame(10)
	game.RoomsCount()
	var conn *websocket.Conn
	player := NewPlayer(conn, "1")
	go func() {
		game.AddPlayer(player)
	}()

	player.RemovePlayerFromRoom()
	roomMulti := NewRoomMulti(2)
	go func() {
		roomMulti.AddPlayer(player)
	}()
	go func() {
		roomMulti.RemovePlayer(player)
	}()
	go func() {
		roomMulti.Run()
	}()
	go func() {
		roomMulti.SelectPlayersRoles()
	}()
	game.AddToRoomMulti(roomMulti)
	roomSingle := NewRoomSingle(2)
	game.AddToRoomSingle(roomSingle)
	go func() {
		roomSingle.AddPlayer(player)
	}()
	go func() {
		roomSingle.RemovePlayer(player)
	}()
	go func() {
		roomSingle.Run()
	}()
}

func TestEngine(t *testing.T) {
	penguinState := CreatePenguin("1")
	fmt.Println(penguinState)
	gunState := CreateGun("1")
	fmt.Println(gunState)
	bulletState := CreateBullet()
	fmt.Println(bulletState)
	fishStates := CreateFishes()
	fmt.Println(fishStates)
	roomMulti := NewRoomMulti(2)
	_ = RunMulti(roomMulti)
	roomState := CreateInitialState(roomMulti)
	roomState.RecalcGun()
	roomState.Gun.Alpha = -1
	roomState.RecalcGun()
	roomState.Gun.Alpha = 360
	roomState.RecalcGun()
	roomState.RotateGun()
	roomState.RecalcGun()
	roomState.RecalcBullet()
	roomState.RecalcPenguin()
	roomState.GetState()
	roomState.GetState()
	roomState.RotatePenguin()
	roomState.RotateGun()
	roomState.RecalcPenguin()
	roomState.Penguin.Alpha = -1
	roomState.RecalcPenguin()
	roomState.Penguin.Alpha = 360
	roomState.RecalcPenguin()
	roomState.RecalcPenguin()
	roomState.RotateGun()
	roomState.Gun.Bullet.DistanceFromCenter = 100*0.8/2 + 1
	roomState.RecalcBullet()
	roomState.Gun.Bullet.Alpha = 30
	roomState.Penguin.Alpha = 30
	roomState.RecalcBullet()
	roomState.RotatePenguin()
	message := roomState.GetState()
	roomMulti.SendRoomState(message)
	var conn *websocket.Conn
	player := NewPlayer(conn, "1")
	player.Finish()
	player.roomMulti = roomMulti
	player.Finish()
	roomMulti.FinishGame(player)
}

func TestHandlers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/single", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(StartSingle)
	// roomMulti := NewRoomMulti(2)
	// roomSingle := NewRoomSingle(2)
	PingGame = NewGame(10)
	handler.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/multi", nil)
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(StartMulti)
	// roomMulti := NewRoomMulti(2)
	roomSingle := NewRoomSingle(2)
	handler.ServeHTTP(w, req)

	// cookie := &http.Cookie{
	// 	Name:     "sessionid",
	// 	Value:    "session",
	// 	Expires:  time.Now().AddDate(0, 0, 1),
	// 	HttpOnly: true,
	// }
	// req, _ = http.NewRequest("GET", "/multi", nil)
	// w = httptest.NewRecorder()
	// req.AddCookie(cookie)
	// handler = http.HandlerFunc(StartMulti)
	// handler.ServeHTTP(w, req)

	for index := 0; index < 11; index++ {
		PingGame.roomsSingle = append(PingGame.roomsSingle, roomSingle)
	}
	req, _ = http.NewRequest("GET", "/multi", nil)
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(StartMulti)
	handler.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/single", nil)
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(StartSingle)
	handler.ServeHTTP(w, req)

}

func TestPlayer(t *testing.T) {
	roomSingle := NewRoomSingle(2)
	// upgrader := &websocket.Upgrader{}
	// w := httptest.NewRecorder()
	// r, _ := http.NewRequest("GET", "/game/multi", nil)
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	t.Error(err)
	// }
	// incomeMessage := new(IncomeMessage)
	var conn *websocket.Conn
	// conn.WriteJSON(incomeMessage)
	player := NewPlayer(conn, "1")
	player.roomSingle = roomSingle
	player.Finish()
	// player.Listen()
}

func Player2(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "newRound"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"
	player := NewPlayer(conn, "1")
	go func() {
		err = player.conn.WriteJSON(incomeMessage)
		fmt.Println(err)
	}()
	// player.Listen()

	//player.in <- incomeMessage
	go player.Listen()
}

func TestPlayer2(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player2))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "newRound"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}
}

// func TestGameRun(t *testing.T) {

// }
