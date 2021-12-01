package game

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type candy struct {
	color termbox.Attribute
}

func (lev *level) initBoard() {
	rand.Seed(time.Now().UnixNano())
	colors := map[int]termbox.Attribute{
		0: termbox.ColorBlue,
		1: termbox.ColorGreen,
		2: termbox.ColorYellow,
		3: termbox.ColorRed,
	}

	for i := 0; i < len(lev.board); i++ {
		for j := 0; j < len(lev.board[i]); j++ {
			lev.board[i][j] = candy{colors[rand.Intn(4)]}
		}
	}

	initBoardAnimation(lev)
}
