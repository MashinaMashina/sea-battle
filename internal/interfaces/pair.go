package interfaces

type GamePair interface {
	SendOpponent(iamIsFirst bool, message GameMessage) error
	GetOpponent(bool) GameClient
	SendAll(message GameMessage) error
	IsFree() bool
	AddClient(client GameClient) error
	Remove(first bool) error
}

