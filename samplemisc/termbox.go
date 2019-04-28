package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}

func drawDot(i int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawLine(0, 0, "Press ESC to exit." + string(i))
	drawLine(2, 1, fmt.Sprintf("date: %v", time.Now()))

	for y := 2 ; y < 30 ; y++ {
		drawLine(0, y, "--------------------------------------------------------------------------------")
		drawLine(0, y, "--------------------------------------------------------------------------------")

	}

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for i := 0 ; i < 3000 ; i++ {
		drawDot(i)
		time.Sleep(1000)
	}

	os.Exit(0)
}
