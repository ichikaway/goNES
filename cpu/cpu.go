package cpu

import (
	"fmt"
	"goNES/cpu/registers"
	"goNES/cpu_interrupts"
	"goNES/cpubus"
	"goNES/util"
)

type Cpu struct {
	CpuBus      cpubus.CpuBus
	Interrupts  cpu_interrupts.Interrupts
	Registers   registers.Registers
	HasBranched bool
}

func NewCpu(cpubus cpubus.CpuBus, interrupts cpu_interrupts.Interrupts) Cpu {

	return Cpu{
		CpuBus:      cpubus,
		Interrupts:  interrupts,
		Registers:   registers.GetDefaultRegisters(),
		HasBranched: false,
	}
}

func (cpu *Cpu) Reset() {
	cpu.Registers = registers.GetDefaultRegisters()
	cpu.Registers.PC = cpu.CpuBus.ReadWord(0xFFFC)
	fmt.Printf("Initialize pc: %04x\n", cpu.Registers.PC)
}

/**
 * Pレジスタのステータス構造体の内容をbit変換してpushする
 * Pレジスタ構造体の中では、boolで管理されているため、それをuint8にキャストしてから8bitに合うようにシフトさせ論理和をとる
 * Negativeがtrueなら、10000000となり、Overflowがtrueなら01000000となる。ORをとると11000000となる。
 */
func (cpu *Cpu) pushStatus() {
	p := cpu.Registers.P
	status := util.Bool2Uint8(p.Negative) << 7 | util.Bool2Uint8(p.Overflow) << 6 |
		util.Bool2Uint8(p.Reserved) << 5 | util.Bool2Uint8(p.Break_mode) << 4 |
		util.Bool2Uint8(p.Decimal_mode) << 3 | util.Bool2Uint8(p.Interrupt) << 2 |
		util.Bool2Uint8(p.Zero) << 1 | util.Bool2Uint8(p.Carry)
	cpu.push(status)
}


func (cpu *Cpu) push(data byte) {
	addr := 0x0100 | int16(cpu.Registers.SP&0xFF)
	cpu.write(int(addr), data)
	cpu.Registers.SP--
}

func (cpu *Cpu) write(addr int, data byte) {
	cpu.CpuBus.WriteByCpu(addr, data)
}

