package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"tanky/internal/protocol"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
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

	style := tcell.StyleDefault.Foreground(color.Green).Bold(true)
	screen.SetContent(5, 3, 'H', nil, style)
	screen.Show()

	for {
		ev := <-screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return
			}
		case *tcell.EventResize:
			screen.Sync()
		}

		b, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
