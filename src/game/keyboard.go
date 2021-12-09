package game

import (
	"github.com/nsf/termbox-go"
)

type keyboardEventType int

const (
	END keyboardEventType = 1 + iota
	NAVIGATE
	SELECT
	MOVE
)

type keyboardEvProcess bool

func (s *keyboardEvProcess) pause() {
	*s = true
}
func (s *keyboardEvProcess) resume() {
	*s = false
}

type keyboardEvent struct {
	eventType keyboardEventType
	key       termbox.Key
}

func wasdControls(isSelected bool, evChan chan keyboardEvent, stopKeyBoardEvent keyboardEvProcess, ch rune) {
	if stopKeyBoardEvent {
		return
	}
	var k termbox.Key

	switch ch {
	case 'w':
		fallthrough
	case 'W':
		k = termbox.KeyArrowUp
	case 'a':
		fallthrough
	case 'A':
		k = termbox.KeyArrowLeft
	case 's':
		fallthrough
	case 'S':
		k = termbox.KeyArrowDown
	case 'd':
		fallthrough
	case 'D':
		k = termbox.KeyArrowRight
	}

	if isSelected {
		evChan <- keyboardEvent{eventType: MOVE, key: k}
	} else {
		evChan <- keyboardEvent{eventType: NAVIGATE, key: k}
	}
}

func listenToKeyboard(isSelected *bool, evChan chan keyboardEvent, stopKeyBoardEvent *keyboardEvProcess) {
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
				if *stopKeyBoardEvent {
					break
				}
				if *isSelected {
					evChan <- keyboardEvent{eventType: MOVE, key: ev.Key}
				} else {
					evChan <- keyboardEvent{eventType: NAVIGATE, key: ev.Key}
				}
			case termbox.KeyEsc:
				evChan <- keyboardEvent{eventType: END, key: ev.Key}
			case termbox.KeySpace:
				if *stopKeyBoardEvent {
					break
				}
				evChan <- keyboardEvent{eventType: SELECT}
			default:
				if gameOver {
					evChan <- keyboardEvent{}
				} else {
					wasdControls(*isSelected, evChan, *stopKeyBoardEvent, ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
