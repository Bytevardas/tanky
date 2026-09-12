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

var availableCommands = []string{"host", "join", "help"}

func main() {
	fmt.Println("staring client")

	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if len(os.Args) < 2 {
		log.Fatal("expecting command to be passed in: host or join <code>")
	}

	var state game.GameState
	switch os.Args[1] {
	case "host":
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandHost, nil))
		state = game.NewGameState(game.Bottom)
	case "join":
		if len(os.Args) < 3 {
			log.Fatal("join command requires room id")
		}
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandJoin, []byte(os.Args[2])))
		state = game.NewGameState(game.Top)
	default:
		log.Fatal("unknown command: expecting host or join <code>")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal("failed to create new screen")
	}

	err = screen.Init()
	if err != nil {
		log.Fatal("failed to create new screen")
	}
	defer screen.Fini()

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
				netMessages = nil
				break
			}
			handleMessage(msg, &state)
		case ev := <-screen.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					return
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

func handleMessage(msg []byte, state *game.GameState) {
	if len(msg) == 0 {
		return
	}
	switch msg[0] {
	case protocol.MsgRoomCode:
		state.RoomCode = string(msg[1:])
	case protocol.MsgStart:
		game.SpawnEnemy(state)
	case protocol.MsgMove:
		if state.Enemy == nil {
			return
		}
		col, row, dir, ok := decodePose(msg[1:])
		if !ok {
			return
		}
		state.Enemy.Col, state.Enemy.Row, state.Enemy.Direction = col, row, dir
		game.AbsorbBullets(state, *state.Enemy)
	case protocol.MsgFire:
		if state.Enemy == nil {
			return
		}
		col, row, dir, ok := decodePose(msg[1:])
		if !ok {
			return
		}
		game.FireFrom(state, game.Tank{Col: col, Row: row, Direction: dir, Side: state.Enemy.Side})
	case protocol.MsgNextLevel:
		if len(msg) < 3 || msg[2] > byte(game.Top) {
			return
		}
		if int(msg[1]) <= state.Level {
			return
		}
		game.NextLevel(state, game.Side(msg[2]))
	}
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

	var dir game.Direction
	switch key.Key() {
	case tcell.KeyUp:
		dir = game.Up
	case tcell.KeyDown:
		dir = game.Down
	case tcell.KeyLeft:
		dir = game.Left
	case tcell.KeyRight:
		dir = game.Right
	case tcell.KeyRune:
		switch key.Str() {
		case "w", "W":
			dir = game.Up
		case "s", "S":
			dir = game.Down
		case "a", "A":
			dir = game.Left
		case "d", "D":
			dir = game.Right
		default:
			return
		}
	default:
		return
	}
	state.Tank.Direction = dir
	game.TryMoveTank(state)
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
