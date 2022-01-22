package interfaces

type GameClient interface {
	Send(message GameMessage)
	AddPair(pair GamePair)
	Listen()
	SetFirst(bool)
}
