package game

import (
	"cadny-cursh/src/utils"
	"time"

	"github.com/nsf/termbox-go"
)

type board [][]candy

type level struct {
	board      board
	posX       int
	posY       int
	cursor     cursor
	isSelected bool
}

type coord struct {
	x, y int
}
type cursor coord

func newLevel(rowCount, colCount, posX, posY int) *level {
	newBoard := make([][]candy, rowCount)
	for i := range newBoard {
		newBoard[i] = make([]candy, colCount)
	}
	l := level{newBoard, posX, posY, cursor{}, false}
	return &l
}

func (lev *level) validateCursor() {
	curs := &lev.cursor
	xmax, ymax := len(lev.board[0])-1, len(lev.board)-1

	if curs.x < 0 {
		curs.x = xmax
	} else if curs.x > xmax {
		curs.x = 0
	}

	if curs.y < 0 {
		curs.y = ymax
	} else if curs.y > ymax {
		curs.y = 0
	}
}

func (lev *level) navigate(kEvent keyboardEvent) {
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
	lev.validateCursor()
}

func (lev *level) toggleSelected() bool {
	lev.isSelected = !lev.isSelected
	return lev.isSelected
}

func (lev *level) move(kEvent keyboardEvent) {
	lev.isSelected = false

	for adj := range adjacentColors {
		lev.board[adjacentColors[adj].index.y][adjacentColors[adj].index.x].color = adjacentColors[adj].color
	}
	lev.render()
	switch kEvent.key {
	case termbox.KeyArrowDown:
		if lev.cursor.y+1 >= len(lev.board) {
			break
		}
		lev.makeMove("down")
	case termbox.KeyArrowUp:
		if lev.cursor.y-1 < 0 {
			break
		}
		lev.makeMove("up")
	case termbox.KeyArrowLeft:
		if lev.cursor.x-1 < 0 {
			break
		}
		lev.makeMove("left")
	case termbox.KeyArrowRight:
		if lev.cursor.x+1 >= len(lev.board[0]) {
			break
		}
		lev.makeMove("right")
	}
	go lev.blinkCursor()
}

func (lev *level) handleKeyboardEvent(kEvent keyboardEvent) bool {
	switch kEvent.eventType {
	case NAVIGATE:
		lev.navigate(kEvent)
	case SELECT:
		if selected := lev.toggleSelected(); selected {
			go lev.blinkAdjacent()
		} else {
			go lev.blinkCursor()
		}
	case MOVE:
		lev.move(kEvent)
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

	go listenToKeyboard(&lev.isSelected, keyboardChan)

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
