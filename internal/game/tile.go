package game

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type TileType uint8

const (
	Empty TileType = iota
	Brick
	Steel
	Trees
	Water
	Ice
	Base
)

func tileGlyph(t TileType) (rune, tcell.Style) {
	switch t {
	case Brick:
		return ' ', tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(200, 80, 20)).
			Background(tcell.NewRGBColor(80, 50, 30))
	case Steel:
		return ' ', tcell.StyleDefault.Background(color.Silver)
	case Trees:
		return ' ', tcell.StyleDefault.Background(color.Green)
	case Water:
		return ' ', tcell.StyleDefault.Background(color.Blue)
	case Ice:
		return ' ', tcell.StyleDefault.Background(color.LightBlue)
	case Base:
		return ' ', tcell.StyleDefault.Background(color.Gray)
	default:
		return ' ', tcell.StyleDefault
	}
}
