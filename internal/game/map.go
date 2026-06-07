package game

import "github.com/gdamore/tcell/v3"

const MapSize = 26

type Map struct {
	Name string
	Grid [MapSize][MapSize]TileType
}

func MapOffset(screen tcell.Screen) (int, int) {
	w, h := screen.Size()
	return (w - MapSize*2) / 2, (h - MapSize) / 2
}

func RenderMap(screen tcell.Screen, m Map) {
	ox, oy := MapOffset(screen)
	for row := range m.Grid {
		for col := range m.Grid[row] {
			char, style := tileGlyph(m.Grid[row][col])
			screen.SetContent(ox+col*2, oy+row, char, nil, style)
			screen.SetContent(ox+col*2+1, oy+row, char, nil, style)
		}
	}
}

var Map1 = Map{
	Name: "Classic",
	Grid: [MapSize][MapSize]TileType{
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Base, Base, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Base, Base, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	},
}

var Map2 = Map{
	Name: "Fortress",
	Grid: [MapSize][MapSize]TileType{
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Base, Base, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Water, Water, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Water, Water, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Water, Water, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Water, Water, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty, Empty, Empty, Brick, Brick, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Brick, Brick, Empty, Empty, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Empty, Empty, Brick, Brick, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Steel, Steel, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Steel, Base, Base, Steel, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	},
}
