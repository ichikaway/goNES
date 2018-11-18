package bus

type PpuBus struct {
	CharacterRam Ram
}

func NewPpuBus(characterRam *Ram) PpuBus {
		return PpuBus{CharacterRam: *characterRam}
}

func (ram PpuBus) readByPpu(addr uint16) byte {
	return ram.CharacterRam.Read(addr)
}

func (ram *PpuBus) writeByPpu(addr uint16, val byte) {
	ram.CharacterRam.Write(addr, val)
}
