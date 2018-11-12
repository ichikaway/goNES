package dma

import (
	"goNES/bus"
	"goNES/ppu"
)

type Dma struct {
	Ram          bus.Ram
	Ppu          ppu.Ppu
	isProcessing bool
	RamAddr      int
}

func NewDma(ram bus.Ram, ppu ppu.Ppu) Dma {
	dma := Dma{
		Ram:          ram,
		Ppu:          ppu,
		isProcessing: false,
		RamAddr:      0x0000,
	}
	return dma
}

func (dma Dma) IsDmaProcessing() bool {
	return dma.isProcessing
}

func (this *Dma) RunDma() {
	if !this.isProcessing {
		return
	}
	for i := 0; i < 0x100; i = (i + 1) | 0 {
		this.Ppu.TransferSprite(i, this.Ram.Read(this.RamAddr+i))
	}
	this.isProcessing = false;
}

func (dma *Dma) Write(data int) {
	dma.RamAddr = data << 8
	dma.isProcessing = true
}
