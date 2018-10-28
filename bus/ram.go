package bus

type Ram struct {
	ram []byte
}

func NewRam(size int) Ram {
	return Ram{ram: make([]byte, size)}
}

func (ram *Ram) Reset() {
	ram.ram = make([]byte, len(ram.ram))
}

func (ram Ram) Read(addr int) byte {
	return ram.ram[addr]
}

func (ram *Ram) Write(addr int, val byte) {
	ram.ram[addr] = val
}
