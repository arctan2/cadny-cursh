package game

import "github.com/nsf/termbox-go"

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
