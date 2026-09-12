package game

import (
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

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

type Side uint8

const (
	Bottom Side = iota
	Top
)

func (s Side) opposite() Side {
	return s ^ 1
}

type Tank struct {
	Col, Row  int
	Direction Direction
	Side      Side
	Kind      TankKind
	LastFire  time.Time
}

var spawns = [2]Tank{
	Bottom: {Col: 13, Row: 22, Direction: Up},
	Top:    {Col: 12, Row: 3, Direction: Down},
}

func NewTank(side Side, kind TankKind) Tank {
	t := spawns[side]
	t.Side = side
	t.Kind = kind
	return t
}

func respawn(t *Tank) {
	spawn := spawns[t.Side]
	t.Col, t.Row, t.Direction = spawn.Col, spawn.Row, spawn.Direction
}

var (
	rowDelta = [4]int{-1, 1, 0, 0}
	colDelta = [4]int{0, 0, -1, 1}
)

func TryMoveTank(state *GameState) {
	t := &state.Tank
	newRow := t.Row + rowDelta[t.Direction]
	newCol := t.Col + colDelta[t.Direction]

	if newRow < 0 || newRow >= MapSize || newCol < 0 || newCol >= MapSize {
		return
	}
	if !tilePassable(state.Map.Grid[newRow][newCol]) {
		return
	}
	if state.Enemy != nil && state.Enemy.Row == newRow && state.Enemy.Col == newCol {
		return
	}

	t.Row = newRow
	t.Col = newCol
	if AbsorbBullets(state, *t) {
		respawn(t)
	}
}

var tankGlyphs = [4]rune{'^', 'v', '<', '>'}

func RenderTank(screen tcell.Screen, t Tank) {
	g := tankGlyphs[t.Direction]
	style := tcell.StyleDefault.Foreground(color.Green)
	if t.Kind == Enemy {
		style = tcell.StyleDefault.Foreground(color.Red)
	}
	ox, oy := MapOffset(screen)
	screen.SetContent(ox+t.Col*2, oy+t.Row, g, nil, style)
	screen.SetContent(ox+t.Col*2+1, oy+t.Row, g, nil, style)
}
