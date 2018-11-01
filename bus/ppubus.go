package bus

type PpuBus struct {
	CharacterRam Ram
}

func NewPpuBus(characterRam *Ram) PpuBus {
		return PpuBus{CharacterRam: *characterRam}
}

func (ram PpuBus) readByPpu(addr int) byte {
	return ram.CharacterRam.Read(addr)
}

func (ram *PpuBus) writeByPpu(addr int, val byte) {
	ram.CharacterRam.Write(addr, val)
}
