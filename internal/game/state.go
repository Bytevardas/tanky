package game

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type GameState struct {
	Map      Map
	Tank     Tank
	Enemy    *Tank
	Bullets  []Bullet
	RoomCode string
}

func Render(screen tcell.Screen, state GameState) {
	RenderMap(screen, state.Map)
	if state.RoomCode != "" {
		ox, oy := MapOffset(screen)
		y := oy - 2
		if y < 0 {
			y = 0
		}
		renderText(screen, ox, y, "Room: "+state.RoomCode, tcell.StyleDefault.Foreground(color.White))
	}
	RenderTank(screen, state.Tank)
	if state.Enemy != nil {
		RenderTank(screen, *state.Enemy)
	}
	for _, b := range state.Bullets {
		RenderBullet(screen, b)
	}
	screen.Show()
}

func renderText(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for i, r := range s {
		screen.SetContent(x+i, y, r, nil, style)
	}
}
