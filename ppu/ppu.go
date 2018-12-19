package ppu

import (
	"goNES/bus"
	"goNES/cpu_interrupts"
)

const SPRITES_NUMBER = 0x100

const CYCLES_PER_LINE = 341

type Sprite [8][8]byte

type Ppu struct {
	Registers       []int
	Cycle           int
	Line            int
	IsValidVramAddr bool
	IsLowerVramAddr bool
	SpriteRamAddr   uint16
	VramAddr        uint16
	Vram            bus.Ram
	VramReadBuf     int
	SpriteRam       bus.Ram
	Bus             bus.PpuBus
	Background      Background

	/** @var \Nes\Ppu\SpriteWithAttribute[] */
	//Sprites

	Palette           PaletteRam
	Interrupts        cpu_interrupts.Interrupts
	IsHrizontalScroll bool
	ScrollX           int
	ScrollY           int
	IsHrizontalMirror bool
}

func NewPpu(ppubus bus.PpuBus, interrupts cpu_interrupts.Interrupts, isHrizontalMirror bool) Ppu {
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
		Background:        NewBackground(),
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

func (this *Ppu)TransferSprite(index uint16, data byte) {
	// The DMA transfer will begin at the current OAM write address.
	// It is common practice to initialize it to 0 with a write to PPU 0x2003 before the DMA transfer.
	// Different starting addresses can be used for a simple OAM cycling technique
	// to alleviate sprite priority conflicts by flickering. If using this technique
	// after the DMA OAMADDR should be set to 0 before the end of vblank to prevent potential OAM corruption
	// (See: Errata).
	// However, due to OAMADDR writes also having a "corruption" effect[5] this technique is not recommended.
	addr := index + this.SpriteRamAddr
	this.SpriteRam.Write(addr % 0x100, data) //256以上のアドレスに入れさせないために256の剰余を求める
}

func (this Ppu) Read(addr uint16) byte {
	//todo
	return 0x0000
}

func (this *Ppu) Write(addr uint16, data byte) {
	//todo
}

func (this *Ppu) Run(cpuCycle int) bool {
	cycle := this.Cycle + cpuCycle
	if cycle < CYCLES_PER_LINE {
		this.Cycle = cycle
		return false
	}

	if this.Line == 0 {
		this.Background.Clear()
		this.buildSprites()
	}

	this.Cycle = cycle - CYCLES_PER_LINE
	this.Line++


	return false
}

func (this *Ppu) buildSprites() {
	var offset uint16 = 0x0000
	var sprite Sprite

	if (this.Registers[0] & 0x08) > 0 {
		offset = 0x1000
	}

	for i := 0 ; i < SPRITES_NUMBER ; i = (i+4)  {
		// INFO: Offset sprite Y position, because First and last 8line is not rendered.
		y := this.SpriteRam.Read(uint16(i))
		if !(y >= 8 && y < 224) {
			return
		}
		spriteId := this.SpriteRam.Read(uint16(i+1))
		//attr := this.SpriteRam.Read(uint16(i+2))
		//x := this.SpriteRam.Read(uint16(i+3))
		sprite = this.buildSprite(spriteId, offset)
	}
	/*
	public function buildSprites()
	{
	$offset = ($this->registers[0] & 0x08) ? 0x1000 : 0x0000;
	for ($i = 0; $i < self::SPRITES_NUMBER; $i = ($i + 4) | 0) {
	// INFO: Offset sprite Y position, because First and last 8line is not rendered.
	$y = $this->spriteRam->read($i) - 8;
	if ($y < 0) {
	return;
	}
	$spriteId = $this->spriteRam->read($i + 1);
	$attr = $this->spriteRam->read($i + 2);
	$x = $this->spriteRam->read($i + 3);
	$sprite = $this->buildSprite($spriteId, $offset);
	$this->sprites[$i / 4] = new SpriteWithAttribute($sprite, $x, $y, $attr, $spriteId);
	}
	*/
}

func (this *Ppu) buildSprite(spriteId uint8, offset uint16) Sprite {
	sprite := Sprite{}
	for i := 0 ; i < 16 ; i++ {
		for j := 0 ; j < 8 ; j++ {
			addr := uint16(spriteId * 16) + uint16(i) + offset
			ram := this.readCharacterRAM(addr)
			if (ram & uint8(0x80 >> uint8(j))) != 0 {
				sprite[i%8][j] += uint8(0x01 << uint8(i/8))
			}
		}
	}
	return sprite
}

func (this *Ppu) readCharacterRAM(addr uint16) byte {
	return this.Bus.ReadByPpu(addr)
}
