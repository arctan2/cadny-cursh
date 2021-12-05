package game

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	colors = map[int]termbox.Attribute{
		0: termbox.ColorBlue,
		1: termbox.ColorGreen,
		2: termbox.ColorYellow,
		3: termbox.ColorRed,
	}
)

func (lev *level) initBoard() {
	rand.Seed(time.Now().UnixNano())
	lev.board.generateRandomCandies()
	lev.board.checkAndReplaceMatches()
	for lev.board.hasNoPossibleMoves() {
		lev.board.generateRandomCandies()
		lev.board.checkAndReplaceMatches()
	}
	initBoardAnimation(lev)
	lev.startBlink()
}

func (lev *level) makeMove(dir string) {
	var xmag, ymag int
	var toAnalizeCoords []coord

	setxymag := func(x, y int) {
		xmag, ymag = x, y
	}

	switch dir {
	case "up":
		setxymag(0, -1)
		toAnalizeCoords = []coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x, lev.cursor.y - 1}}
	case "down":
		setxymag(0, 1)
		toAnalizeCoords = []coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x, lev.cursor.y + 1}}
	case "left":
		setxymag(-2, 0)
		toAnalizeCoords = []coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x - 1, lev.cursor.y}}
	case "right":
		setxymag(2, 0)
		toAnalizeCoords = []coord{{lev.cursor.x, lev.cursor.y}, {lev.cursor.x + 1, lev.cursor.y}}
	}

	lev.swapAnimation(xmag, ymag)
	if isMatched, affectedColumns := lev.board.analizeAtCoords(toAnalizeCoords); !isMatched {
		lev.swapAnimation(xmag, ymag)
	} else {
		lev.board.fillVacancies(affectedColumns)
		prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
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

func (b board) analizeAtCoords(coords []coord) (bool, between) {
	var verticRange, horizRange between
	var matched bool
	var color termbox.Attribute
	var affectedColumns between

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
		if affectedColumns.from > horizRange.from {
			affectedColumns.from = horizRange.from
		}
		if affectedColumns.to < horizRange.to {
			affectedColumns.to = horizRange.to
		}
		b.destroyMatches(horizRange, verticRange, c, &matched)
	}

	return matched, affectedColumns
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

func (b board) fallCandiesAndFillRandom(x int) {
	_, lastVacant := b.getMax()

	for ; b[lastVacant][x].color != defaultColor; lastVacant-- {
		if lastVacant == 0 {
			break
		}
	}

	lastCandy := lastVacant
	for ; b[lastCandy][x].color == defaultColor; lastCandy-- {
		if lastCandy == 0 {
			break
		}
	}

	for ; lastCandy >= 0; lastCandy, lastVacant = lastCandy-1, lastVacant-1 {
		b[lastVacant][x].color, b[lastCandy][x].color = b[lastCandy][x].color, b[lastVacant][x].color
	}

	for i := 0; b[i][x].color == defaultColor; i++ {
		b[i][x] = candy{colors[rand.Intn(4)]}
	}

}

func (b board) fillVacancies(btw between) {
	for i := btw.from; i <= btw.to; i++ {
		b.fallCandiesAndFillRandom(i)
	}
}
