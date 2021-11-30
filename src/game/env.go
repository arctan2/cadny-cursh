package game

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	board [8][8]candy
)

type candy struct {
	color termbox.Attribute
}

func initBoard() {
	rand.Seed(time.Now().UnixNano())
	colors := map[int]termbox.Attribute{
		0: termbox.ColorBlue,
		1: termbox.ColorGreen,
		2: termbox.ColorYellow,
		3: termbox.ColorRed,
	}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			board[i][j] = candy{colors[rand.Intn(4)]}
		}
	}
}
