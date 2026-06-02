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

type Tank struct {
	Col, Row           int
	SpawnCol, SpawnRow int
	Direction          Direction
	Kind               TankKind
	LastFire           time.Time
}

func NewPlayer(col, row int, dir Direction) Tank {
	return Tank{Col: col, Row: row, SpawnCol: col, SpawnRow: row, Direction: dir, Kind: Player}
}

var (
	rowDelta = [4]int{-1, 1, 0, 0}
	colDelta = [4]int{0, 0, -1, 1}
)

func TryMoveTank(m Map, t *Tank) {
	newRow := t.Row + rowDelta[t.Direction]
	newCol := t.Col + colDelta[t.Direction]

	if newRow < 0 || newRow >= MapSize || newCol < 0 || newCol >= MapSize {
		return
	}
	if !tilePassable(m.Grid[newRow][newCol]) {
		return
	}

	t.Row = newRow
	t.Col = newCol
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
