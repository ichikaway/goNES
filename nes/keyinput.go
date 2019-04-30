package nes

import (
	"github.com/nsf/termbox-go"
	"goNES/bus"
	"os"
)

func keyEvent(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}

func getKeyinput(keyCh chan termbox.Key) [8]bool {
	buttons := [8]bool{}
	select {
	case key := <-keyCh:
		mu.Lock()
		switch key {
		case termbox.KeyEsc, termbox.KeyCtrlC: //終了
			mu.Unlock()
			os.Exit(0)
		case termbox.KeyF12:
			buttons[bus.ButtonStart] = true
			break
		case termbox.KeyF11:
			buttons[bus.ButtonSelect] = true
			break
		case termbox.KeyArrowUp:
			buttons[bus.ButtonUp] = true
			break
		case termbox.KeyArrowDown:
			buttons[bus.ButtonDown] = true
			break
		case termbox.KeyArrowLeft:
			buttons[bus.ButtonLeft] = true
			break
		case termbox.KeyArrowRight:
			buttons[bus.ButtonRight] = true
			break
		case termbox.KeyEnter:
			buttons[bus.ButtonA] = true
			break
		case termbox.KeySpace:
			buttons[bus.ButtonB] = true
			break
		}
		mu.Unlock()
		break
	default:
		break
	}
	return buttons
}
