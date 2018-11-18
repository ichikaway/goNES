package cpu

func (cpu *Cpu) fetchByte() byte {
	data := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()
	return data
}

func (cpu *Cpu) fetchWord() uint16 {
	lower := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()

	upper := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()

	return uint16(upper << 8 | lower)
}
