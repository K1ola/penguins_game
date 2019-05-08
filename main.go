package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func setConfig() (string, int) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("auth")
	var port string
	var maxRooms int
	if err := viper.ReadInConfig(); err != nil {
		port = ":8083"
		maxRooms = 10
	} else {
		port = ":" + viper.GetString("port")
		maxRooms = viper.GetInt("maxRooms")
	}
	return port, maxRooms
}

func main() {
	//to local package in local parametr (will be tested)
	port, maxRooms := setConfig()
	PingGame = InitGame(uint(maxRooms))
	go PingGame.Run()

	router := mux.NewRouter()
	gameRouter := router.PathPrefix("/game").Subrouter()
	//TODO
	//router.Use(mw.PanicMiddleware)
	//router.Use(mw.CORSMiddleware)
	//router.Use(mw.AuthMiddleware)

	//router.HandleFunc("/", RootHandler)
	gameRouter.HandleFunc("/ws", StartWS)

	LogMsg("GameServer started at", port)

	http.ListenAndServe(port, handlers.LoggingHandler(os.Stdout, router))
}
