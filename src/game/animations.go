package game

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

var prevCellColor termbox.Attribute

func setBg(x, y int, candyColor termbox.Attribute) {
	termbox.SetBg(x, y, candyColor)
	termbox.SetBg(x+1, y, candyColor)
}

/*
[
	[a, b, c],
	[p, q, r],
	[x, y, z]
]

x -> -1, p -> -3, a -> -5
x ->  2, p -> -2, a -> -4
x ->  3, p -> -1, a -> -3
x ->  4, p ->  2, a -> -2
x ->  5, p ->  3, a -> -1
x ->  6, p ->  4, a ->  2
*/

func fallAnimation(lev *level, colIdx, paintIdx int, wg *sync.WaitGroup, mut *sync.Mutex) {
	var rowCount int = len(lev.board)
	candiesPosY := make([]int, len(lev.board))
	for i, yPos := 0, (rowCount*2-1)*-1; i < rowCount; i, yPos = i+1, yPos+2 {
		candiesPosY[i] = yPos
	}

	for iterCount := 0; iterCount < rowCount*2; iterCount++ {
		for i := rowCount - 1; i >= 0; i-- {
			setBg(paintIdx, candiesPosY[i], lev.board[i][colIdx].color)
			setBg(paintIdx, candiesPosY[i]-1, defaultColor)

			candiesPosY[i]++
			if candiesPosY[i] == 0 {
				candiesPosY[i] = lev.posY
			}
		}
		mut.Lock()
		termbox.Flush()
		mut.Unlock()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}
	wg.Done()
}

func initBoardAnimation(lev *level) {
	var wg sync.WaitGroup
	var mut sync.Mutex

	paintIdx := lev.posX
	for colIdx := range lev.board[0] {
		wg.Add(1)
		go fallAnimation(lev, colIdx, paintIdx, &wg, &mut)
		paintIdx += 4
	}
	wg.Wait()
}

func blinkCursor(lev *level) {
	for {
		if lev.board[lev.cursor.y][lev.cursor.x].color == defaultColor {
			lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
		} else {
			prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
			lev.board[lev.cursor.y][lev.cursor.x].color = defaultColor
		}
		time.Sleep(time.Millisecond * 300)
	}
}
