package game

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

var (
	defaultColor         = termbox.ColorDefault
	pointsDisplayFrom    = coord{1, 2}
	msgDisplayFrom       = coord{25, 2}
	movesLeftDisplayFrom = coord{1, 3}
)

func (lev *level) render() {
	termbox.Clear(defaultColor, defaultColor)

	renderBoard(lev, lev.posX, lev.posY)
	renderPoints(lev.visuals.points, pointsDisplayFrom)
	renderMsg(lev.visuals.msg, msgDisplayFrom)
	renderMovesLeft(lev.movesLeft, movesLeftDisplayFrom)

	termbox.Flush()
}

func renderMsg(msg string, pos coord) {
	renderText(pos.x, pos.y, msg)
}

func renderPoints(points int, pos coord) {
	x := pos.x
	text := "points: " + strconv.Itoa(points)
	renderText(x, pos.y, text)
}

func renderMovesLeft(movesLeft int, pos coord) {
	x := pos.x
	text := "moves: " + strconv.Itoa(movesLeft)
	renderText(x, pos.y, text)
}

func renderText(x, y int, text string) {
	for i := range text {
		termbox.SetChar(x, y, rune(text[i]))
		x++
	}
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

func renderGameOver(points int) {
	termbox.Clear(defaultColor, defaultColor)
	renderText(5, 5, "Game Over!")
	renderText(5, 6, "score: "+strconv.Itoa(points))
	renderText(5, 10, "press any key to exit")
	termbox.Flush()
}
