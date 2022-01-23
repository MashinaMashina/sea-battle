package interfaces

type MapCell interface {
	AliveCell() bool
	HasShip() bool
	SetAliveCell(bool)
	SetHasShip(bool)
}

type GameCellMap interface {
	From(json []byte) error
	ExistsShip(x, y int) bool
	Validate() error
	HasAliveShips() bool
	DamageCell(x, y int) error
}
