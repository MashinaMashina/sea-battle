package interfaces

type GameMessage interface {
	GetCode() string
	GetMessage() []byte
	SetCode(string)
	SetMessage(interface{})
	GetJSON() ([]byte, error)
}
