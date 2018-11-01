package ppu

import (
	"goNES/bus"
	"goNES/cpu"
)

const SPRITES_NUMBER = 0x100

type Ppu struct {
	Registers       []int
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
	IsHrizontalMirror bool
}

func NewPpu(ppubus bus.PpuBus, interrupts cpu.Interrupts, isHrizontalMirror bool) Ppu {
	ppu := Ppu{
		Registers:         make([]int, 7),
		Cycle:             0,
		Line:              0,
		IsValidVramAddr:   false,
		IsLowerVramAddr:   false,
		IsHrizontalScroll: true,
		VramAddr:          0x0000,
		Vram:              bus.NewRam(0x2000),
		VramReadBuf:       0,
		SpriteRam:         bus.NewRam(0x100),
		SpriteRamAddr:     0,
		//Background: []
		//Sprites: []
		Bus:               ppubus,
		Interrupts:        interrupts,
		IsHrizontalMirror: isHrizontalMirror,
		ScrollX:           0,
		ScrollY:           0,
		Palette:           NewPaletteRam(),
	}
	return ppu
}
