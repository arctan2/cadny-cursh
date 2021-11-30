package game

import (
	"time"

	"github.com/nsf/termbox-go"
)

var (
	keyboardChan = make(chan keyboardEvent)
)

func Start() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	initBoard()
	render()

	go listenToKeyboard(keyboardChan)

mainloop:
	for {
		select {
		case e := <-keyboardChan:
			switch e.eventType {
			case END:
				break mainloop
			}

		default:
			render()
			time.Sleep(time.Millisecond * 10)
		}
	}
}
