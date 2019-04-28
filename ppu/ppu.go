package ppu

import (
	"goNES/bus"
	"goNES/cpu_interrupts"
)


// PPU power up state
// see. https://wiki.nesdev.com/w/index.php/PPU_power_up_state
//
// Memory map
/*
| addr           |  description               |
+----------------+----------------------------+
| 0x0000-0x0FFF  |  Pattern table#0           |
| 0x1000-0x1FFF  |  Pattern table#1           |
| 0x2000-0x23BF  |  Name table                |
| 0x23C0-0x23FF  |  Attribute table           |
| 0x2400-0x27BF  |  Name table                |
| 0x27C0-0x27FF  |  Attribute table           |
| 0x2800-0x2BBF  |  Name table                |
| 0x2BC0-0x2BFF  |  Attribute table           |
| 0x2C00-0x2FBF  |  Name Table                |
| 0x2FC0-0x2FFF  |  Attribute Table           |
| 0x3000-0x3EFF  |  mirror of 0x2000-0x2EFF   |
| 0x3F00-0x3F0F  |  background Palette        |
| 0x3F10-0x3F1F  |  sprite Palette            |
| 0x3F20-0x3FFF  |  mirror of 0x3F00-0x3F1F   |
*/
/*
  Control Register1 0x2000
| bit  | description                                 |
+------+---------------------------------------------+
|  7   | Assert NMI when VBlank 0: disable, 1:enable |
|  6   | PPU master/slave, always 1                  |
|  5   | Sprite size 0: 8x8, 1: 8x16                 |
|  4   | Bg pattern table 0:0x0000, 1:0x1000         |
|  3   | sprite pattern table 0:0x0000, 1:0x1000     |
|  2   | PPU memory increment 0: +=1, 1:+=32         |
|  1-0 | Name table 0x00: 0x2000                     |
|      |            0x01: 0x2400                     |
|      |            0x02: 0x2800                     |
|      |            0x03: 0x2C00                     |
*/
/*
  Control Register2 0x2001
| bit  | description                                 |
+------+---------------------------------------------+
|  7-5 | Background color  0x00: Black               |
|      |                   0x01: Green               |
|      |                   0x02: Blue                |
|      |                   0x04: Red                 |
|  4   | Enable sprite                               |
|  3   | Enable background                           |
|  2   | Sprite mask       render left end           |
|  1   | Background mask   render left end           |
|  0   | Display type      0: color, 1: mono         |
*/


const SPRITES_NUMBER = 0x100

const SPRITES_ARRAY_MAX = 64

const CYCLES_PER_LINE = 341

type Sprite [8][8]byte

type Ppu struct {
	Registers       []byte
	Cycle           int
	Line            int
	IsValidVramAddr bool
	IsLowerVramAddr bool
	SpriteRamAddr   uint16
	VramAddr        uint16
	Vram            bus.Ram
	VramReadBuf     byte
	SpriteRam       bus.Ram
	Bus             bus.PpuBus
	Background      Background
	Sprites         []SpriteWithAttribute
    RenderingData   RenderingData
	Palette           PaletteRam
	Interrupts        *cpu_interrupts.Interrupts
	IsHrizontalScroll bool
	ScrollX           byte
	ScrollY           byte
	IsHrizontalMirror bool
}

func NewPpu(ppubus bus.PpuBus, interrupts *cpu_interrupts.Interrupts, isHrizontalMirror bool) Ppu {
	ppu := Ppu{
		Registers:         make([]byte, 8),
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
		Sprites:           make([]SpriteWithAttribute, SPRITES_ARRAY_MAX),
		Bus:               ppubus,
		Interrupts:        interrupts,
		IsHrizontalMirror: isHrizontalMirror,
		ScrollX:           0,
		ScrollY:           0,
		Palette:           NewPaletteRam(),
		RenderingData:     RenderingData{},
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



func (this *Ppu) clearVblank() {
	this.Registers[0x02] &= 0x7F
}

func (this *Ppu) clearSpriteHit() {
	this.Registers[0x02] &= 0xbf
}

func (this Ppu) Read(addr uint16) byte {
	/*
        | bit  | description                                 |
        +------+---------------------------------------------+
        | 7    | 1: VBlank clear by reading this register    |
        | 6    | 1: sprite hit                               |
        | 5    | 0: less than 8, 1: 9 or more                |
        | 4-0  | invalid                                     |
        |      | bit4 VRAM write flag [0: success, 1: fail]  |
	*/

	if addr == 0x0002 {
		this.IsHrizontalScroll = true
		data := this.Registers[0x02]
		this.clearVblank()
		return data
	}

	// Write OAM data here. Writes will increment OAMADDR after the write
	// reads during vertical or forced blanking return the value from OAM at that address but do not increment.
	if addr == 0x0004 {
		return this.SpriteRam.Read(this.SpriteRamAddr)
	}
	if addr == 0x0007 {
		return this.readVram()
	}

	return 0x0000
}

func (this Ppu) vramOffset() uint16 {
	if (this.Registers[0x00] & 0x04) == 0x04 {
		return 32
	}
	return 1
}


func (this Ppu) calcVramAddr() uint16 {
	if this.VramAddr >= 0x3000 && this.VramAddr < 0x3f00 {
		return this.VramAddr - 0x3000
	}
	return this.VramAddr - 0x2000
}

func (this *Ppu) readVram() byte {
	buf := this.VramReadBuf
	if this.VramAddr >= 0x2000 {
		addr := this.calcVramAddr()
		this.VramAddr += this.vramOffset()
		if addr >= 0x3F00 {
			return this.Vram.Read(addr)
		}
		this.VramReadBuf = this.Vram.Read(addr)
	} else {
		this.VramReadBuf = this.readCharacterRAM(this.VramAddr)
		this.VramAddr += this.vramOffset()
	}
	return buf
}



func (this *Ppu) Write(addr uint16, data byte) {
	if addr == 0x0003 {
		this.writeSpriteRamAddr(data)
	}
	if addr == 0x0004 {
		this.writeSpriteRamData(data)
	}
	if addr == 0x0005 {
		this.writeScrollData(data)
	}
	if addr == 0x0006 {
		this.writeVramAddr(data)
	}
	if addr == 0x0007 {
		this.writeVramData(data)
	}
	this.Registers[addr] = data
}


func (this *Ppu) writeVramData(data byte) {
	if this.VramAddr >= 0x2000 {
		if this.VramAddr >= 0x3f00 && this.VramAddr < 0x4000 {
			this.Palette.Write(this.VramAddr - 0x3f00, data)
		} else {
			this.Vram.Write(this.calcVramAddr(), data)
		}
	} else {
		this.writeCharacterRAM(this.VramAddr, data)
	}
	this.VramAddr += this.vramOffset()
}

func (this *Ppu) writeVramAddr(data byte) {
	if this.IsLowerVramAddr {
		this.VramAddr += uint16(data)
		this.IsLowerVramAddr = false
		this.IsValidVramAddr = true
	} else {
		this.VramAddr = uint16(data) << 8
		this.IsLowerVramAddr = true
		this.IsValidVramAddr = false
	}
}

func (this *Ppu) writeScrollData(data byte) {
	if this.IsHrizontalScroll {
		this.IsHrizontalScroll = false
		this.ScrollX = data
	} else {
		this.ScrollY = data
		this.IsHrizontalScroll = true
	}
}

func (this *Ppu) writeSpriteRamData(data byte) {
	this.SpriteRam.Write(this.SpriteRamAddr, data)
	this.SpriteRamAddr += 1
}

func (this *Ppu) writeSpriteRamAddr(data byte) {
	this.SpriteRamAddr = uint16(data)
}

func (this Ppu) isBackgroundEnable() bool {
	return (this.Registers[0x01] & 0x08) == 0x08
}
func (this Ppu) isSpriteEnable() bool {
	return (this.Registers[0x01] & 0x10) == 0x10
}

func (this Ppu) hasSpriteHit() bool {
	y := this.SpriteRam.Read(0)
	if int(y) == this.Line && this.isBackgroundEnable() && this.isSpriteEnable() {
		return true
	}
	return false
}

func (this *Ppu) setVblank() {
	this.Registers[0x02] |= 0x80
}
func (this Ppu) hasVblankIrqEnabled() bool {
	return this.Registers[0] & 0x80 == 0x80
}


func (this *Ppu) setSpriteHit() {
	this.Registers[0x02] |= 0x40
}

func (this Ppu) nameTableId() byte {
	return this.Registers[0x00] & 0x03
}

func (this Ppu) scrollTileX() int {
	/*
    Name table id and address
    +------------+------------+
    |            |            |
    |  0(0x2000) |  1(0x2400) |
    |            |            |
    +------------+------------+
    |            |            |
    |  2(0x2800) |  3(0x2C00) |
    |            |            |
    +------------+------------+
	*/
	return (int(this.ScrollX) + ((int(this.nameTableId()) % 2) * 256)) / 8
}

func (this Ppu) scrollTileY() int {
	return (int(this.ScrollY) + (int(this.nameTableId()) / 2 * 240)) / 8
}

func (this Ppu) tileY() int {
	return (this.Line / 8) + this.scrollTileY()
}

func getBlockId(tileX int, tileY int) byte {
	return uint8(((tileX % 4) / 2) + ((tileY % 4) / 2) * 2)
}

func (this Ppu) getAttribute(tileX int, tileY int, offset int) byte {
	addr := (tileX / 4) + ((tileY / 4) * 8) + 0x03c0 + offset
	return this.Vram.Read(uint16(addr))
}

func (this Ppu) getSpriteId(tileX int, tileY int, offset int) byte {
	tileNumber := tileY * 32 + tileX
	spriteAddr := this.mirrorDownSpriteAddr(uint16(tileNumber + offset))
	return this.Vram.Read(spriteAddr)
}

func (this Ppu) mirrorDownSpriteAddr(addr uint16) uint16 {
	if !this.IsHrizontalMirror {
		return addr
	}
	if addr >= 0x0400 && addr < 0x0800 || addr >= 0x0c00 {
		return addr - 0x400
	}
	return addr
}

func (this Ppu) backgroundTableOffset() uint16 {
	if (this.Registers[0] & 0x10) == 0x10 {
		return 0x1000
	}
	return 0x0000
}


func (this *Ppu) Run(cpuCycle int) bool {
	cycle := this.Cycle + cpuCycle
	this.Cycle = cycle
	if cycle < CYCLES_PER_LINE {
		return false
	}
//fmt.Println("cycle: ", cycle, " line: ", this.Line)
	if this.Line == 0 {
		this.Background.Clear()
		this.buildSprites()
	}

	if this.Cycle >= CYCLES_PER_LINE {
		this.Cycle = cycle - CYCLES_PER_LINE
		this.Line++

		if this.hasSpriteHit() {
			this.setSpriteHit()
		}
		if this.Line <= 240 && this.Line % 8 == 0 && this.ScrollY <= 240 {
			this.buildBackground()
		}
		if this.Line == 241 {
			//log.Println("before setVblank", this.Registers[0x02])
			this.setVblank()
			//log.Println("after setVblank", this.Registers[0x02])
			if this.hasVblankIrqEnabled() {
				this.Interrupts.AssertNmi()
				//log.Println("in hasVblankIrqEnabled: ", this.Interrupts.Nmi)
			}
		}

		if this.Line == 262 {
			this.clearVblank()
			this.clearSpriteHit()
			this.Line = 0
			this.Interrupts.DeassertNmi()

			background := Background{}
			if this.isBackgroundEnable() {
				background = this.Background
			}
			sprites := []SpriteWithAttribute{}
			if this.isSpriteEnable() {
				sprites = this.Sprites
			}
			this.RenderingData = RenderingData{
				Palette:    this.Palette.Read(),
				Background: background,
				Sprites:    sprites,
			}
			return true
		}
	}

	return false
}

func (this *Ppu) buildBackground() {
	// INFO: Horizontal offsets range from 0 to 255. "Normal" vertical offsets range from 0 to 239,
	// while values of 240 to 255 are treated as -16 through -1 in a way, but tile data is incorrectly
	// fetched from the attribute table.

	clampedTileY := this.tileY() % 30
	tableIdOffset := 0
	if ((this.tileY() / 30) % 2) > 0 {
		tableIdOffset = 2
	}
	// background of a line.
	// Build viewport + 1 tile for background scroll.
	for x := 0 ; x < 32 + 1 ; x = (x+1) {
		tileX := x + this.scrollTileX()
		clampedTileX := tileX % 32
		nameTableId := ((tileX / 32) % 2) + tableIdOffset
		offsetAddrByNameTable := nameTableId * 0x400
		tile := this.buildTile(clampedTileX, clampedTileY, offsetAddrByNameTable)
		this.Background.Tiles = append(this.Background.Tiles, tile)
	}
}

func (this Ppu) buildTile(tileX int, tileY int, offset int) Tile {
	// INFO see. http://hp.vector.co.jp/authors/VA042397/nes/ppu.html
	blockId := getBlockId(tileX, tileY)
	spriteId := this.getSpriteId(tileX, tileY, offset)
	attr := this.getAttribute(tileX, tileY, offset)
	paletteId := (attr >> (blockId * 2)) & 0x03
	sprite := this.buildSprite(spriteId, this.backgroundTableOffset())

	return Tile{
		Sprite:    sprite,
		Scroll_x:  this.ScrollX,
		Scroll_y:  this.ScrollY,
		PaletteId: int(paletteId),
	}
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
		if y < 8 {
			return
		}

		spriteId := this.SpriteRam.Read(uint16(i+1))
		attr := this.SpriteRam.Read(uint16(i+2))
		x := this.SpriteRam.Read(uint16(i+3))
		sprite = this.buildSprite(spriteId, offset)
		this.Sprites[i/4] = NewStripeWithAttribute(sprite, x, y, attr, spriteId)
		//fmt.Println(sprite)
	}
}

func (this *Ppu) buildSprite(spriteId uint8, offset uint16) Sprite {
	sprite := Sprite{}
	for i := 0 ; i < 16 ; i++ {
		for j := 0 ; j < 8 ; j++ {
			addr := uint16(spriteId) * 16 + uint16(i) + offset
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

func (this *Ppu) writeCharacterRAM(addr uint16, data byte) {
	this.Bus.WriteByPpu(addr, data)
}
