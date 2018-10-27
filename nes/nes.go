package nes

type Nes struct {
	ProgramRom   []byte
	CharacterRom []byte
	isHorizontalMirror bool
	mapper uint8
}

func New(data []byte) Nes {
	return parse(data)
}

