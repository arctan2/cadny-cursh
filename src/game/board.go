package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type candy struct {
	color termbox.Attribute
}

func (b board) generateRandomCandies(colors map[int]termbox.Attribute) {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			b[i][j] = candy{colors[rand.Intn(4)]}
		}
	}
}

func (b board) getMax() (int, int) {
	return len(b[0]) - 1, len(b) - 1
}

type matchFunc func(int, int, termbox.Attribute) bool

func (b board) testPossibles(possibles possibleType, color termbox.Attribute, testFunc matchFunc, ch chan bool, mut *sync.Mutex) {
	mut.Lock()
	for _, possible := range possibles {
		if testFunc(possible.x, possible.y, color) {
			ch <- true
			mut.Unlock()
			return
		}
	}
	mut.Unlock()
	ch <- false
}

func (b board) hasNoPossibleMoves() bool {
	for y := range b {
		for x := range b[y] {
			horizPossibles := getHorizPossibles(x, y)
			verticPossibles := getVerticPossibles(x, y)

			hChan := make(chan bool)
			vChan := make(chan bool)
			var mut sync.Mutex

			go b.testPossibles(horizPossibles, b[y][x].color, b.hasHorizMatchFrom, hChan, &mut)
			go b.testPossibles(verticPossibles, b[y][x].color, b.hasVericMatchFrom, vChan, &mut)

			if hRes, vRes := <-hChan, <-vChan; hRes || vRes {
				return false
			}
		}
	}
	return true
}

func (b board) hasHorizMatchFrom(x, y int, color termbox.Attribute) bool {
	xmax, ymax := b.getMax()
	if y > ymax || y < 0 {
		return false
	}
	if x+1 <= xmax && x+2 <= xmax && x >= 0 {
		if b[y][x+1].color == color && b[y][x+2].color == color {
			return true
		}
	}
	return false
}

func (b board) hasVericMatchFrom(x, y int, color termbox.Attribute) bool {
	xmax, ymax := b.getMax()
	if x > xmax || x < 0 {
		return false
	}
	if y+1 <= ymax && y+2 <= ymax && y >= 0 {
		if b[y+1][x].color == color && b[y+2][x].color == color {
			return true
		}
	}
	return false
}

func (b board) scanAndReplaceMatches(y int, colors map[int]termbox.Attribute, wg *sync.WaitGroup, mut *sync.Mutex) {
	for x := range b[y] {
		color := b[y][x].color
		if b.hasVericMatchFrom(x, y, color) {
			mut.Lock()
			for b[y+2][x].color == color {
				b[y+2][x] = candy{colors[rand.Intn(4)]}
			}
			mut.Unlock()
		}
		if b.hasHorizMatchFrom(x, y, color) {
			mut.Lock()
			for b[y][x+2].color == color {
				b[y][x+2] = candy{colors[rand.Intn(4)]}
			}
			mut.Unlock()
		}
	}
	wg.Done()
}

func (b board) checkAndReplaceMatches(colors map[int]termbox.Attribute) {
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < len(b); i++ {
		wg.Add(1)
		go b.scanAndReplaceMatches(i, colors, &wg, &mut)
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

	lev.board.generateRandomCandies(colors)
	lev.board.checkAndReplaceMatches(colors)
	fmt.Println(lev.board.hasNoPossibleMoves())
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
