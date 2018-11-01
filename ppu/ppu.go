package ppu

import (
	"goNES/bus"
	"goNES/cpu"
)

const SPRITES_NUMBER = 0x100

type Ppu struct {
	Cycle           int
	Line            int
	IsValidVramAddr bool
	IsLowerVramAddr bool
	SpriteRamAddr   int
	VramAddr        int
	Vram            bus.Ram
	VramReadBuf     int
	SpriteRam       bus.Ram
	Bus             bus.PpuBus

	/** @var \Nes\Ppu\Tile[] */
	//Background

	/** @var \Nes\Ppu\SpriteWithAttribute[] */
	//Sprites

	Palette           PaletteRam
	Interrupts        cpu.Interrupts
	IsHrizontalScroll bool
	ScrollX           int
	ScrollY           int
	isHrizontalMirror bool
}

func NewPpu(bus bus.PpuBus, interrupts cpu.Interrupts, isHrizontalMirror bool) Ppu {
	ppu := Ppu{}
	return ppu
}
