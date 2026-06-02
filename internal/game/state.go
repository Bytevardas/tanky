package game

import "github.com/gdamore/tcell/v3"

type GameState struct {
	Map     Map
	Tank    Tank
	Bullets []Bullet
}

func Render(screen tcell.Screen, state GameState) {
	RenderMap(screen, state.Map)
	RenderTank(screen, state.Tank)
	for _, b := range state.Bullets {
		RenderBullet(screen, b)
	}
	screen.Show()
}
