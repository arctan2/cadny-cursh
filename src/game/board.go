package game

import (
	"math/rand"
	"sync"
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
	lev.checkNoPossibleMoves(false)
	initBoardAnimation(lev)
	lev.recursiveDestroyCandies()
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
	if isMatched, affectedColumns := lev.board.analizeAtCoords(toAnalizeCoords, pos{lev.coordX, lev.coordY}); !isMatched {
		lev.swapAnimation(xmag, ymag)
	} else {
		lev.fillVacancies(affectedColumns)
		lev.recursiveDestroyCandies()
		prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
	}
	lev.checkNoPossibleMoves(true)
}

func (lev *level) checkNoPossibleMoves(animate bool) {
	animTime := time.Millisecond * 20
	for lev.board.hasNoPossibleMoves() {
		if animate {
			for y := range lev.board {
				for x := range lev.board[y] {
					setBg(lev.coordX(x), lev.coordY(y), defaultColor)
					termbox.Flush()
					time.Sleep(animTime)
				}
			}
		}
		lev.board.generateRandomCandies()
		if animate {
			for y := range lev.board {
				for x := range lev.board[y] {
					setBg(lev.coordX(x), lev.coordY(y), lev.board[y][x].color)
					termbox.Flush()
					time.Sleep(animTime)
				}
			}
		}
		lev.board.checkAndReplaceMatches()
	}
}

func (lev *level) recursiveDestroyCandies() {
	for y := range lev.board {
		for x := range lev.board[y] {
			if isMatched, affectedColumns := lev.board.analizeAtCoords([]coord{{x, y}}, pos{lev.coordX, lev.coordY}); isMatched {
				lev.fillVacancies(affectedColumns)
				lev.recursiveDestroyCandies()
				prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
			}
		}
	}
}

type pos struct {
	x, y func(int) int
}

type between struct {
	from, to int
}

func (b board) loopHoriz(x, y, mag int, color termbox.Attribute) int {
	var val int = x
	for i, newx := 1, 0; ; i++ {
		newx = x + (i * mag)
		if !b.xyNotInBounds(newx, y) {
			if b[y][newx].color == color {
				val = newx
			} else {
				break
			}
		} else {
			break
		}
	}
	return val
}

func (b board) loopVertic(x, y, mag int, color termbox.Attribute) int {
	var val int = y
	for i, newy := 1, 0; ; i++ {
		newy = y + (i * mag)
		if !b.xyNotInBounds(x, newy) {
			if b[newy][x].color == color {
				val = newy
			} else {
				break
			}
		} else {
			break
		}
	}
	return val
}

func (b board) analizeAtCoords(coords []coord, pos pos) (bool, between) {
	var verticRange, horizRange between
	var matched bool
	var color termbox.Attribute
	var affectedColumns between = between{-1, -1}

	for _, c := range coords {
		color = b[c.y][c.x].color
		horizRange = between{
			from: b.loopHoriz(c.x, c.y, -1, color),
			to:   b.loopHoriz(c.x, c.y, 1, color),
		}
		verticRange = between{
			from: b.loopVertic(c.x, c.y, -1, color),
			to:   b.loopVertic(c.x, c.y, 1, color),
		}

		b.destroyMatches(horizRange, verticRange, c, &matched, pos)
		if affectedColumns.from == -1 {
			affectedColumns.from = horizRange.from
			affectedColumns.to = horizRange.to
		} else if matched {
			if affectedColumns.from > horizRange.from || affectedColumns.from == -1 {
				affectedColumns.from = horizRange.from
			}
			if affectedColumns.to < horizRange.to || affectedColumns.to == -1 {
				affectedColumns.to = horizRange.to
			}
		}
	}

	return matched, affectedColumns
}

func (b board) destroyMatches(horizRange, verticRange between, c coord, matched *bool, pos pos) {
	if horizRange.to-horizRange.from > 1 {
		for i := horizRange.from; i <= horizRange.to; i++ {
			b[c.y][i].color = defaultColor
			setBg(pos.x(i), pos.y(c.y), defaultColor)
		}
		*matched = true
	}
	if verticRange.to-verticRange.from > 1 {
		for i := verticRange.from; i <= verticRange.to; i++ {
			b[i][c.x].color = defaultColor
			setBg(pos.x(c.x), pos.y(i), defaultColor)
		}
		*matched = true
	}
	termbox.Flush()
}

func (lev *level) fallCandiesAndFillRandom(x int) (int, int) {
	b := lev.board
	_, lastVacant := b.getMax()

	for ; b[lastVacant][x].color != defaultColor; lastVacant-- {
		if lastVacant == 0 {
			break
		}
	}

	lastCandy := lastVacant

	for ; b[lastCandy][x].color == defaultColor; lastCandy-- {
		if lastCandy == 0 {
			lastCandy--
			break
		}
	}

	for i, j := lastCandy, lastVacant; i >= 0; i, j = i-1, j-1 {
		b[j][x].color, b[i][x].color = b[i][x].color, b[j][x].color
	}

	for i := 0; b[i][x].color == defaultColor; i++ {
		b[i][x] = candy{colors[rand.Intn(4)]}
		if i+1 == len(lev.board) {
			break
		}
	}

	return lastVacant, lastCandy
}

func (b board) hasVacancies(x int) bool {
	for y := 0; y < len(b); y++ {
		if b[y][x].color == defaultColor {
			return true
		}
	}
	return false
}

func (lev *level) fillVacancies(btw between) {
	var wg sync.WaitGroup
	var mut sync.Mutex

	for i := btw.from; i <= btw.to; i++ {
		for lev.board.hasVacancies(i) {
			lv, lc := lev.fallCandiesAndFillRandom(i)
			if lv <= 0 {
				continue
			}
			wg.Add(1)
			go lev.fallAnim(lv, lc, i, &wg, &mut)
		}
	}
	wg.Wait()
}

func (lev *level) fallAnim(lv, lc, x int, wg *sync.WaitGroup, mut *sync.Mutex) {
	rowCount := lv + 1
	candiesPosY := make([]int, rowCount)
	paintIdx := lev.coordX(x)

	for yPos, candIdx := ((lc+2)*2 + 1), lv; candIdx >= 0; yPos, candIdx = yPos-2, candIdx-1 {
		if yPos < lev.posY {
			yPos = -2
			for ; candIdx >= 0; yPos, candIdx = yPos-2, candIdx-1 {
				candiesPosY[candIdx] = yPos
			}
			break
		}
		candiesPosY[candIdx] = yPos
	}

	fall(lev, candiesPosY, rowCount, x, paintIdx, 50, mut)
	wg.Done()
}
