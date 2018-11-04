package nes

import (
	"goNES/bus"
	"goNES/cpu"
	"goNES/cpubus"
	"goNES/dma"
	"goNES/ppu"
)

type Rom struct {
	ProgramRom         []byte
	CharacterRom       []byte
	isHorizontalMirror bool
	mapper             uint8
}

type Nes struct {
	Rom          Rom
	Ram          bus.Ram
	characterMem bus.Ram
	ProgramRom   bus.Rom
	PpuBus       bus.PpuBus
	Interrupts   cpu.Interrupts
	Dma          dma.Dma
	Ppu          ppu.Ppu
	CpuBus       cpubus.CpuBus
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

	nes.ProgramRom = bus.NewRom(nes.Rom.ProgramRom)
	nes.PpuBus = bus.NewPpuBus(&nes.characterMem)
	nes.Interrupts = cpu.NewInterrupts()

	nes.Ppu = ppu.NewPpu(nes.PpuBus, nes.Interrupts, nes.Rom.isHorizontalMirror)

	nes.Dma = dma.NewDma(nes.Ram, nes.Ppu)

	nes.CpuBus = cpubus.NewCpuBus(
		nes.Ram,
		nes.ProgramRom,
		nes.Ppu,
		nes.Dma,
		)
}
