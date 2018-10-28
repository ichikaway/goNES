package nes

import (
	"goNes/bus"
)

type Rom struct {
	ProgramRom   []byte
	CharacterRom []byte
	isHorizontalMirror bool
	mapper uint8
}

type Nes struct {
	Rom Rom
	Ram bus.Ram
	characterMem bus.Ram
	ProgramRom bus.Rom
}

func New(data []byte) Nes {
	return Nes{Rom: parse(data)}
}

func (nes *Nes) Load() {
	nes.Ram = bus.NewRam(2048)
	nes.characterMem = bus.NewRam(0x4000)
	for i := 0; i < len(nes.Rom.CharacterRom); i++ {
		nes.characterMem.Write(i, nes.Rom.CharacterRom[i])
	}

}

