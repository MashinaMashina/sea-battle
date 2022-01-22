package interfaces

type GamePair interface {
	SendOpponent(message GameMessage) error
	SendAll(message GameMessage) error
	IsFree() bool
	AddClient(client GameClient) error
	Remove(first bool) error
}

