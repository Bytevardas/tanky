package game

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type GameState struct {
	Map       Map
	Tank      Tank
	Enemy     *Tank
	Bullets   []Bullet
	RoomCode  string
	Level     int
	Score     [2]int
	LevelOver bool
	Winner    Side
}

func NewGameState(side Side) GameState {
	return GameState{Map: Levels[0], Tank: NewTank(side, Player)}
}

func SpawnEnemy(state *GameState) {
	enemy := NewTank(state.Tank.Side.opposite(), Enemy)
	state.Enemy = &enemy
}

func NextLevel(state *GameState, winner Side) {
	state.Score[winner]++
	state.Level++
	state.Map = Levels[state.Level%len(Levels)]
	state.Bullets = state.Bullets[:0]
	state.LevelOver = false
	respawn(&state.Tank)
	if state.Enemy != nil {
		respawn(state.Enemy)
	}
}

func Render(screen tcell.Screen, state GameState) {
	screen.Clear()
	RenderMap(screen, state.Map)
	renderHUD(screen, state)
	RenderTank(screen, state.Tank)
	if state.Enemy != nil {
		RenderTank(screen, *state.Enemy)
	}
	for _, b := range state.Bullets {
		RenderBullet(screen, b)
	}
	screen.Show()
}

func renderHUD(screen tcell.Screen, state GameState) {
	ox, oy := MapOffset(screen)
	style := tcell.StyleDefault.Foreground(color.White)
	if state.RoomCode != "" {
		RenderText(screen, ox, max(oy-2, 0), "Room: "+state.RoomCode, style)
	}
	you := state.Score[state.Tank.Side]
	enemy := state.Score[state.Tank.Side.opposite()]
	level := fmt.Sprintf("Level %d: %s   You %d - %d Enemy", state.Level+1, state.Map.Name, you, enemy)
	RenderText(screen, ox, max(oy-1, 0), level, style)
}

func RenderText(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for i, r := range []rune(s) {
		screen.SetContent(x+i, y, r, nil, style)
	}
}
