package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func StartWS(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= int(maxRooms) {
		//TODO check response on the client side
		LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	upgrader := &websocket.Upgrader{}

	//check for multi in micro!!!!!

	//cookie, err := r.Cookie("sessionid")
	//if err != nil {
	//	cookie.Value = "Anonumys"
	//}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		LogMsg("Connection error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	LogMsg("Connected to client")

	//TODO remove hardcore, get from front player value
	player := NewPlayer(conn)
	go player.Listen()
}
