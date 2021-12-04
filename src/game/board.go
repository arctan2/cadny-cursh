package game

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func (lev *level) initBoard() {
	rand.Seed(time.Now().UnixNano())
	colors := map[int]termbox.Attribute{
		0: termbox.ColorBlue,
		1: termbox.ColorGreen,
		2: termbox.ColorYellow,
		3: termbox.ColorRed,
	}

	lev.board.generateRandomCandies(colors)
	lev.board.checkAndReplaceMatches(colors)
	if lev.board.hasNoPossibleMoves() {
		for lev.board.hasNoPossibleMoves() {
			lev.board.generateRandomCandies(colors)
			lev.board.checkAndReplaceMatches(colors)
		}
	}
	initBoardAnimation(lev)
	lev.startBlink()
}

func (lev *level) makeMove(dir string) {
	switch dir {
	case "up":
		lev.swapAnimation(0, -1)
		if !lev.board.analizeAtCoords([]coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x, lev.cursor.y - 1}}) {
			lev.swapAnimation(0, -1)
		}
	case "down":
		lev.swapAnimation(0, 1)
		if !lev.board.analizeAtCoords([]coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x, lev.cursor.y + 1}}) {
			lev.swapAnimation(0, 1)
		}
	case "left":
		lev.swapAnimation(-2, 0)
		if !lev.board.analizeAtCoords([]coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x - 1, lev.cursor.y}}) {
			lev.swapAnimation(-2, 0)
		}
	case "right":
		lev.swapAnimation(2, 0)
		if lev.board.analizeAtCoords([]coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x + 1, lev.cursor.y}}) {
			lev.swapAnimation(2, 0)
		}
	}
}

type between struct {
	from, to int
}

func (b board) loopHorizTwice(x, y, mag int, color termbox.Attribute) int {
	var val int = x
	for i, newx := 1, 0; i <= 2; i++ {
		newx = x + (i * mag)
		if !b.xyNotInBounds(newx, y) {
			if b[y][newx].color == color {
				val = newx
			} else {
				break
			}
		}
	}
	return val
}

func (b board) loopVerticTwice(x, y, mag int, color termbox.Attribute) int {
	var val int = y
	for i, newy := 1, 0; i <= 2; i++ {
		newy = y + (i * mag)
		if !b.xyNotInBounds(x, newy) {
			if b[newy][x].color == color {
				val = newy
			} else {
				break
			}
		}
	}
	return val
}

func (b board) analizeAtCoords(coords []coord) bool {
	var verticRange, horizRange between
	var matched bool
	var color termbox.Attribute
	for _, c := range coords {
		color = b[c.y][c.x].color
		horizRange = between{
			from: b.loopHorizTwice(c.x, c.y, -1, color),
			to:   b.loopHorizTwice(c.x, c.y, 1, color),
		}
		verticRange = between{
			from: b.loopVerticTwice(c.x, c.y, -1, color),
			to:   b.loopVerticTwice(c.x, c.y, 1, color),
		}
		b.destroyMatches(horizRange, verticRange, c, &matched)
	}

	return matched
}

func (b board) destroyMatches(horizRange, verticRange between, c coord, matched *bool) {
	if horizRange.to-horizRange.from > 1 {
		for i := horizRange.from; i <= horizRange.to; i++ {
			b[c.y][i].color = defaultColor
		}
		*matched = true
	}
	if verticRange.to-verticRange.from > 1 {
		for i := verticRange.from; i <= verticRange.to; i++ {
			b[i][c.x].color = defaultColor
		}
		*matched = true
	}
}
