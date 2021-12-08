package game

import "github.com/nsf/termbox-go"

func (lev level) isCursorNotInBounds() bool {
	return lev.cursor.y < 0 || lev.cursor.y > lev.ymax ||
		lev.cursor.x < 0 || lev.cursor.x > lev.xmax
}

func (lev level) coordX(x int) int {
	return 4*x + lev.posX
}

func (lev level) coordY(y int) int {
	return 2*y + lev.posY
}

type adjacentColorMap map[string]adjacentCell

type adjacentCell struct {
	color termbox.Attribute
	index coord
}

var prevCellColor termbox.Attribute
var adjacentColors adjacentColorMap

func (adj adjacentColorMap) repaintCells(lev *level) {
	for adj := range adjacentColors {
		lev.board[adjacentColors[adj].index.y][adjacentColors[adj].index.x].color = adjacentColors[adj].color
	}
	lev.render()
}

func setBg(x, y int, candyColor termbox.Attribute) {
	termbox.SetBg(x, y, candyColor)
	termbox.SetBg(x+1, y, candyColor)
}
