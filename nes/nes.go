package nes

import "fmt"

const NES_HEADER_SIZE int = 0x0010;

const PROGRAM_ROM_SIZE int = 0x4000;   //16KB unit プログラムROMの単位
const CHARACTER_ROM_SIZE int = 0x2000; //8KB unit キャラクターROMの単位

type Nes struct {
	ProgramRom   []byte
	CharacterRom []byte
	isHorizontalMirror bool
}

func New(data []byte) Nes {
	nes := Nes{}
	nes.parse(data)
	return nes
}

func (nesData *Nes) parse(nes []byte) {
	fmt.Println("---------- Parse --------------")
	programRomPages := nes[4];
	characterRomPages := nes[5];
	isHorizontalMirror := nes[6] &^ 0x01 // !(nes[6] & 0x01)
	mapper := (((nes[6] & 0xF0) >> 4) | nes[7]&0xF0)

	characterRomStart := NES_HEADER_SIZE + (int(programRomPages) * PROGRAM_ROM_SIZE)
	characterRomEnd := characterRomStart + (int(characterRomPages) * CHARACTER_ROM_SIZE)

	fmt.Printf("Program ROM pages: %d\n", programRomPages);
	fmt.Printf("Character ROM pages: %d\n", characterRomPages);
	fmt.Printf("Character ROM start:  0x%x (%d)\n", characterRomStart, characterRomStart);
	fmt.Printf("Character ROM end:  0x%x (%d)\n", characterRomEnd, characterRomEnd);
	fmt.Printf("isHrizontalMirror: %d\n", isHorizontalMirror);
	fmt.Printf("Mapper: %d\n", mapper);

	if isHorizontalMirror == 1 {
		nesData.isHorizontalMirror = true
	}

	nesData.ProgramRom = nes[NES_HEADER_SIZE : (characterRomStart - 1)]
	nesData.CharacterRom = nes[characterRomStart : characterRomEnd]

	fmt.Printf("ProgramRom size: %d\n", len(nesData.ProgramRom));
	fmt.Printf("CharacterRom size: %d\n", len(nesData.CharacterRom));
}
