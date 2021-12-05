package game

import (
	"math/rand"
	"sync"

	"github.com/nsf/termbox-go"
)

type candy struct {
	color termbox.Attribute
}

func (b board) generateRandomCandies() {
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

func (b board) testPossibles(possibles possibleType, color termbox.Attribute, testFunc matchFunc, ch chan bool) {
	for _, possible := range possibles {
		if testFunc(possible.x, possible.y, color) {
			ch <- true
			return
		}
	}
	ch <- false
}

func (b board) hasNoPossibleMoves() bool {
	for y := range b {
		for x := range b[y] {
			hChan := make(chan bool)
			vChan := make(chan bool)
			eqChan := make(chan bool)

			go b.testPossibles(getHorizPossibles(x, y), b[y][x].color, b.hasHorizMatchFrom, hChan)
			go b.testPossibles(getVerticPossibles(x, y), b[y][x].color, b.hasVericMatchFrom, vChan)
			go b.anyOfGroupMatch(getMidPossibles(x, y), b[y][x].color, eqChan)

			if <-hChan || <-vChan || <-eqChan {
				return false
			}
		}
	}
	return true
}

func (b board) xyNotInBounds(x, y int) bool {
	xmax, ymax := b.getMax()
	if y > ymax || y < 0 || x > xmax || x < 0 {
		return true
	}
	return false
}

func (b board) hasHorizMatchFrom(x, y int, color termbox.Attribute) bool {
	if b.xyNotInBounds(x, y) {
		return false
	}

	xmax, _ := b.getMax()
	if x+1 <= xmax {
		if b[y][x].color == color && b[y][x+1].color == color {
			return true
		}
	}
	return false
}

func (b board) hasVericMatchFrom(x, y int, color termbox.Attribute) bool {
	if b.xyNotInBounds(x, y) {
		return false
	}
	_, ymax := b.getMax()
	if y+1 <= ymax {
		if b[y][x].color == color && b[y+1][x].color == color {
			return true
		}
	}
	return false
}

func (b board) anyOfGroupMatch(groups []possibleType, color termbox.Attribute, eqChan chan bool) {
	for _, group := range groups {
		s := areAllSame(b, group, color)
		if s {
			eqChan <- true
		}
	}
	eqChan <- false
}

func areAllSame(b board, group []coord, color termbox.Attribute) bool {
	for _, c := range group {
		if b.xyNotInBounds(c.x, c.y) {
			return false
		}
		if b[c.y][c.x].color != color {
			return false
		}
	}
	return true
}

func (b board) scanAndReplaceMatches(y int, wg *sync.WaitGroup, mut *sync.Mutex) {
	for x := range b[y] {
		color := b[y][x].color
		if b.hasVericMatchFrom(x, y+1, color) {
			mut.Lock()
			for b[y+2][x].color == color {
				b[y+2][x] = candy{colors[rand.Intn(4)]}
			}
			mut.Unlock()
		}
		if b.hasHorizMatchFrom(x+1, y, color) {
			mut.Lock()
			for b[y][x+2].color == color {
				b[y][x+2] = candy{colors[rand.Intn(4)]}
			}
			mut.Unlock()
		}
	}
	wg.Done()
}

func (b board) checkAndReplaceMatches() {
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < len(b); i++ {
		wg.Add(1)
		go b.scanAndReplaceMatches(i, &wg, &mut)
	}
	wg.Wait()
}
