package game

import (
	"cadny-cursh/src/utils"
	"time"

	"github.com/nsf/termbox-go"
)

type board [][]candy

type level struct {
	board                  board
	posX, posY, xmax, ymax int
	cursor                 cursor
	isSelected             bool
	blinkCh                chan bool
}

func (lev *level) startBlink() {
	go lev.blinkCursor(lev.blinkCh)
}

func (lev *level) stopBlink() {
	lev.blinkCh <- true
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
	l := level{
		board:      newBoard,
		posX:       posX,
		posY:       posY,
		xmax:       colCount - 1,
		ymax:       rowCount - 1,
		cursor:     cursor{},
		isSelected: false,
		blinkCh:    make(chan bool),
	}
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

func (lev *level) repaintCurCell() {
	if lev.board[lev.cursor.y][lev.cursor.x].color == defaultColor {
		lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
	}
}

func (lev *level) navigate(kEvent keyboardEvent) {
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
	adjacentColors.repaintCells(lev)
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
}

func (lev *level) handleKeyboardEvent(kEvent keyboardEvent, kbProc *keyboardEvProcess) bool {
	lev.stopBlink()
	lev.repaintCurCell()
	switch kEvent.eventType {
	case NAVIGATE:
		lev.navigate(kEvent)
	case SELECT:
		if selected := lev.toggleSelected(); selected {
			go lev.blinkAdjacent()
		} else {
			adjacentColors.repaintCells(lev)
		}
	case MOVE:
		kbProc.pause()
		lev.move(kEvent)
		kbProc.resume()
	case END:
		return true
	}
	lev.startBlink()
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

	lev := newLevel(5, 5, 8, 5)

	lev.initBoard()

	var keyboardChan chan keyboardEvent = make(chan keyboardEvent)

	var kbProc keyboardEvProcess = false

	go listenToKeyboard(&lev.isSelected, keyboardChan, &kbProc)

mainloop:
	for {
		select {
		case e := <-keyboardChan:
			if breakLoop := lev.handleKeyboardEvent(e, &kbProc); breakLoop {
				break mainloop
			}
		default:
			lev.render()
			time.Sleep(time.Millisecond * 10)
		}
	}
}
