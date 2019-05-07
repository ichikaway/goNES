package cpu

import (
	"fmt"
	"goNES/cpu/registers"
	"goNES/cpu_interrupts"
	"goNES/cpubus"
	"goNES/util"
	)

type Cpu struct {
	CpuBus      *cpubus.CpuBus
	Interrupts  *cpu_interrupts.Interrupts
	Registers   *registers.Registers
	HasBranched bool
	Opcode      map[byte]Opcode
}

func NewCpu(cpubus *cpubus.CpuBus, interrupts *cpu_interrupts.Interrupts) *Cpu {

	return &Cpu{
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
	//fmt.Printf("Initialize pc: %04x\n", cpu.Registers.PC)
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

func (cpu *Cpu) popStatus() {
	val := cpu.pop()
	cpu.Registers.P.Negative = val & 0x80 == 0x80
	cpu.Registers.P.Overflow = val & 0x40 == 0x40
	cpu.Registers.P.Reserved = val & 0x20 == 0x20
	cpu.Registers.P.Break_mode = val & 0x10 == 0x10
	cpu.Registers.P.Decimal_mode = val & 0x08 == 0x08
	cpu.Registers.P.Interrupt = val & 0x04 == 0x04
	cpu.Registers.P.Zero = val & 0x02 == 0x02
	cpu.Registers.P.Carry = val & 0x01 == 0x01
}

func (cpu *Cpu) push(data byte) {
	addr := uint16(0x0100 | uint16(cpu.Registers.SP&0xFF))
	cpu.write(addr, data)
	cpu.Registers.SP--
}

func (cpu *Cpu) pop() byte {
	cpu.Registers.SP++
	addr := 0x0100 | uint16(cpu.Registers.SP)
	return cpu.read(addr)
}

func (cpu *Cpu) popPc() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	cpu.Registers.PC = upper << 8 | lower
}

func (cpu Cpu) read(addr uint16) byte {
	return cpu.CpuBus.ReadByCpu(addr)
}

func (cpu *Cpu) write(addr uint16, data byte) {
	cpu.CpuBus.WriteByCpu(addr, data)
}

func (cpu *Cpu) branch(addr uint16) {
	cpu.Registers.PC = addr
	cpu.HasBranched = true
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
		return 0x0000, 0
	case Implied:
		return 0x0000, 0
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
		baseAddr := uint16((cpu.fetchByte() + cpu.Registers.X) & 0xFF)
		addr := uint16(cpu.CpuBus.ReadByCpu(baseAddr)) + (uint16(cpu.CpuBus.ReadByCpu((baseAddr + 1) & 0xFF)) << 8)
		cycle := 0
		if (addr & 0xFF00) != (baseAddr & 0xFF00) {
			cycle = 1
		}
		return addr & 0xFFFF, cycle
	case PostIndexedIndirect:
		data := uint16(cpu.fetchByte())
		baseAddr := uint16(cpu.CpuBus.ReadByCpu(data)) + (uint16(cpu.CpuBus.ReadByCpu((data + 1) & 0x00FF)) * 0x100)
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

	//fmt.Println("base: ", opecode, "data:", data, "mode:", mode)
	//printOpecode(opecode)
	//printAddressingMode(mode)
	//fmt.Println("data: ", data)

	switch opecode {
	case LDA:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		this.Registers.A = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case LDX:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case LDY:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
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
		val := data
		if mode != Immediate {
			val = uint16(this.read(data))
		}
		registerA := this.Registers.A
		computed := val + uint16(registerA) + uint16(util.Bool2Uint8(this.Registers.P.Carry))

		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(computed))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(computed))

		this.Registers.P.Carry = false
		if computed > 0xFF {
			this.Registers.P.Carry = true
		}

		/*
		this.Registers.P.Overflow = false
		if !(((registerA ^ uint8(val)) & 0x80) != 0 && ((registerA ^ uint8(computed)) & 0x80) != 0) {
			this.Registers.P.Overflow = true
		}
		*/
		this.Registers.P.Overflow = !(((registerA ^ uint8(val)) & 0x80) != 0 && ((registerA ^ uint8(computed)) & 0x80) != 0)
		this.Registers.A = uint8(computed)

	case AND:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		computed := this.Registers.A & val
		this.Registers.P.Negative = registers.UpdateNegativeBy(computed)
		this.Registers.P.Zero = registers.UpdateZeroBy(computed)
		this.Registers.A = computed
	case ASL:
		if mode == Accumulator {
			acc := this.Registers.A
			shifted := uint8(acc << 1)

			this.Registers.P.Carry = (acc & 0x80) == 0x80
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.Registers.A = shifted
		} else {
			fetched := this.read(data)
			shifted := uint8(fetched << 1)

			this.Registers.P.Carry = (fetched & 0x80) == 0x80
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.write(data, shifted)
		}
	case BIT:
		val := this.read(data)
		acc := this.Registers.A

		//fmt.Println("BIT read data: ", val)
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(acc & val)
		this.Registers.P.Overflow = (val & 0x40) == 0x40
	case CMP:
		this.compare(data, mode, this.Registers.A)
	case CPX:
		this.compare(data, mode, this.Registers.X)
	case CPY:
		this.compare(data, mode, this.Registers.Y)
	case DEC:
		val := int8(this.read(data)) -1
		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(val))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(val))
		this.write(data, uint8(val))
	case DEX:
		val := int8(this.Registers.X) -1
		this.Registers.X = uint8(val)
		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(val))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(val))
	case DEY:
		val := int8(this.Registers.Y) -1
		this.Registers.Y = uint8(val)
		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(val))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(val))
	case EOR:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		computed := this.Registers.A ^ val
		this.Registers.P.Negative = registers.UpdateNegativeBy(computed)
		this.Registers.P.Zero = registers.UpdateZeroBy(computed)
		this.Registers.A = computed
	case INC:
		val := this.read(data) + 1
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
		this.write(data, val)
	case INX:
		val := this.Registers.X + 1
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case INY:
		val := this.Registers.Y + 1
		this.Registers.Y = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case LSR:
		if mode == Accumulator {
			acc := this.Registers.A
			shifted := uint8(acc >> 1)

			this.Registers.P.Carry = (acc & 0x01) == 0x01
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.Registers.A = shifted
		} else {
			fetched := this.read(data)
			shifted := uint8(fetched >> 1)

			this.Registers.P.Carry = (fetched & 0x01) == 0x01
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.write(data, shifted)
		}
	case ORA:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		computed := this.Registers.A | val
		this.Registers.P.Negative = registers.UpdateNegativeBy(computed)
		this.Registers.P.Zero = registers.UpdateZeroBy(computed)
		this.Registers.A = computed
	case ROL:
		if mode == Accumulator {
			acc := this.Registers.A
			shifted := rotateToLeft(this.Registers.P.Carry, acc)

			this.Registers.P.Carry = (acc & 0x80) == 0x80
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.Registers.A = shifted
		} else {
			fetched := this.read(data)
			shifted := rotateToLeft(this.Registers.P.Carry, fetched)

			this.Registers.P.Carry = (fetched & 0x80) == 0x80
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.write(data, shifted)
		}
	case ROR:
		if mode == Accumulator {
			acc := this.Registers.A
			shifted := rotateToRight(this.Registers.P.Carry, acc)

			this.Registers.P.Carry = (acc & 0x01) == 0x01
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.Registers.A = shifted
		} else {
			fetched := this.read(data)
			shifted := rotateToRight(this.Registers.P.Carry, fetched)

			this.Registers.P.Carry = (fetched & 0x01) == 0x01
			this.Registers.P.Negative = registers.UpdateNegativeBy(shifted)
			this.Registers.P.Zero = registers.UpdateZeroBy(shifted)
			this.write(data, shifted)
		}
	case SBC:
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		computed := int(this.Registers.A) - int(val) - int(util.Bool2Uint8(!this.Registers.P.Carry))
		registerA := this.Registers.A

		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(computed))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(computed))
		this.Registers.A = uint8(computed)

		this.Registers.P.Carry = false
		if computed >= 0 {
			this.Registers.P.Carry = true
		}

		this.Registers.P.Overflow = false
		if ((registerA ^ val) & 0x80) != 0 && ((registerA ^ uint8(computed)) & 0x80) != 0 {
			this.Registers.P.Overflow = true
		}
	case PHA:
		this.push(this.Registers.A)
	case PHP:
		this.Registers.P.Break_mode = true
		this.pushStatus()
	case PLA:
		val := this.pop()
		this.Registers.A = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case PLP:
		this.Registers.P.Reserved = true
		this.popStatus()
	case JMP:
		this.Registers.PC = data
	case JSR:
		pc := this.Registers.GetPc() - 1
		this.push(uint8(pc >> 8))
		this.push(uint8(pc))
		this.Registers.PC = data
	case RTS:
		//fmt.Println(this.Registers.PC)
		this.popPc()
		//fmt.Println(this.Registers.PC)
		this.Registers.IncrementPc()
		//fmt.Println(this.Registers.PC)
	case RTI:
		this.popStatus()
		this.popPc()
		this.Registers.P.Reserved = true
	case BCC:
		if !this.Registers.P.Carry {
			this.branch(data)
		}
	case BCS:
		if this.Registers.P.Carry {
			this.branch(data)
		}
	case BEQ:
		if this.Registers.P.Zero {
			this.branch(data)
		}
	case BMI:
		if this.Registers.P.Negative {
			this.branch(data)
		}
	case BNE:
		if !this.Registers.P.Zero {
			this.branch(data)
		}
	case BPL:
		if !this.Registers.P.Negative {
			this.branch(data)
		}
	case BVS:
		if this.Registers.P.Overflow {
			this.branch(data)
		}
	case BVC:
		if !this.Registers.P.Overflow {
			this.branch(data)
		}
	case CLD:
		this.Registers.P.Decimal_mode = false
	case CLC:
		this.Registers.P.Carry = false
	case CLI:
		this.Registers.P.Interrupt = false
	case CLV:
		this.Registers.P.Overflow = false
	case SEC:
		this.Registers.P.Carry = true
	case SEI:
		this.Registers.P.Interrupt = true
	case SED:
		this.Registers.P.Decimal_mode = true
	case BRK:
		//fmt.Println("PC1: ",this.Registers.PC)
		interrupt := this.Registers.P.Interrupt
		this.Registers.IncrementPc()
		//fmt.Println("PC2: ",this.Registers.PC)
		pc := this.Registers.GetPc()
		this.push(uint8(pc >> 8))
		this.push(uint8(pc))
		this.Registers.P.Break_mode = true
		this.pushStatus()
		this.Registers.P.Interrupt = true
		if !interrupt {
			fetched := this.CpuBus.ReadWord(0xFFFE)
			this.Registers.PC = fetched
		}
		this.Registers.DecrementPc()
		//fmt.Println("PC3: ",this.Registers.PC)
	case NOP:
		return
	//ここから下はunofficialな実装
	case LAX:
		val := this.read(data)
		this.Registers.A = val
		this.Registers.X = val
		this.Registers.P.Negative = registers.UpdateNegativeBy(val)
		this.Registers.P.Zero = registers.UpdateZeroBy(val)
	case SLO:
		val := this.read(data)
		this.Registers.P.Carry = (val & 0x80) == 0x80
		val2 := (val << 1) & 0xFF
		this.Registers.A = val2
		this.Registers.P.Negative = registers.UpdateNegativeBy(val2)
		this.Registers.P.Zero = registers.UpdateZeroBy(val2)
		this.write(data, val2)
	case RLA:
		return
	case SRE:
		return
	case RRA:
		return
	default:
		fmt.Println("opecode: ", opecode)
		printOpecode(opecode)
		panic("no instruction!")
	}

}

func rotateToRight(carry bool, data uint8) uint8 {
	//((v >> 1) as Data | if registers.get_carry() { 0x80 } else { 0x00 }) as Data
	v := data >> 1
	c := uint8(0x00)
	if carry {
		c = 0x80
	}
	return v | c
}

func rotateToLeft(carry bool, data uint8) uint8 {
	//((v << 1) as Data | if registers.get_carry() { 0x01 } else { 0x00 }) as Data
	v := data << 1
	c := uint8(0x00)
	if carry {
		c = 0x01
	}
	return v | c
}

func (this *Cpu) compare(data uint16, mode int, registerVal byte) {
		val := uint8(data)
		if mode != Immediate {
			val = this.read(data)
		}
		compared := int16(registerVal) - int16(val)
		this.Registers.P.Carry = compared >= 0
		this.Registers.P.Negative = registers.UpdateNegativeBy(uint8(compared))
		this.Registers.P.Zero = registers.UpdateZeroBy(uint8(compared))
}

func (cpu *Cpu) Run() int {
	//log.Println("------------- CPU run ---------------")
	//log.Println("run Pc: ", cpu.Registers.PC)

	//log.Println("cpu Nmi: ", cpu.Interrupts.Nmi)
	if cpu.Interrupts.IsNmiAssert() {
		//log.Println("processNmi()")
		cpu.processNmi()
	}
	if cpu.Interrupts.IsIrqAssert() {
		//log.Println("processIRQ()")
		cpu.processIrq()
	}

	opcode := cpu.fetchByte()
	//fmt.Println("ope: ", opcode)
	opc := cpu.Opcode[opcode]
	//fmt.Println(opc)
	data, additionalCycle := cpu.getAddrOrDataWithAdditionalCycle(opc.mode)
	//fmt.Println(data, additionalCycle)

	/*
	log.Println(
		"PC: ", cpu.Registers.GetPc(),
		" opcode: ", getOpecodeName(opc.name),
		" addr: ", data,
		" mode: ", getAddressingMode(opc.mode),
		)
	*/

	cpu.execInstruction(opc.name, data, opc.mode)

	cycle := opc.cycle + additionalCycle
	if cpu.HasBranched {
		cycle++
	}


	return cycle
}
