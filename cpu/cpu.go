package cpu

import (
	"fmt"
	"goNES/cpu/registers"
	"goNES/cpu_interrupts"
	"goNES/cpubus"
)

type Cpu struct {
	CpuBus     cpubus.CpuBus
	Interrupts cpu_interrupts.Interrupts
	Registers  registers.Registers
	HasBranched bool
}

func NewCpu(cpubus cpubus.CpuBus, interrupts cpu_interrupts.Interrupts) Cpu {

	return Cpu{
		CpuBus:     cpubus,
		Interrupts: interrupts,
		Registers:  registers.GetDefaultRegisters(),
		HasBranched: false,

	}
}

func (cpu *Cpu) Reset() {
	cpu.Registers = registers.GetDefaultRegisters()
	cpu.Registers.PC = cpu.CpuBus.ReadWord(0xFFFC)
	fmt.Printf("Initialize pc: %04x\n", cpu.Registers.PC)
}