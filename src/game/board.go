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
	case "down":
		lev.swapAnimation(0, 1)
	case "left":
		lev.swapAnimation(-2, 0)
	case "right":
		lev.swapAnimation(2, 0)
	}
}
