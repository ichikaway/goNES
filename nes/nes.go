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
}

func New(data []byte) Nes {
	return Nes{Rom: parse(data)}
}

func (nes *Nes) Load() {
	nes.Ram = bus.NewRam(2048)
}

