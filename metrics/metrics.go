package metrics

import "github.com/prometheus/client_golang/prometheus"

var PlayersCountInGame = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "players_in_game",
	Help: "Number of players successfully added to the Game.",
})

var ActiveRooms = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "active_rooms",
	Help: "Number of active rooms in the Game.",
})
