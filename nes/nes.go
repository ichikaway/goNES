package nes

type Rom struct {
	ProgramRom   []byte
	CharacterRom []byte
	isHorizontalMirror bool
	mapper uint8
}

func New(data []byte) Rom {
	return parse(data)
}

