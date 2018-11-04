package cpubus

import (
	"goNES/bus"
	"goNES/dma"
	"goNES/ppu"
)

type CpuBus struct {
	Ram bus.Ram
	ProgramRom bus.Rom
	Ppu ppu.Ppu
	//KeyPad Keypad
	Dma dma.Dma
}

func NewCpuBus(ram bus.Ram, programRom bus.Rom, ppu ppu.Ppu, dma dma.Dma) CpuBus {

	cpuBus := CpuBus{
		Ram: ram,
		ProgramRom: programRom,
		Ppu: ppu,
		Dma: dma,
	}
	return cpuBus
}

