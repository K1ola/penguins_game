package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestConfig(t *testing.T) {

	go func() {
		main()
	}()
}

func TestGame(t *testing.T) {
	InitGame(20)
	game := NewGame(20)
	game.RoomsCount()
	var conn *websocket.Conn
	player := NewPlayer(conn)
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

}
