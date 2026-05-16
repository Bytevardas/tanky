package main

import (
	"log"
	"net"

	"github.com/gdamore/tcell/v3"
)

func main() {
	network, err := net.Listen("tcp", "8080")
	if err != nil {
		log.Fatal("failed to establish tcp")
	}

	network.Accept()

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal("failer to create new screen")
	}

	err = screen.Init()
	if err != nil {
		log.Fatal("failer to create new screen")
	}
	defer screen.Fini()

	style := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
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
	}
}
