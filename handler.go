package main

import (
	"game/helpers"
	"game/models"
	"net/http"

	"github.com/gorilla/websocket"

)

var AuthManager models.AuthCheckerClient
var user *models.User


func StartSingle(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		helpers.LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	//user = new(models.User)

	//ctx := context.Background()
	//user, _ = models.AuthManager.Get(ctx, &models.JWT{Token: cookie.Value})

	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogMsg("Connection error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	helpers.LogMsg("Connected to client")

	//TODO remove hardcore, get from front player value
	player := NewPlayer(conn)
	//player.ID = user.Login
	go player.Listen()
}

func StartMulti(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		helpers.LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	//user = new(models.User)

	//cookie, err := r.Cookie("sessionid")
	//fmt.Println(cookie)
	// fmt.Println(cookie.Value)

	//if err != nil || cookie.Value == "" {
	//	helpers.LogMsg("No Cookie in Multi")
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//} else {
	//	ctx := context.Background()

		//user, err = models.AuthManager.GetUser(ctx, &models.JWT{Token: cookie.Value})
		//if err != nil {
		//	helpers.LogMsg("Invalid Cookie in Multi")
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}

		//check if such user already in game
		//fmt.Println(PingGame.Players)
		//for _, player := range PingGame.Players {
		//	if player.ID == user.Login {
		//		helpers.LogMsg("Such user already in game")
		//		//TODO and what's next with this user?
		//		w.WriteHeader(http.StatusForbidden)
		//		return
		//	}
		//}
	//}
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogMsg("Connection error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	helpers.LogMsg("Connected to client")

	player := NewPlayer(conn)
	go player.Listen()
}
