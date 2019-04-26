package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

func keyinput(key chan<- byte){
	for{
		in := bufio.NewReader(os.Stdin)
		b, _ := in.ReadByte()
		key <- b

		/*
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		input, _ := consoleReader.ReadByte()
		//fmt.Println(input)
		key <- input
		*/
	}

}

func keyEvent(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}

func main() {
	fmt.Println("keypad test")

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	keyCh := make(chan termbox.Key)
	go keyEvent(keyCh)

	for {
		time.Sleep(50 * time.Millisecond)
		select {
			case key := <-keyCh:
				mu.Lock()
				switch key {
				case termbox.KeyEsc, termbox.KeyCtrlC: //終了
					mu.Unlock()
					return
				}
				mu.Unlock()
				fmt.Println(key)
				break
			default:
				//fmt.Println("---")
				break
		}
	}

/*	key := make(chan byte)
	go keyinput(key)
	for i:=0;i<1000;i++ {
		//fmt.Print(" ")
		select{
		case button := <- key:
			fmt.Print(button)
		default:
			time.Sleep(1000 * time.Millisecond)
		}

	}*/

}