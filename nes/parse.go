package nes

import "fmt"

const NES_HEADER_SIZE int = 0x0010;

const PROGRAM_ROM_SIZE int = 0x4000;   //16KB unit プログラムROMの単位
const CHARACTER_ROM_SIZE int = 0x2000; //8KB unit キャラクターROMの単位

func parse(nes []byte) Rom {

	fmt.Println("---------- Parse --------------")
	programRomPages := nes[4];
	characterRomPages := nes[5];
	isHorizontalMirror := nes[6] &^ 0x01 // !(nes[6] & 0x01)
	mapper := ((nes[6] & 0xF0) >> 4) | nes[7]&0xF0

	characterRomStart := NES_HEADER_SIZE + (int(programRomPages) * PROGRAM_ROM_SIZE)
	characterRomEnd := characterRomStart + (int(characterRomPages) * CHARACTER_ROM_SIZE)

	fmt.Printf("Program ROM pages: %d\n", programRomPages);
	fmt.Printf("Character ROM pages: %d\n", characterRomPages);
	fmt.Printf("Character ROM start:  0x%x (%d)\n", characterRomStart, characterRomStart);
	fmt.Printf("Character ROM end:  0x%x (%d)\n", characterRomEnd, characterRomEnd);
	fmt.Printf("isHrizontalMirror: %d\n", isHorizontalMirror);
	fmt.Printf("Mapper: %d\n", mapper);
	nesData := Rom{}

	//スライスは、[n:m]の場合、mはm-1番目までのため、プログラムROMの終端はcharacterRomStartで良い
	nesData.ProgramRom = nes[NES_HEADER_SIZE:(characterRomStart)]
	nesData.CharacterRom = nes[characterRomStart:characterRomEnd]
	nesData.mapper = mapper
	if isHorizontalMirror == 1 {
		nesData.isHorizontalMirror = true
	}

	fmt.Printf("ProgramRom size: %x, %d\n", len(nesData.ProgramRom), len(nesData.ProgramRom));
	fmt.Printf("CharacterRom size: %x, %d\n", len(nesData.CharacterRom), len(nesData.CharacterRom));
	return nesData
}
