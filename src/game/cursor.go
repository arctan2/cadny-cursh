package game

import "time"

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
	adjacentColors = getAdjacentMap(lev.board, &lev.cursor, lev.xmax, lev.ymax)
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

func (lev *level) blinkCursor(stopBlink <-chan bool) {
blinkloop:
	for {
		select {
		case <-stopBlink:
			break blinkloop
		case <-time.After(time.Millisecond * 300):
			if lev.isCursorNotInBounds() {
				continue blinkloop
			}
			if lev.board[lev.cursor.y][lev.cursor.x].color == defaultColor {
				lev.board[lev.cursor.y][lev.cursor.x].color = prevCellColor
			} else {
				prevCellColor = lev.board[lev.cursor.y][lev.cursor.x].color
				lev.board[lev.cursor.y][lev.cursor.x].color = defaultColor
			}
		}
	}
}
