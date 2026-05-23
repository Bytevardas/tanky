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

var passable = [7]bool{
	true,  // Empty
	false, // Brick
	false, // Steel
	true,  // Trees
	false, // Water
	true,  // Ice
	false, // Base
}

func tilePassable(t TileType) bool {
	return passable[t]
}

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
