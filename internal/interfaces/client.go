package interfaces

const StageWaitOpponent = 1
const StageSettingShips = 2
const StageReady = 3
const StageWaitOpponentFair = 4
const StageShoot = 5
const StageEnd = 6

type GameClient interface {
	Send(message GameMessage)
	AddPair(pair GamePair)
	Listen()
	SetFirst(bool)
	SetStage(stage uint8)
	GetStage() uint8
	SendOpponent(m GameMessage) error
	GetOpponent() GameClient
	SendAll(m GameMessage) error
	ShootStage()
	GetCellMap() GameCellMap
	SuccessShoot(x, y int, m GameMessage)
	GotSuccessShoot(x, y int, m GameMessage)
	FailShoot(x, y int, m GameMessage)
	GotFailShoot(x, y int, m GameMessage)
	Win()
	Defeat()
	CloseConn()
}