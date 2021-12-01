package game

import (
	"cadny-cursh/src/utils"
	"time"

	"github.com/nsf/termbox-go"
)

type level struct {
	board [][]candy
	posX  int
	posY  int
}

func newLevel(rowCount, colCount, posX, posY int) *level {
	newBoard := make([][]candy, rowCount)
	for i := range newBoard {
		newBoard[i] = make([]candy, colCount)
	}
	l := level{newBoard, posX, posY}
	return &l
}

func handleKeyboardEvent(kEvent keyboardEvent) bool {
	switch kEvent.eventType {
	case END:
		return true
	}
	return false
}

func Start() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer func() {
		termbox.Clear(defaultColor, defaultColor)
		termbox.Flush()
		termbox.Close()
		utils.Clrscr()
	}()

	lev := newLevel(8, 8, 3, 2)

	lev.initBoard()

	var keyboardChan chan keyboardEvent = make(chan keyboardEvent)

	go listenToKeyboard(keyboardChan)

mainloop:
	for {
		select {
		case e := <-keyboardChan:
			if breakLoop := handleKeyboardEvent(e); breakLoop {
				break mainloop
			}
		default:
			lev.render()
			time.Sleep(time.Millisecond * 10)
		}
	}
}
