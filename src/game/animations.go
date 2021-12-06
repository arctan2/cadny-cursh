package game

import (
	"math"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

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

func fall(lev *level, candiesPosY []int, rowCount, x, paintIdx int, duration float64, mut *sync.Mutex) {
	for iterCount := candiesPosY[0]; iterCount < 1; iterCount++ {
		for i := rowCount - 1; i >= 0; i-- {
			setBg(paintIdx, candiesPosY[i], lev.board[i][x].color)
			setBg(paintIdx, candiesPosY[i]-1, defaultColor)

			candiesPosY[i]++
			if candiesPosY[i] == 0 {
				candiesPosY[i] = lev.posY
			}
		}
		mut.Lock()
		termbox.Flush()
		mut.Unlock()
		time.Sleep(time.Millisecond * time.Duration(int(duration)))
	}
}

func fallAnimation(lev *level, x int, iterTo int, duration float64, wg *sync.WaitGroup, mut *sync.Mutex) {
	rowCount := lev.ymax + 1
	candiesPosY := make([]int, rowCount)
	paintIdx := coordX(lev.posX, x)

	for i, yPos := 0, (rowCount*2-1)*-1; i < rowCount; i, yPos = i+1, yPos+2 {
		candiesPosY[i] = yPos
	}

	fall(lev, candiesPosY, rowCount, x, paintIdx, duration, mut)

	wg.Done()
}

func initBoardAnimation(lev *level) {
	var wg sync.WaitGroup
	var mut sync.Mutex

	for colIdx := range lev.board[0] {
		wg.Add(1)
		go fallAnimation(
			lev, colIdx, len(lev.board)*2,
			float64(math.Pow(float64(colIdx+2), 2.2))-float64(colIdx)+float64(50),
			&wg, &mut,
		)
	}
	wg.Wait()
}

type cellState struct {
	x, y  int
	color termbox.Attribute
}

func (lev *level) swapAnimation(x, y int) {
	sleep := func() {
		time.Sleep(time.Millisecond * 100)
	}

	xidx := 0

	if x == -2 {
		xidx = -1
	} else if x == 2 {
		xidx = 1
	}
	curColor, adjColor := lev.board[lev.cursor.y][lev.cursor.x].color, lev.board[lev.cursor.y+y][lev.cursor.x+xidx].color

	sequence := [][]cellState{
		{
			{coordX(lev.posX, lev.cursor.x), coordY(lev.posY, lev.cursor.y), defaultColor},
			{coordX(lev.posX, lev.cursor.x) + x, coordY(lev.posY, lev.cursor.y) + y, curColor},
		},
		{
			{coordX(lev.posX, lev.cursor.x) + x*2, coordY(lev.posY, lev.cursor.y) + y*2, curColor},
			{coordX(lev.posX, lev.cursor.x) + x, coordY(lev.posY, lev.cursor.y) + y, adjColor},
		},
		{
			{coordX(lev.posX, lev.cursor.x) + x, coordY(lev.posY, lev.cursor.y) + y, defaultColor},
			{coordX(lev.posX, lev.cursor.x), coordY(lev.posY, lev.cursor.y), adjColor},
		},
	}

	for _, stateGroup := range sequence {
		for _, state := range stateGroup {
			setBg(state.x, state.y, state.color)
		}
		termbox.Flush()
		sleep()
	}

	lev.board[lev.cursor.y][lev.cursor.x].color, lev.board[lev.cursor.y+y][lev.cursor.x+xidx].color = adjColor, curColor
}
