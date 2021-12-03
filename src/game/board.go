package game

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type candy struct {
	color termbox.Attribute
}

func (lev *level) generateRandomCandies(colors map[int]termbox.Attribute) {
	for i := 0; i < len(lev.board); i++ {
		for j := 0; j < len(lev.board[i]); j++ {
			lev.board[i][j] = candy{colors[rand.Intn(4)]}
		}
	}
}

func (lev *level) scanAndReplaceMatches(y int, colors map[int]termbox.Attribute, wg *sync.WaitGroup) {
	for x := range lev.board[y] {
		color := lev.board[y][x].color
		ymax := len(lev.board) - 1
		xmax := len(lev.board[0]) - 1
		if y+1 <= ymax && y+2 <= ymax {
			if lev.board[y+1][x].color == color && lev.board[y+2][x].color == color {
				for lev.board[y+2][x].color == color {
					lev.board[y+2][x] = candy{colors[rand.Intn(4)]}
				}
			}
		}
		if x+1 <= xmax && x+2 <= xmax {
			if lev.board[y][x+1].color == color && lev.board[y][x+2].color == color {
				for lev.board[y][x+2].color == color {
					lev.board[y][x+2] = candy{colors[rand.Intn(4)]}
				}
			}
		}
	}
	wg.Done()
}

func (lev *level) checkAndReplaceMatches(colors map[int]termbox.Attribute) {
	var wg sync.WaitGroup
	for i := 0; i < len(lev.board); i++ {
		wg.Add(1)
		go lev.scanAndReplaceMatches(i, colors, &wg)
	}
	wg.Wait()
}

func (lev *level) initBoard() {
	rand.Seed(time.Now().UnixNano())
	colors := map[int]termbox.Attribute{
		0: termbox.ColorBlue,
		1: termbox.ColorGreen,
		2: termbox.ColorYellow,
		3: termbox.ColorRed,
	}

	lev.generateRandomCandies(colors)
	lev.checkAndReplaceMatches(colors)
	initBoardAnimation(lev)
	go lev.blinkCursor()
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
