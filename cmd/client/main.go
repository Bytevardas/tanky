package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"tanky/internal/game"
	"tanky/internal/protocol"

	"github.com/gdamore/tcell/v3"
)

func main() {
	reason := run(os.Args[1:])
	if reason != "" {
		fmt.Println(reason)
	}
}

func run(args []string) string {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal("failed to create new screen")
	}
	if err := screen.Init(); err != nil {
		log.Fatal("failed to create new screen")
	}
	defer screen.Fini()

	var option menuOption
	var code string
	switch {
	case len(args) == 0:
		option, code = runMenu(screen)
	case args[0] == "host":
		option = optionHost
	case args[0] == "join" && len(args) > 1:
		option, code = optionJoin, args[1]
	default:
		return "usage: client [host | join <code>]"
	}
	if option == optionExit {
		return ""
	}

	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		return "Could not reach server: " + err.Error()
	}
	defer conn.Close()

	var state game.GameState
	switch option {
	case optionHost:
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandHost, nil))
		state = game.NewGameState(game.Bottom)
	case optionJoin:
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandJoin, []byte(code)))
		state = game.NewGameState(game.Top)
	}
	return play(screen, conn, state)
}

func play(screen tcell.Screen, conn net.Conn, state game.GameState) string {
	game.Render(screen, state)

	netMessages := make(chan []byte)
	go func() {
		for {
			msg, err := protocol.ReadMessage(conn)
			if err != nil {
				close(netMessages)
				return
			}
			netMessages <- msg
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		before := state.Tank
		select {
		case <-ticker.C:
			game.UpdateBullets(&state)
		case msg, ok := <-netMessages:
			if !ok {
				return "Connection lost"
			}
			if reason := handleMessage(msg, &state); reason != "" {
				return reason
			}
		case ev := <-screen.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					return ""
				}
				handleInput(ev, conn, &state)
			case *tcell.EventResize:
				screen.Sync()
			}
		}
		if state.LevelOver {
			finishLevel(conn, &state)
		}
		if poseChanged(before, state.Tank) {
			t := state.Tank
			protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgMove, encodePose(t.Col, t.Row, t.Direction)))
		}
		game.Render(screen, state)
	}
}

func finishLevel(conn net.Conn, state *game.GameState) {
	winner := state.Winner
	game.NextLevel(state, winner)
	protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgNextLevel, []byte{byte(state.Level), byte(winner)}))
}

func handleMessage(msg []byte, state *game.GameState) string {
	if len(msg) == 0 {
		return ""
	}
	switch msg[0] {
	case protocol.MsgError:
		return string(msg[1:])
	case protocol.MsgPeerLeft:
		return "Opponent left"
	case protocol.MsgRoomCode:
		state.RoomCode = string(msg[1:])
	case protocol.MsgStart:
		game.SpawnEnemy(state)
	case protocol.MsgMove:
		if state.Enemy == nil {
			return ""
		}
		col, row, dir, ok := decodePose(msg[1:])
		if !ok {
			return ""
		}
		state.Enemy.Col, state.Enemy.Row, state.Enemy.Direction = col, row, dir
		game.AbsorbBullets(state, *state.Enemy)
	case protocol.MsgFire:
		if state.Enemy == nil {
			return ""
		}
		col, row, dir, ok := decodePose(msg[1:])
		if !ok {
			return ""
		}
		game.FireFrom(state, game.Tank{Col: col, Row: row, Direction: dir, Side: state.Enemy.Side})
	case protocol.MsgNextLevel:
		if len(msg) < 3 || msg[2] > byte(game.Top) {
			return ""
		}
		if int(msg[1]) <= state.Level {
			return ""
		}
		game.NextLevel(state, game.Side(msg[2]))
	}
	return ""
}

func handleInput(key *tcell.EventKey, conn net.Conn, state *game.GameState) {
	if key.Key() == tcell.KeyRune && key.Str() == " " {
		if !game.FireBullet(state) {
			return
		}
		t := state.Tank
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgFire, encodePose(t.Col, t.Row, t.Direction)))
		return
	}

	dir, ok := keyDirection(key)
	if !ok {
		return
	}
	state.Tank.Direction = dir
	game.TryMoveTank(state)
}

func keyDirection(key *tcell.EventKey) (game.Direction, bool) {
	switch key.Key() {
	case tcell.KeyUp:
		return game.Up, true
	case tcell.KeyDown:
		return game.Down, true
	case tcell.KeyLeft:
		return game.Left, true
	case tcell.KeyRight:
		return game.Right, true
	case tcell.KeyRune:
		switch key.Str() {
		case "w", "W":
			return game.Up, true
		case "s", "S":
			return game.Down, true
		case "a", "A":
			return game.Left, true
		case "d", "D":
			return game.Right, true
		}
	}
	return 0, false
}

func poseChanged(before, after game.Tank) bool {
	return before.Col != after.Col || before.Row != after.Row || before.Direction != after.Direction
}

func encodePose(col, row int, dir game.Direction) []byte {
	return []byte{byte(col), byte(row), byte(dir)}
}

func decodePose(b []byte) (col, row int, dir game.Direction, ok bool) {
	if len(b) < 3 {
		return 0, 0, 0, false
	}
	col, row, dir = int(b[0]), int(b[1]), game.Direction(b[2])
	if col >= game.MapSize || row >= game.MapSize || dir > game.Right {
		return 0, 0, 0, false
	}
	return col, row, dir, true
}
