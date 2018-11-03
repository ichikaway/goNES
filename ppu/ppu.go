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

func (this *Ppu)TransferSprite(index int, data byte) {
	// The DMA transfer will begin at the current OAM write address.
	// It is common practice to initialize it to 0 with a write to PPU 0x2003 before the DMA transfer.
	// Different starting addresses can be used for a simple OAM cycling technique
	// to alleviate sprite priority conflicts by flickering. If using this technique
	// after the DMA OAMADDR should be set to 0 before the end of vblank to prevent potential OAM corruption
	// (See: Errata).
	// However, due to OAMADDR writes also having a "corruption" effect[5] this technique is not recommended.
	addr := index + this.SpriteRamAddr
	this.SpriteRam.Write(addr % 0x100, data)
}