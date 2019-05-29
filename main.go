package main

import (
	"game/helpers"
	"game/metrics"
	mw "game/middleware"
	"game/models"
	"google.golang.org/grpc"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func setConfig() (string, int, string) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("game")
	var port, authAddress string
	var maxRooms int
	if err := viper.ReadInConfig(); err != nil {
		port = ":8085"
		maxRooms = 10
		authAddress = "127.0.0.1:8083"
	} else {
		port = ":" + viper.GetString("port")
		maxRooms = viper.GetInt("maxRooms")
		authAddress = viper.GetString("auth")
	}
	return port, maxRooms, authAddress
}

func main() {
	port, maxRooms, authAddress := setConfig()
	PingGame = InitGame(uint(maxRooms))
	go PingGame.Run()

	grcpConn, err := grpc.Dial(
		authAddress,
		grpc.WithInsecure(),
	)
	if err != nil {
		helpers.LogMsg("Can`t connect to grpc")
		return
	}
	defer grcpConn.Close()

	models.AuthManager = models.NewAuthCheckerClient(grcpConn)

	router := mux.NewRouter()
	gameRouter := router.PathPrefix("/game").Subrouter()
	//userRouter := router.PathPrefix("/data").Subrouter()

	//TODO
	//router.Use(mw.PanicMiddleware)
	gameRouter.Use(mw.CORSMiddleware)
	//router.Use(mw.AuthMiddleware)

	//router.HandleFunc("/", RootHandler)
	prometheus.MustRegister(metrics.PlayersCountInGame, metrics.ActiveRooms)
	router.Handle("/metrics", promhttp.Handler())
	gameRouter.HandleFunc("/single", StartSingle)
	gameRouter.HandleFunc("/multi", StartMulti)
	//userRouter.HandleFunc("/checkSingleWs", CheckWsSingle).Methods("GET", "OPTIONS")
	//userRouter.HandleFunc("/checkMultiWs", CheckWsMulti).Methods("GET", "OPTIONS")

	helpers.LogMsg("GameServer started at", port)

	http.ListenAndServe(port, handlers.LoggingHandler(os.Stdout, router))
}
