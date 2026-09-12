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
				handleInput(ev, &state)
			case *tcell.EventResize:
				screen.Sync()
			}
		}
		if state.LevelOver {
			game.NextLevel(&state, state.Winner)
		}
		game.Render(screen, state)
	}
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
	}
}

func handleInput(key *tcell.EventKey, state *game.GameState) {
	if key.Key() == tcell.KeyRune && key.Str() == " " {
		game.FireBullet(state)
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
	game.TryMoveTank(state.Map, &state.Tank)
}
