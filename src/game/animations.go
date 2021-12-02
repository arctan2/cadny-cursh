package game

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

var prevCellColor termbox.Attribute
var adjacentColors adjacentColorMap

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

func (lev *level) blinkCursor() {
	for {
		if lev.isSelected {
			lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
			break
		}
		if lev.board[lev.cursor.y][lev.cursor.x].color == defaultColor {
			lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
		} else {
			prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
			lev.board[lev.cursor.y][lev.cursor.x].color = defaultColor
		}
		time.Sleep(time.Millisecond * 300)
	}
}

type adjacentColorMap map[string]adjacentCell

type adjacentCell struct {
	color termbox.Attribute
	index coord
}

func getAdjacentMap(board [][]candy, curs *cursor, xmax, ymax int) adjacentColorMap {
	adjacentColors := make(adjacentColorMap)
	if curs.x-1 >= 0 {
		adjacentColors["left"] = adjacentCell{board[curs.y][curs.x-1].color, coord{curs.x - 1, curs.y}}
	}
	if curs.x+1 <= xmax {
		adjacentColors["right"] = adjacentCell{board[curs.y][curs.x+1].color, coord{curs.x + 1, curs.y}}
	}
	if curs.y-1 >= 0 {
		adjacentColors["up"] = adjacentCell{board[curs.y-1][curs.x].color, coord{curs.x, curs.y - 1}}
	}
	if curs.y+1 <= ymax {
		adjacentColors["down"] = adjacentCell{board[curs.y+1][curs.x].color, coord{curs.x, curs.y + 1}}
	}

	return adjacentColors
}

func (lev *level) blinkAdjacent() {
	adjacentColors = getAdjacentMap(lev.board, &lev.cursor, len(lev.board[0])-1, len(lev.board)-1)
	for {
		if !lev.isSelected {
			break
		}
		for adj := range adjacentColors {
			if lev.board[adjacentColors[adj].index.y][adjacentColors[adj].index.x].color == defaultColor {
				lev.board[adjacentColors[adj].index.y][adjacentColors[adj].index.x].color = adjacentColors[adj].color
			} else {
				lev.board[adjacentColors[adj].index.y][adjacentColors[adj].index.x].color = defaultColor
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}

/*
0 -> 3
1 -> 7
2 -> 11

4x + 3

0 -> 2
1 -> 4
2 -> 6

2x + 2

*/

func coordX(posX, x int) int {
	return 4*x + posX
}

func coordY(posY, y int) int {
	return 2*y + posY
}

type cellState struct {
	x, y  int
	color termbox.Attribute
}

func (lev *level) swapAnimation(x, y int) {
	sleep := func() {
		time.Sleep(time.Millisecond * 100)
	}

	var xidx int = 0

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
