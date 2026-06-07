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

	switch os.Args[1] {
	case "host":
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandHost, nil))
	case "join":
		if len(os.Args) < 3 {
			log.Fatal("join command requires room id")
		}
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.CommandJoin, []byte(os.Args[2])))
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

	state := game.GameState{
		Map:  game.Map1,
		Tank: game.NewPlayer(13, 22, game.Up),
	}
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
			game.Render(screen, state)
		case msg, ok := <-netMessages:
			if !ok {
				netMessages = nil
				break
			}
			if len(msg) == 0 {
				break
			}
			switch msg[0] {
			case protocol.MsgRoomCode:
				state.RoomCode = string(msg[1:])
			case protocol.MsgStart:
				enemy := game.NewEnemy(12, 3, game.Down)
				state.Enemy = &enemy
			}
			game.Render(screen, state)
		case ev := <-screen.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					return
				}
				handleInput(ev, &state)
				game.Render(screen, state)
			case *tcell.EventResize:
				screen.Sync()
			}
		}
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
