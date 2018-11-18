package dma

import (
	"goNES/bus"
	"goNES/ppu"
)

type Dma struct {
	Ram          bus.Ram
	Ppu          ppu.Ppu
	isProcessing bool
	Register     byte
}

func NewDma(ram bus.Ram, ppu ppu.Ppu) Dma {
	dma := Dma{
		Ram:          ram,
		Ppu:          ppu,
		isProcessing: false,
		Register:     0x0000,
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

	addr := uint16(this.Register) << 8
	for i := uint16(0x0000); i < 0x100; i++  {
		this.Ppu.TransferSprite(i, this.Ram.Read(addr+i))
	}
	this.isProcessing = false;
}

func (dma *Dma) Write(data byte) {
	dma.Register = data
	dma.isProcessing = true
}
