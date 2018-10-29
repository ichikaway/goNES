package bus

type CharacterRam struct {
	CharacterRam Ram
}

func NewPpuBus(characterRam *Ram) CharacterRam {
		return CharacterRam{CharacterRam: *characterRam}
}

func (ram CharacterRam) readByPpu(addr int) byte {
	return ram.CharacterRam.Read(addr)
}

func (ram *CharacterRam) writeByPpu(addr int, val byte) {
	ram.CharacterRam.Write(addr, val)
}
