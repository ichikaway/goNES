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
	Opcode      map[byte]Opcode
}

func NewCpu(cpubus cpubus.CpuBus, interrupts cpu_interrupts.Interrupts) Cpu {

	return Cpu{
		CpuBus:      cpubus,
		Interrupts:  interrupts,
		Registers:   registers.GetDefaultRegisters(),
		HasBranched: false,
		Opcode:      getOpCodes(),
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
	status := util.Bool2Uint8(p.Negative)<<7 | util.Bool2Uint8(p.Overflow)<<6 |
		util.Bool2Uint8(p.Reserved)<<5 | util.Bool2Uint8(p.Break_mode)<<4 |
		util.Bool2Uint8(p.Decimal_mode)<<3 | util.Bool2Uint8(p.Interrupt)<<2 |
		util.Bool2Uint8(p.Zero)<<1 | util.Bool2Uint8(p.Carry)
	cpu.push(status)
}

func (cpu *Cpu) push(data byte) {
	addr := uint16(0x0100 | uint16(cpu.Registers.SP&0xFF))
	cpu.write(addr, data)
	cpu.Registers.SP--
}

func (cpu *Cpu) write(addr uint16, data byte) {
	cpu.CpuBus.WriteByCpu(addr, data)
}



func (cpu *Cpu) processIrq() {
	if cpu.Registers.P.Interrupt {
		return
	}
	cpu.Interrupts.DeassertIrq()
	cpu.Registers.P.Break_mode = false
	cpu.push(byte(cpu.Registers.PC >> 8))
	cpu.push(byte(cpu.Registers.PC))
	cpu.pushStatus()
	cpu.Registers.P.Interrupt = true
	cpu.Registers.PC = cpu.CpuBus.ReadWord(0xFFFE)
}

func (cpu *Cpu) processNmi() {
	cpu.Interrupts.DeassertNmi()
	cpu.Registers.P.Break_mode = false
	cpu.push(byte(cpu.Registers.PC >> 8))
	cpu.push(byte(cpu.Registers.PC))
	cpu.pushStatus()
	cpu.Registers.P.Interrupt = true
	cpu.Registers.PC = cpu.CpuBus.ReadWord(0xFFFA)
}

/**
	CPUのアドレッシングモードの説明
	http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#addressing
 */
func (cpu *Cpu) getAddrOrDataWithAdditionalCycle(mode int) (uint16, int){
	switch mode {
	case Accumulator:
		return 0x00, 0
	case Implied:
		return 0x00, 0
	case Immediate:
		return uint16(cpu.fetchByte()), 0
	case Relative:
		base := uint16(cpu.fetchByte())
		cycle := 0
		if base & 0xff00 != cpu.Registers.PC & 0xff00 {
			cycle = 1
		}
		if base < 0x80 {
			return base + cpu.Registers.GetPc(), cycle
		}
		return base + cpu.Registers.GetPc() - 256, cycle
	case ZeroPage:
		return uint16(cpu.fetchByte()), 0
	case ZeroPageX:
		addr := uint16(cpu.fetchByte())
		return (addr + uint16(cpu.Registers.X)) & 0xFF, 0
	case ZeroPageY:
		addr := uint16(cpu.fetchByte())
		return (addr + uint16(cpu.Registers.Y)) & 0xFF, 0
	case Absolute:
		return cpu.fetchWord(), 0
	case AbsoluteX:
		addr := cpu.fetchWord()
		cycle := 0
		if (addr & 0xFF00) != ((addr + uint16(cpu.Registers.X)) & 0xFF00) {
			cycle = 1
		}
		return (addr + uint16(cpu.Registers.X)) & 0xFFFF, cycle
	case AbsoluteY:
		addr := cpu.fetchWord()
		cycle := 0
		if (addr & 0xFF00) != ((addr + uint16(cpu.Registers.Y)) & 0xFF00) {
			cycle = 1
		}
		return (addr + uint16(cpu.Registers.Y)) & 0xFFFF, cycle
	case PreIndexedIndirect:
		baseAddr := uint16(cpu.fetchByte() + cpu.Registers.X) & 0xFF
		addr := (uint16(cpu.CpuBus.ReadByCpu(baseAddr)) + uint16(cpu.CpuBus.ReadByCpu(baseAddr+1)) & 0xFF) << 8
		cycle := 0
		if (addr & 0xFF00) != (baseAddr & 0xFF00) {
			cycle = 1
		}
		return addr & 0xFFFF, cycle
	case PostIndexedIndirect:
		data := uint16(cpu.fetchByte())
		baseAddr := uint16(cpu.CpuBus.ReadByCpu(data)) + uint16(cpu.CpuBus.ReadByCpu(data + 1)) & 0x00FF
		addr := baseAddr + uint16(cpu.Registers.Y)
		cycle := 0
		if (addr & 0xFF00) != (baseAddr & 0xFF00) {
			cycle = 1
		}
		return addr & 0xFFFF, cycle
	case IndirectAbsolute:
		addr := cpu.fetchWord()
		upper := uint16(cpu.CpuBus.ReadByCpu((addr & 0xFF00) | uint16(((addr & 0xFF) + 1) & 0xFF)))
		addr2 := uint16(cpu.CpuBus.ReadByCpu(addr)) + (upper << 8)
		return addr2 & 0xFFFF, 0
	default:
		panic("no match Addressing Mode")
	}
}

func (this *Cpu) execInstruction(opecode int, data uint16, mode int) {
	this.HasBranched = false

	switch opecode {
	case LDA:
		val := uint8(data)
		if mode != Immediate {
			val = this.CpuBus.ReadByCpu(data)
		}
		this.Registers.A = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case LDX:
		val := uint8(data)
		if mode != Immediate {
			val = this.CpuBus.ReadByCpu(data)
		}
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case LDY:
		val := uint8(data)
		if mode != Immediate {
			val = this.CpuBus.ReadByCpu(data)
		}
		this.Registers.Y = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case STA:
		this.write(data, this.Registers.A)
	case STX:
		this.write(data, this.Registers.X)
	case STY:
		this.write(data, this.Registers.Y)
	case TAX:
		val := this.Registers.A
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case TAY:
		val := this.Registers.A
		this.Registers.Y = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case TSX:
		val := this.Registers.SP
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case TXA:
		val := this.Registers.X
		this.Registers.A = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case TXS:
		this.Registers.SP = this.Registers.X
	case TYA:
		val := this.Registers.Y
		this.Registers.A = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case ADC:
		val := uint8(data)
		if mode != Immediate {
			val = this.CpuBus.ReadByCpu(data)
		}
		computed := val + this.Registers.A + util.Bool2Uint8(this.Registers.P.Carry)
		registerA := this.Registers.A

		this.Registers.P.Negative = registers.UpdateNegativeBy(computed)
		this.Registers.P.Zero = registers.UpdateZeroBy(computed)
		this.Registers.A = computed

		this.Registers.P.Carry = false
		if computed > 0xFF {
			this.Registers.P.Carry = true
		}

		this.Registers.P.Overflow = false
		if ((registerA ^ val) & 0x80) == 0 && ((registerA ^ computed) & 0x80) != 0 {
			this.Registers.P.Overflow = true
		}
	case AND:
		val := uint8(data)
		if mode != Immediate {
			val = this.CpuBus.ReadByCpu(data)
		}
		computed := this.Registers.A & val
		this.Registers.P.Negative = registers.UpdateNegativeBy(computed)
		this.Registers.P.Zero = registers.UpdateZeroBy(computed)
		this.Registers.A = computed
	}

}

func (cpu *Cpu) Run() int {
	if cpu.Interrupts.IsNmiAssert() {
		cpu.processNmi()
	}
	if cpu.Interrupts.IsIrqAssert() {
		cpu.processIrq()
	}

	opcode := cpu.fetchByte()
	opc := cpu.Opcode[opcode]
	fmt.Println(opc)
	data, additionalCycle := cpu.getAddrOrDataWithAdditionalCycle(opc.mode)
	fmt.Println(data, additionalCycle)

	cpu.execInstruction(opc.name, data, opc.mode)

	cycle := opc.cycle + additionalCycle
	if cpu.HasBranched {
		cycle++
	}

	return cycle
}
