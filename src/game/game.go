package game

import (
	"cadny-cursh/src/utils"
	"time"

	"github.com/nsf/termbox-go"
)

type level struct {
	board  [][]candy
	posX   int
	posY   int
	cursor cursor
}

type cursor struct {
	x, y int
}

func newLevel(rowCount, colCount, posX, posY int) *level {
	newBoard := make([][]candy, rowCount)
	for i := range newBoard {
		newBoard[i] = make([]candy, colCount)
	}
	l := level{newBoard, posX, posY, cursor{}}
	return &l
}

func validateCursor(curs *cursor, xmax, ymax int) {
	if curs.x < 0 {
		curs.x = xmax - 1
	} else if curs.x >= xmax {
		curs.x = 0
	}

	if curs.y < 0 {
		curs.y = ymax - 1
	} else if curs.y >= ymax {
		curs.y = 0
	}
}

func (lev *level) handleKeyboardEvent(kEvent keyboardEvent) bool {
	switch kEvent.eventType {
	case MOVE:
		if lev.board[lev.cursor.y][lev.cursor.x].color == defaultColor {
			lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
		}
		switch kEvent.key {
		case termbox.KeyArrowDown:
			lev.cursor.y++
		case termbox.KeyArrowUp:
			lev.cursor.y--
		case termbox.KeyArrowLeft:
			lev.cursor.x--
		case termbox.KeyArrowRight:
			lev.cursor.x++
		}
		validateCursor(&lev.cursor, len(lev.board[0]), len(lev.board))
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
			if breakLoop := lev.handleKeyboardEvent(e); breakLoop {
				break mainloop
			}
		default:
			lev.render()
			time.Sleep(time.Millisecond * 10)
		}
	}
}
