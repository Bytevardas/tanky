package game

import (
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

const fireCooldown = 500 * time.Millisecond

type Bullet struct {
	Col, Row  int
	Direction Direction
}

func RenderBullet(screen tcell.Screen, b Bullet) {
	ox, oy := MapOffset(screen)
	style := tcell.StyleDefault.Foreground(color.White)
	g := '|'
	if b.Direction == Left || b.Direction == Right {
		g = '-'
	}
	screen.SetContent(ox+b.Col*2, oy+b.Row, g, nil, style)
	screen.SetContent(ox+b.Col*2+1, oy+b.Row, g, nil, style)
}

func FireBullet(state *GameState) {
	if time.Since(state.Tank.LastFire) < fireCooldown {
		return
	}
	state.Tank.LastFire = time.Now()

	t := state.Tank
	b := Bullet{
		Col:       t.Col + colDelta[t.Direction],
		Row:       t.Row + rowDelta[t.Direction],
		Direction: t.Direction,
	}
	if bulletHit(state, b) {
		return
	}
	state.Bullets = append(state.Bullets, b)
}

func UpdateBullets(state *GameState) {
	write := 0
	for read := 0; read < len(state.Bullets); read++ {
		b := state.Bullets[read]
		b.Row += rowDelta[b.Direction]
		b.Col += colDelta[b.Direction]
		if bulletHit(state, b) {
			continue
		}
		state.Bullets[write] = b
		write++
	}
	state.Bullets = state.Bullets[:write]
}

func bulletHit(state *GameState, b Bullet) bool {
	if b.Row < 0 || b.Row >= MapSize || b.Col < 0 || b.Col >= MapSize {
		return true
	}
	tile := state.Map.Grid[b.Row][b.Col]
	if tile == Brick {
		state.Map.Grid[b.Row][b.Col] = Empty
		return true
	}
	if tile == Steel {
		return true
	}
	return false
}
