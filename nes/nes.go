package nes

const NES_HEADER_SIZE int = 0x0010;

const PROGRAM_ROM_SIZE int = 0x4000;   //16KB unit プログラムROMの単位
const CHARACTER_ROM_SIZE int = 0x2000; //8KB unit キャラクターROMの単位

type Nes struct {
	ProgramRom   []byte
	CharacterRom []byte
	isHorizontalMirror bool
	mapper uint8
}

func New(data []byte) Nes {
	return parse(data)
}

