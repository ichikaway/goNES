package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func keyinput(key chan<- byte){
	for{
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		input, _ := consoleReader.ReadByte()
		//fmt.Println(input)
		key <- input
	}

}

func main() {
	fmt.Println("keypad test")
	key := make(chan byte)

	go keyinput(key)

	for i:=0;i<1000;i++ {
		//fmt.Print(" ")
		select{
		case button := <- key:
			fmt.Print(button)
		default:
			time.Sleep(1000 * time.Millisecond)
		}

	}

}