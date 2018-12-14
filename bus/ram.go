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

func (ram Ram) Read(addr uint16) byte {
	return ram.ram[addr]
}

func (ram *Ram) Write(addr uint16, val byte) {
	ram.ram[addr] = val
}


func (ram Ram) Size () int {
	return len(ram.ram)
}
