package game

import (
	"github.com/nsf/termbox-go"
)

type keyboardEventType int

// Keyboard events
const (
	END keyboardEventType = 1 + iota
	MOVE
)

type keyboardEvent struct {
	eventType keyboardEventType
	key       termbox.Key
}

func listenToKeyboard(evChan chan keyboardEvent) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft:
				fallthrough
			case termbox.KeyArrowDown:
				fallthrough
			case termbox.KeyArrowRight:
				fallthrough
			case termbox.KeyArrowUp:
				evChan <- keyboardEvent{eventType: MOVE, key: ev.Key}
			case termbox.KeyEsc:
				evChan <- keyboardEvent{eventType: END, key: ev.Key}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}