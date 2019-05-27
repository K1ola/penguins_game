package main

import (
	"fmt"
	"game/helpers"
	"game/models"
	"net/http"

	"github.com/gorilla/websocket"

	"golang.org/x/net/context"
)

var AuthManager models.AuthCheckerClient
var user *models.User

func CheckWsSingle(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		helpers.LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	user = new(models.User)

	cookie, err := r.Cookie("sessionid")
	if err != nil || cookie.Value == "" {
		user.Login = "Anonumys"
	} else {
		ctx := context.Background()
		user, _ = models.AuthManager.GetUser(ctx, &models.JWT{Token: cookie.Value})
	}
	w.WriteHeader(http.StatusOK)
}

func CheckWsMulti(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		helpers.LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	user = new(models.User)
	
	cookie, err := r.Cookie("sessionid")
	fmt.Println(cookie)
	// fmt.Println(cookie.Value)

	if err != nil || cookie.Value == "" {
		helpers.LogMsg("No Cookie in Multi")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		ctx := context.Background()

		user, err = models.AuthManager.GetUser(ctx, &models.JWT{Token: cookie.Value})
		if err != nil {
			helpers.LogMsg("Invalid Cookie in Multi")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//check if such user already in game
		fmt.Println(PingGame.Players)
		for _, player := range PingGame.Players {
			if player.ID == user.Login {
				helpers.LogMsg("Such user already in game")
				//TODO and what's next with this user?
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func StartSingle(w http.ResponseWriter, r *http.Request) {
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
	player := NewPlayer(conn, user.Login, user)
	player.ID = user.Login
	go player.Listen()
}

func StartMulti(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogMsg("Connection error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	helpers.LogMsg("Connected to client")

	player := NewPlayer(conn, user.Login, user)
	go player.Listen()
}
