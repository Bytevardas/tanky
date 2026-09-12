package game

import (
	"slices"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

const fireCooldown = 500 * time.Millisecond

type Bullet struct {
	Col, Row  int
	Direction Direction
	Owner     Side
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

func FireBullet(state *GameState) bool {
	if time.Since(state.Tank.LastFire) < fireCooldown {
		return false
	}
	state.Tank.LastFire = time.Now()

	FireFrom(state, state.Tank)
	return true
}

func FireFrom(state *GameState, shooter Tank) {
	b := Bullet{
		Col:       shooter.Col + colDelta[shooter.Direction],
		Row:       shooter.Row + rowDelta[shooter.Direction],
		Direction: shooter.Direction,
		Owner:     shooter.Side,
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

// Each client only respawns its own tank; the peer learns of it through the
// next move message, so an enemy hit just consumes the bullet here.
func bulletHit(state *GameState, b Bullet) bool {
	if b.Row < 0 || b.Row >= MapSize || b.Col < 0 || b.Col >= MapSize {
		return true
	}
	if hitsTank(b, state.Tank) {
		respawn(&state.Tank)
		return true
	}
	if state.Enemy != nil && hitsTank(b, *state.Enemy) {
		return true
	}
	switch state.Map.Grid[b.Row][b.Col] {
	case Brick:
		state.Map.Grid[b.Row][b.Col] = Empty
		return true
	case Base:
		state.LevelOver = true
		state.Winner = baseSide(b.Row).opposite()
		return true
	case Steel:
		return true
	}
	return false
}

func hitsTank(b Bullet, t Tank) bool {
	return b.Owner != t.Side && b.Row == t.Row && b.Col == t.Col
}

func AbsorbBullets(state *GameState, t Tank) bool {
	before := len(state.Bullets)
	state.Bullets = slices.DeleteFunc(state.Bullets, func(b Bullet) bool {
		return hitsTank(b, t)
	})
	return len(state.Bullets) < before
}
