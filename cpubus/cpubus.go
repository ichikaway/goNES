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


func (this CpuBus) ReadWord(addr uint16) uint16 {
	/**
	fn read_word(&mut self, addr: u16) -> u16 {
		let lower = self.read(addr) as u16;
		let upper = self.read(addr + 1) as u16;
		(upper << 8 | lower) as u16
	}
	 */
	lower := uint16(this.ReadByCpu(addr))
	upper := uint16(this.ReadByCpu(addr + 1))

	// 下記のようにupper、lowerをuint16に変換せずbyte型のまま8bitシフトさせると
	// 10000000が00000000になってしまうため事前にuint16にしておかないといけない
	//lower := this.ReadByCpu(addr)
	//upper := this.ReadByCpu(addr + 1)
	//fmt.Printf("%b\n",lower)
	//fmt.Printf("%b\n",upper)
	//fmt.Printf("%b\n",upper << 8)
	//fmt.Printf("%b\n", upper << 8 | lower)
	return uint16(upper << 8 | lower)
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

func (this CpuBus) ReadByCpu(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x0800:
		return this.Ram.Read(addr)
	case 0x0800 <= addr && addr <= 0x1FFF:
		return this.Ram.Read(addr - 0x0800)
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


func (this *CpuBus) WriteByCpu(addr uint16, data byte) {
	switch {
	case 0x0000 <= addr && addr < 0x0800:
		this.Ram.Write(addr, data)
	case 0x0800 <= addr && addr <= 0x1FFF:
		this.Ram.Write(addr - 0x0800, data)
	case 0x2000 <= addr && addr <= 0x3FFF:
		this.Ppu.Write(addr-0x2000, data)
	case addr == 0x4014:
		this.Dma.Write(data)
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

