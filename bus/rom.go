package bus

type Rom struct {
	rom []byte
}

func NewRom(data []byte) Rom {
	return Rom{rom: data}
}

func (rom Rom) Read(addr int) byte {
	return rom.rom[addr]
}

func (rom Rom) size() int {
	return len(rom.rom)
}
