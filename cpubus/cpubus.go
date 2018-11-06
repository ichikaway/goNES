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

func (this CpuBus) ReadByCpu(addr int) byte {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return this.Ram.Read(addr)
	case 0x2000 <= addr && addr <= 0x3FFF:
		return this.Ppu.Read(addr - 0x2000)
	case addr == 0x4016:
		// todo  0x4016 => self.keypad.read(),
		return 0x0000
	case addr == 0x4017:
		// todo 2payler
		return 0x0000
	case 0x4000 <= addr && addr <= 0x401F:
		//todo  0x4000...0x401F => self.apu.read(addr - 0x4000),
		return 0x0000
	case 0x8000 <= addr && addr <= 0xBFFF:
		return this.ProgramRom.Read(addr - 0x8000)

	case 0xC000 <= addr && addr <= 0xFFFF:
		if this.ProgramRom.Size() <= 0x4000 {
			return this.ProgramRom.Read(addr - 0xC000)
		}
		return this.ProgramRom.Read(addr - 0x8000)
	}
	if addr < 0x0800 {
		return this.Ram.Read(addr)
	}
	return 0x0000
}


func (this *CpuBus) WriteByCpu(addr int, data byte) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		this.Ram.Write(addr&0x07FF, data)
	case 0x2000 <= addr && addr <= 0x3FFF:
		this.Ppu.Write(addr-0x2000, data)
	case addr == 0x4014:
		this.Dma.Write(int(data))
	case addr == 0x4016:
		// todo  0x4016 => self.keypad.write(data),
	case addr == 0x4017:
		// todo 2payler
	case 0x4000 <= addr && addr <= 0x401F:
		//todo  0x4000...0x401F => self.apu.write(addr - 0x4000, data),
	case 0x8000 <= addr && addr <= 0xFFFF:
		//todo
		//println!("switch bank to {}", data);
		//self.mmc.set_bank(data);
	}
}

