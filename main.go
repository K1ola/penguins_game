package main

import (
	"game/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func main() {
	//to local package in local parametr (will be tested)
	PingGame = InitGame(10)
	go PingGame.Run()

	router := mux.NewRouter()
	gameRouter := router.PathPrefix("/game").Subrouter()
	//TODO
	//router.Use(mw.PanicMiddleware)
	//router.Use(mw.CORSMiddleware)
	//router.Use(mw.AuthMiddleware)

	//router.HandleFunc("/", RootHandler)
	gameRouter.HandleFunc("/ws", StartWS)
	prometheus.MustRegister(metrics.PlayersCountInGame, metrics.ActiveRooms)
	router.Handle("/metrics", promhttp.Handler())


	LogMsg("GameServer started at :8085")

	http.ListenAndServe(":8085", handlers.LoggingHandler(os.Stdout, router))
}

