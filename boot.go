package main

import (
	"fmt"
	"goNES/nes"
	"io/ioutil"
	"os"
)

const ERROR_EXIT int = 1

func main() {
	fmt.Println("goNes: NES emulator written in golang.")

	if len(os.Args) < 2 {
		fmt.Println("Error: No NES file.")
		os.Exit(ERROR_EXIT)
	}
	filename := os.Args[1]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error: can not open NES file.", filename)
		os.Exit(ERROR_EXIT)
	}

	nes := nes.New(data)
	nes.Load()
	//nes.Ram.Write(1, 1)
	//nes.Ram.Reset()
	fmt.Println(nes.Ram)
	//println(data[1])

}
