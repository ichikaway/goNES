package bus

type Rom struct {
	rom []byte
}

func NewRom(data []byte) Rom {
	return Rom{rom: data}
}

func (rom Rom) Read(addr uint16) byte {
	return rom.rom[addr]
}

func (rom Rom) Size() int {
	return len(rom.rom)
}
