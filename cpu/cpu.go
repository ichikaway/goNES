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
