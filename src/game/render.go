package game

import "github.com/nsf/termbox-go"

var (
	defaultColor   = termbox.ColorDefault
	boardStartPosX = 3
	boardStartPosY = 2
)

func (l *level) render() {
	termbox.Clear(defaultColor, defaultColor)

	renderBoard(l, boardStartPosX, boardStartPosY)

	termbox.Flush()
}

func renderBoard(lev *level, startX, startY int) {
	x, y := startX, startY
	for _, row := range lev.board {
		for _, candy := range row {
			setBg(x, y, candy.color)
			x += 4
		}
		y += 2
		x = startX
	}
}
