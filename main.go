package main

import (
	"game/metrics"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func setConfig() (string, int) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("game")
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
	prometheus.MustRegister(metrics.PlayersCountInGame, metrics.ActiveRooms)
	router.Handle("/metrics", promhttp.Handler())
	gameRouter.HandleFunc("/single", StartSingle)
	gameRouter.HandleFunc("/multi", StartMulti)

	LogMsg("GameServer started at", port)

	http.ListenAndServe(port, handlers.LoggingHandler(os.Stdout, router))
}
