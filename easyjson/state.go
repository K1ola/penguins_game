package easyjson

type GameMode string
const (
	SINGLE GameMode = "SINGLE"
	MULTI  GameMode = "MULTI"
)

type ClientCommand string
const (
	NEWPLAYER ClientCommand = "newPlayer"
	NEWCOMMAND ClientCommand = "newCommand"
	NEWROUND ClientCommand = "newRound"
)

type ServerMessage string
const (
	WAIT  ServerMessage = "SIGNAL_TO_WAIT_OPPONENT"
	START  ServerMessage = "SIGNAL_START_THE_GAME"
	FINISHGAME ServerMessage = "SIGNAL_FINISH_GAME"
	FINISHROUND ServerMessage = "SIGNAL_FINISH_ROUND"
	STATE ServerMessage = "SIGNAL_NEW_GAME_STATE"
)

type ClientRole string
const (
	PENGUIN ClientRole = "PENGUIN"
	GUN ClientRole = "GUN"
)

type GameResult string
const (
	WIN GameResult = "WIN"
	LOST GameResult = "LOST"
)

type GameCurrentState string
const (
	INITIALIZED GameCurrentState = "INITIALIZED"	//game is ready for players
	PICKINGUP GameCurrentState = "PICKINGUP"	//game is collecting players to start
	WAITING  GameCurrentState = "WAITING"	//game is waiting players for new round
	RUNNING  GameCurrentState = "RUNNING"	//game is running
	FINISHED GameCurrentState = "FINISHED"	//game has finished
)

