package main

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

type TankKind uint8

const (
	Player TankKind = iota
	Enemy
)

type Tank struct {
	X, Y      int
	Direction Direction
	Kind      TankKind
}
