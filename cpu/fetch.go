package cpu

import "fmt"

func (cpu *Cpu) fetchByte() byte {
	data := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()
	return data
}

func (cpu *Cpu) fetchWord() uint16 {
	fmt.Println("fetchWord PC1: ", cpu.Registers.GetPc())
	lower := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()

	fmt.Println("fetchWord PC2: ", cpu.Registers.GetPc())
	upper := cpu.CpuBus.ReadByCpu(cpu.Registers.GetPc())
	cpu.Registers.IncrementPc()


	fmt.Println("fetchWord lower: ", lower)
	fmt.Println("fetchWord upper: ", upper)
	fmt.Println("fetchWord upper shift: ", uint16(upper) << 8)
	fmt.Println("fetchWord result: ", uint16(upper) << 8 | uint16(lower))

	// upper, lower共に uint16にしておかないといけない
	// upperで8bitシフトするのでuint16の型にしてからシフトしないと0になるため
	return uint16(upper) << 8 | uint16(lower)
}
