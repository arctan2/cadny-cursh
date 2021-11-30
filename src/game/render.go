package game

import "github.com/nsf/termbox-go"

var (
	defaultColor = termbox.ColorDefault
)

func render() {
	termbox.Clear(defaultColor, defaultColor)

	renderBoard()

	termbox.Flush()
}

func renderBoard() {
	x, y := 1, 1

	for _, row := range board {
		for _, candy := range row {
			termbox.SetBg(x, y, candy.color)
			termbox.SetBg(x+1, y, candy.color)
			x += 4
		}
		y += 2
		x = 1
	}
}
