package main

import (
	"game/models"
	"net/http"

	"github.com/gorilla/websocket"

	"golang.org/x/net/context"
)

func StartSingle(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	user := new(models.User)
	cookie, err := r.Cookie("sessionid")
	if err != nil || cookie.Value == "" {
		user.Login = "Anonumys"
	} else {
		ctx := context.Background()

		user, _ = models.AuthManager.GetUser(ctx, &models.JWT{Token: cookie.Value})
		cookie.Value = user.Login
	}

	upgrader := &websocket.Upgrader{}

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
	player.ID = user.Login
	go player.Listen()
}

func StartMulti(w http.ResponseWriter, r *http.Request) {
	if PingGame.RoomsCount() >= 10 {
		//TODO check response on the client side
		LogMsg("Too many clients")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many clients"))
		return
	}

	var user *models.User
	cookie, err := r.Cookie("sessionid")
	if err != nil || cookie.Value == "" {
		LogMsg("No Cookie in Multi")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		//grcpConn, err := grpc.Dial(
		//	"127.0.0.1:8083",
		//	grpc.WithInsecure(),
		//)
		//if err != nil {
		//	helpers.LogMsg("Can`t connect to grpc")
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}
		//defer grcpConn.Close()
		//
		//authManager := models.NewAuthCheckerClient(grcpConn)
		ctx := context.Background()

		user, err = models.AuthManager.GetUser(ctx, &models.JWT{Token: cookie.Value})
		if err != nil {
			LogMsg("Invalid Cookie in Multi")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cookie.Value = user.Login
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
