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

	player = NewPlayer(conn, "1")
	roomMulti = NewRoomMulti(2)
	player.ID = "user6"
	roomMulti.Players["test"] = player
	player.roomMulti = roomMulti
	roomMulti.Players["test"].Type = "PENGUIN"
	CreateInitialState(roomMulti)
	roomMulti.Players["test"].Type = "GUN"
	CreateInitialState(roomMulti)
	_ = RunMulti(roomMulti)
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
	player := NewPlayer(conn, "1")
	go player.Listen()
	message := &OutcomeMessage{}
	message.Type = "SIGNAL_START_THE_GAME"
	player.out <- message

	// message.Type = "SIGNAL_FINISH_GAME"
	// player.out <- message

	// message.Type = "SIGNAL_NEW_GAME_STATE"
	// player.out <- message

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
	incomeMessage.Type = "new"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

	incomeMessage.Type = "newPlayer"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

// func TestHandlerStartSingle(t *testing.T) {
// 	s := httptest.NewServer(http.HandlerFunc(StartSingle))
// 	defer s.Close()

// 	u := "ws" + strings.TrimPrefix(s.URL, "http")

// 	// Connect to the server
// 	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
// 	if err != nil {
// 		t.Fatalf("%v", err)
// 	}
// 	defer ws.Close()
// 	// incomeMessage := new(IncomeMessage)
// 	// incomeMessage.Type = "newPlayer"
// 	// incomeMessage.Payload.Name = "user6"
// 	// incomeMessage.Payload.Mode = "MULTI"

// 	// if err := ws.WriteJSON(""); err != nil {
// 	// 	t.Fatalf("%v", err)
// 	// }
// }

func Player3(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	player := NewPlayer(conn, "1")
	go player.Listen()
	message := &OutcomeMessage{}
	// message.Type = "SIGNAL_START_THE_GAME"
	// player.out <- message

	message.Type = "SIGNAL_FINISH_GAME"
	player.out <- message

	// message.Type = "SIGNAL_NEW_GAME_STATE"
	// player.out <- message

}

func TestPlayer3(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player3))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "new"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

	incomeMessage.Type = "newPlayer"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

func Player4(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	player := NewPlayer(conn, "1")
	go player.Listen()
	message := &OutcomeMessage{}
	// message.Type = "SIGNAL_START_THE_GAME"
	// player.out <- message

	// message.Type = "SIGNAL_FINISH_GAME"
	// player.out <- message

	message.Type = "SIGNAL_NEW_GAME_STATE"
	player.out <- message

}

func TestPlayer4(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player4))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "new"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

	incomeMessage.Type = "newPlayer"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

func Player5(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	player := NewPlayer(conn, "1")
	go player.Listen()
	message := &OutcomeMessage{}
	// message.Type = "SIGNAL_START_THE_GAME"
	// player.out <- message

	// message.Type = "SIGNAL_FINISH_GAME"
	// player.out <- message

	message.Type = "SIGNAL_TO_WAIT_OPPONENT"
	player.out <- message

}

func TestPlayer5(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player5))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "new"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

	incomeMessage.Type = "newPlayer"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

func Player6(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	player := NewPlayer(conn, "1")
	go player.Listen()
	message := &OutcomeMessage{}
	// message.Type = "SIGNAL_START_THE_GAME"
	// player.out <- message

	// message.Type = "SIGNAL_FINISH_GAME"
	// player.out <- message

	message.Type = "S"
	player.out <- message

}

func TestPlayer6(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player5))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)
	incomeMessage.Type = "new"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

	incomeMessage.Type = "newPlayer"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

func Player7(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	player := NewPlayer(conn, "1")
	roomMulti := NewRoomMulti(2)
	player.ID = "user6"
	roomMulti.Players["test"] = player
	player.roomMulti = roomMulti
	roomMulti.Players["test"].Type = "PENGUIN"
	incomeMessage := &IncomeMessage{}
	incomeMessage.Payload.Name = player.ID
	roomMulti.ProcessCommand(incomeMessage)

	roomMulti.Players["test"].Type = "GUN"
	roomMulti.ProcessCommand(incomeMessage)

	roomMulti.Players["test"].Type = "def"
	roomMulti.ProcessCommand(incomeMessage)

	incomeMessage.Payload.Name = ""
	roomMulti.ProcessCommand(incomeMessage)

	roomMulti.SelectPlayersRoles()
	roomMulti.SelectPlayersRoles()
	roomMulti.SelectPlayersRoles()
	roomMulti.SelectPlayersRoles()

	// roomSingle := NewRoomSingle(2)
	// player.roomSingle = roomSingle
	go player.Listen()
	message := &OutcomeMessage{}
	// roomMulti.SendRoomState(message)

	// message.Type = "SIGNAL_START_THE_GAME"
	// player.out <- message

	// message.Type = "SIGNAL_FINISH_GAME"
	// player.out <- message

	message.Type = "S"
	player.out <- message

	player.RemovePlayerFromRoom()

}

func TestPlayer7(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Player7))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)

	incomeMessage.Type = "newCommand"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}

func TestRecalcBullet(t *testing.T) {
	roomMulti := NewRoomMulti(2)
	rs := CreateInitialState(roomMulti)
	rs.Gun.Bullet.Alpha = 30
	rs.Penguin.Alpha = 30
	rs.Gun.Bullet.DistanceFromCenter = 2000
	rs.RecalcBullet()
}

func TestPlayer8(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(StartSingle))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	incomeMessage := new(IncomeMessage)

	incomeMessage.Type = "newCommand"
	incomeMessage.Payload.Name = "user6"
	incomeMessage.Payload.Mode = "MULTI"

	if err := ws.WriteJSON(incomeMessage); err != nil {
		t.Fatalf("%v", err)
	}

}
