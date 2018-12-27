package main

/**
 NES ROMファイルからキャラクターROMデータを切り出してPNGファイル出力するプログラム
 */

import (
	"fmt"
	"goNES/nes"
	"io/ioutil"
	"os"
)

const ERROR_EXIT int = 1

const DEFAULT_CANVAS_WIDTH = 800
const PIXEL_RATIO = 1

type Sprite [8][8]byte

func main() {
	fmt.Println("goNes toPng")

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

	nesData := nes.New(data)
	fmt.Println(nesData)

	spritesPerRow := DEFAULT_CANVAS_WIDTH / (8 * PIXEL_RATIO)
	spritesNum := len(nesData.Rom.CharacterRom) / 16
	rowNum := (spritesNum / spritesPerRow) + 1

	height := rowNum * 8 * PIXEL_RATIO

}


func buildSprite(charRom []byte, spriteId uint8) Sprite {
	sprite := Sprite{}
	for i := 0 ; i < 16 ; i++ {
		for j := 0 ; j < 8 ; j++ {
			addr := uint16(spriteId * 16) + uint16(i)
			//fmt.Println(addr)
			ram := charRom[addr]
			if (ram & uint8(0x80 >> uint8(j))) != 0 {
				sprite[i%8][j] += uint8(0x01 << uint8(i/8))
			}
		}
	}
	return sprite
}
