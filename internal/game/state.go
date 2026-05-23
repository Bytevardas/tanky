package game

import "github.com/gdamore/tcell/v3"

type GameState struct {
	Map  Map
	Tank Tank
}

func Render(screen tcell.Screen, state GameState) {
	RenderMap(screen, state.Map)
	RenderTank(screen, state.Tank)
	screen.Show()
}
