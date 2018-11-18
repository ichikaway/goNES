package ppu

import "goNES/bus"

type PaletteRam struct {
	PaletteRam bus.Ram
}

func NewPaletteRam() PaletteRam {
	ram := bus.NewRam(0x20)
	return PaletteRam{PaletteRam: ram}
}

func isSpriteMirror(addr uint16) bool {
	return (addr == 0x10) || (addr == 0x14) || (addr == 0x18) || (addr == 0x1c)
}

func isBackgroundMirror(addr uint16) bool {
	return (addr == 0x04) || (addr == 0x08) || (addr == 0x0c)
}

func (this PaletteRam) Read() []byte {
	length := this.PaletteRam.Size()

	ret := make([]byte, length)

	for i := 0; i < length; i++ {
		if isSpriteMirror(uint16(i)) {
			ret[i] = this.PaletteRam.Read(uint16(i - 0x10))
		} else if isBackgroundMirror(uint16(i)) {
			ret[i] = this.PaletteRam.Read(0x00)
		} else {
			ret[i] = this.PaletteRam.Read(uint16(i))
		}
	}
	return ret
}

func getPaletteAddr(addr uint16) uint16{
	mirrorDowned := ((addr & 0xFF) % 0x20)
	if isSpriteMirror(mirrorDowned) {
		return 	mirrorDowned - 0x10
	}
	return mirrorDowned
}

func (this *PaletteRam) Write(addr uint16, data byte) {
	this.PaletteRam.Write(getPaletteAddr(addr) ,data)
}

