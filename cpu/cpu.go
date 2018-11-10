package cpu

import (
	"goNES/cpu_interrupts"
	"goNES/cpubus"
)

type Cpu struct {
	CpuBus cpubus.CpuBus
	Interrupts cpu_interrupts.Interrupts
}

func NewCpu(cpubus cpubus.CpuBus, interrupts cpu_interrupts.Interrupts) Cpu {

	return Cpu {
		CpuBus: cpubus,
		Interrupts: interrupts,
	}
}
