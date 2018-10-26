package nes

import "fmt"

const NES_HEADER_SIZE int = 0x0010;
const PROGRAM_ROM_SIZE int = 0x4000; //16KB unit プログラムROMの単位
const CHARACTER_ROM_SIZE int = 0x2000; //8KB unit キャラクターROMの単位

type Nes struct {
	rom []byte
	ProgramRom []byte
	CharacterRom []byte
}

func New(data []byte) Nes{
	nes := Nes{rom:data}
	nes.parse()
	return nes
}

func (nesData *Nes) parse () {
	nes := nesData.rom

	programRomPages := nes[4];
	characterRomPages := nes[5];
	isHorizontalMirror := nes[6] &^ 0x01 // !(nes[6] & 0x01)
	mapper := (((nes[6] & 0xF0) >> 4) | nes[7] & 0xF0)

	characterRomStart := NES_HEADER_SIZE + (int(programRomPages) * PROGRAM_ROM_SIZE)
	characterRomEnd   := characterRomStart + (int(characterRomPages) * CHARACTER_ROM_SIZE)

	fmt.Println("isHorizontalMirror: ", isHorizontalMirror)
	fmt.Println("mapper: ", mapper)
	fmt.Println("start: ", characterRomStart)
	fmt.Println("end", characterRomEnd)
}
