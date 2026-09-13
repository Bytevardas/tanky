package main

import (
	"strings"
	"unicode/utf8"

	"tanky/internal/game"
	"tanky/internal/protocol"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type menuOption int

const (
	optionHost menuOption = iota
	optionJoin
	optionExit
)

var menuLabels = [3]string{optionHost: "Host game", optionJoin: "Join game", optionExit: "Exit"}

var banner = []string{
	"████████╗ █████╗ ███╗   ██╗██╗  ██╗██╗   ██╗",
	"╚══██╔══╝██╔══██╗████╗  ██║██║ ██╔╝╚██╗ ██╔╝",
	"   ██║   ███████║██╔██╗ ██║█████╔╝  ╚████╔╝ ",
	"   ██║   ██╔══██║██║╚██╗██║██╔═██╗   ╚██╔╝  ",
	"   ██║   ██║  ██║██║ ╚████║██║  ██╗   ██║   ",
	"   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝   ",
}

var (
	accentStyle = tcell.StyleDefault.Foreground(color.Green)
	textStyle   = tcell.StyleDefault.Foreground(color.White)
	hintStyle   = tcell.StyleDefault.Foreground(color.Gray)
)

func runMenu(screen tcell.Screen) (menuOption, string) {
	selected := optionHost
	for {
		renderMenu(screen, selected)
		key, ok := (<-screen.EventQ()).(*tcell.EventKey)
		if !ok {
			screen.Sync()
			continue
		}
		if dir, ok := keyDirection(key); ok {
			selected = moveSelection(selected, dir)
			continue
		}
		switch key.Key() {
		case tcell.KeyEscape:
			return optionExit, ""
		case tcell.KeyEnter:
			if selected != optionJoin {
				return selected, ""
			}
			if code, ok := promptRoomCode(screen); ok {
				return optionJoin, code
			}
		}
	}
}

func moveSelection(selected menuOption, dir game.Direction) menuOption {
	switch dir {
	case game.Up:
		return max(selected-1, optionHost)
	case game.Down:
		return min(selected+1, optionExit)
	}
	return selected
}

func promptRoomCode(screen tcell.Screen) (string, bool) {
	defer screen.HideCursor()
	code := ""
	for {
		renderPrompt(screen, code)
		key, ok := (<-screen.EventQ()).(*tcell.EventKey)
		if !ok {
			screen.Sync()
			continue
		}
		switch key.Key() {
		case tcell.KeyEscape:
			return "", false
		case tcell.KeyEnter:
			if len(code) == protocol.RoomCodeLength {
				return code, true
			}
		case tcell.KeyBackspace:
			if len(code) > 0 {
				code = code[:len(code)-1]
			}
		case tcell.KeyRune:
			if len(code) == protocol.RoomCodeLength || len(key.Str()) != 1 {
				continue
			}
			if strings.Contains(protocol.RoomCodeChars, key.Str()) {
				code += key.Str()
			}
		}
	}
}

func renderMenu(screen tcell.Screen, selected menuOption) {
	screen.Clear()
	y := renderBanner(screen)
	for option, label := range menuLabels {
		marker, style := "  ", textStyle
		if menuOption(option) == selected {
			marker, style = "> ", accentStyle
		}
		renderCentered(screen, y+option, marker+label, style)
	}
	renderCentered(screen, y+len(menuLabels)+1, "W/S or arrows to move, Enter to select, Esc to quit", hintStyle)
	screen.Show()
}

func renderPrompt(screen tcell.Screen, code string) {
	screen.Clear()
	y := renderBanner(screen)
	label := "Room code: "
	x := renderCentered(screen, y, label+code+strings.Repeat("_", protocol.RoomCodeLength-len(code)), textStyle)
	screen.ShowCursor(x+len(label)+len(code), y)
	renderCentered(screen, y+2, "Enter to join, Esc to go back", hintStyle)
	screen.Show()
}

func renderBanner(screen tcell.Screen) int {
	_, h := screen.Size()
	menuHeight := len(banner) + 1 + len(menuLabels)
	top := max((h-menuHeight)/2, 0)
	for i, line := range banner {
		renderCentered(screen, top+i, line, accentStyle)
	}
	return top + len(banner) + 1
}

func renderCentered(screen tcell.Screen, y int, s string, style tcell.Style) int {
	w, _ := screen.Size()
	x := max((w-utf8.RuneCountInString(s))/2, 0)
	game.RenderText(screen, x, y, s, style)
	return x
}
