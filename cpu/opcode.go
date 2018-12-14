package cpu

type Opcode struct {
	name  int
	mode  int
	cycle int
}

const (
	Immediate = iota
	ZeroPage
	Relative
	Implied
	Absolute
	Accumulator
	ZeroPageX
	ZeroPageY
	AbsoluteX
	AbsoluteY
	PreIndexedIndirect
	PostIndexedIndirect
	IndirectAbsolute
)

const (
	LDA = iota
	LDX
	LDY
	STA
	STX
	STY
	TXA
	TYA
	TXS
	TAY
	TAX
	TSX
	PHP
	PLP
	PHA
	PLA
	ADC
	SBC
	CPX
	CPY
	CMP
	AND
	EOR
	ORA
	BIT
	ASL
	LSR
	ROL
	ROR
	INX
	INY
	INC
	DEX
	DEY
	DEC
	CLC
	CLI
	CLV
	SEC
	SEI
	NOP
	BRK
	JSR
	JMP
	RTI
	RTS
	BPL
	BMI
	BVC
	BVS
	BCC
	BCS
	BNE
	BEQ
	SED
	CLD
	LAX
	SAX
	DCP
	ISB
	SLO
	RLA
	SRE
	RRA
)

func getOpCodes() map[byte]Opcode {
	cycles := []int {
		7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
		6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
		6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
		6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
		2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
		2, 6, 2, 6, 4, 4, 4, 4, 2, 4, 2, 5, 5, 4, 5, 5,
		2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
		2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
		2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
		2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
		2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	}

	m := map[byte]Opcode{
		0xA9: Opcode{name: LDA, mode: Immediate, cycle: cycles[0xA9]},
		0xA5: Opcode{name: LDA, mode: ZeroPage, cycle: cycles[0xA5]},
		0xB5: Opcode{name: LDA, mode: ZeroPageX, cycle: cycles[0xB5]},
		0xAD: Opcode{name: LDA, mode: Absolute, cycle: cycles[0xAD]},
		0xBD: Opcode{name: LDA, mode: AbsoluteX, cycle: cycles[0xBD]},
		0xB9: Opcode{name: LDA, mode: AbsoluteY, cycle: cycles[0xB9]},
		0xA1: Opcode{name: LDA, mode: PreIndexedIndirect, cycle: cycles[0xA1]},
		0xB1: Opcode{name: LDA, mode: PostIndexedIndirect, cycle: cycles[0xB1]},
		0xA2: Opcode{name: LDX, mode: Immediate, cycle: cycles[0xA2]},
		0xA6: Opcode{name: LDX, mode: ZeroPage, cycle: cycles[0xA6]},
		0xAE: Opcode{name: LDX, mode: Absolute, cycle: cycles[0xAE]},
		0xB6: Opcode{name: LDX, mode: ZeroPageY, cycle: cycles[0xB6]},
		0xBE: Opcode{name: LDX, mode: AbsoluteY, cycle: cycles[0xBE]},
		0xA0: Opcode{name: LDY, mode: Immediate, cycle: cycles[0xA0]},
		0xA4: Opcode{name: LDY, mode: ZeroPage, cycle: cycles[0xA4]},
		0xAC: Opcode{name: LDY, mode: Absolute, cycle: cycles[0xAC]},
		0xB4: Opcode{name: LDY, mode: ZeroPageX, cycle: cycles[0xB4]},
		0xBC: Opcode{name: LDY, mode: AbsoluteX, cycle: cycles[0xBC]},
		0x85: Opcode{name: STA, mode: ZeroPage, cycle: cycles[0x85]},
		0x8D: Opcode{name: STA, mode: Absolute, cycle: cycles[0x8D]},
		0x95: Opcode{name: STA, mode: ZeroPageX, cycle: cycles[0x95]},
		0x9D: Opcode{name: STA, mode: AbsoluteX, cycle: cycles[0x9D]},
		0x99: Opcode{name: STA, mode: AbsoluteY, cycle: cycles[0x99]},
		0x81: Opcode{name: STA, mode: PreIndexedIndirect, cycle: cycles[0x81]},
		0x91: Opcode{name: STA, mode: PostIndexedIndirect, cycle: cycles[0x91]},
		0x86: Opcode{name: STX, mode: ZeroPage, cycle: cycles[0x86]},
		0x8E: Opcode{name: STX, mode: Absolute, cycle: cycles[0x8E]},
		0x96: Opcode{name: STX, mode: ZeroPageY, cycle: cycles[0x96]},
		0x84: Opcode{name: STY, mode: ZeroPage, cycle: cycles[0x84]},
		0x8C: Opcode{name: STY, mode: Absolute, cycle: cycles[0x8C]},
		0x94: Opcode{name: STY, mode: ZeroPageX, cycle: cycles[0x94]},
		0x8A: Opcode{name: TXA, mode: Implied, cycle: cycles[0x8A]},
		0x98: Opcode{name: TYA, mode: Implied, cycle: cycles[0x98]},
		0x9A: Opcode{name: TXS, mode: Implied, cycle: cycles[0x9A]},
		0xA8: Opcode{name: TAY, mode: Implied, cycle: cycles[0xA8]},
		0xAA: Opcode{name: TAX, mode: Implied, cycle: cycles[0xAA]},
		0xBA: Opcode{name: TSX, mode: Implied, cycle: cycles[0xBA]},
		0x08: Opcode{name: PHP, mode: Implied, cycle: cycles[0x08]},
		0x28: Opcode{name: PLP, mode: Implied, cycle: cycles[0x28]},
		0x48: Opcode{name: PHA, mode: Implied, cycle: cycles[0x48]},
		0x68: Opcode{name: PLA, mode: Implied, cycle: cycles[0x68]},
		0x69: Opcode{name: ADC, mode: Immediate, cycle: cycles[0x69]},
		0x65: Opcode{name: ADC, mode: ZeroPage, cycle: cycles[0x65]},
		0x6D: Opcode{name: ADC, mode: Absolute, cycle: cycles[0x6D]},
		0x75: Opcode{name: ADC, mode: ZeroPageX, cycle: cycles[0x75]},
		0x7D: Opcode{name: ADC, mode: AbsoluteX, cycle: cycles[0x7D]},
		0x79: Opcode{name: ADC, mode: AbsoluteY, cycle: cycles[0x79]},
		0x61: Opcode{name: ADC, mode: PreIndexedIndirect, cycle: cycles[0x61]},
		0x71: Opcode{name: ADC, mode: PostIndexedIndirect, cycle: cycles[0x71]},
		0xE9: Opcode{name: SBC, mode: Immediate, cycle: cycles[0xE9]},
		0xE5: Opcode{name: SBC, mode: ZeroPage, cycle: cycles[0xE5]},
		0xED: Opcode{name: SBC, mode: Absolute, cycle: cycles[0xED]},
		0xF5: Opcode{name: SBC, mode: ZeroPageX, cycle: cycles[0xF5]},
		0xFD: Opcode{name: SBC, mode: AbsoluteX, cycle: cycles[0xFD]},
		0xF9: Opcode{name: SBC, mode: AbsoluteY, cycle: cycles[0xF9]},
		0xE1: Opcode{name: SBC, mode: PreIndexedIndirect, cycle: cycles[0xE1]},
		0xF1: Opcode{name: SBC, mode: PostIndexedIndirect, cycle: cycles[0xF1]},
		0xE0: Opcode{name: CPX, mode: Immediate, cycle: cycles[0xE0]},
		0xE4: Opcode{name: CPX, mode: ZeroPage, cycle: cycles[0xE4]},
		0xEC: Opcode{name: CPX, mode: Absolute, cycle: cycles[0xEC]},
		0xC0: Opcode{name: CPY, mode: Immediate, cycle: cycles[0xC0]},
		0xC4: Opcode{name: CPY, mode: ZeroPage, cycle: cycles[0xC4]},
		0xCC: Opcode{name: CPY, mode: Absolute, cycle: cycles[0xCC]},
		0xC9: Opcode{name: CMP, mode: Immediate, cycle: cycles[0xC9]},
		0xC5: Opcode{name: CMP, mode: ZeroPage, cycle: cycles[0xC5]},
		0xCD: Opcode{name: CMP, mode: Absolute, cycle: cycles[0xCD]},
		0xD5: Opcode{name: CMP, mode: ZeroPageX, cycle: cycles[0xD5]},
		0xDD: Opcode{name: CMP, mode: AbsoluteX, cycle: cycles[0xDD]},
		0xD9: Opcode{name: CMP, mode: AbsoluteY, cycle: cycles[0xD9]},
		0xC1: Opcode{name: CMP, mode: PreIndexedIndirect, cycle: cycles[0xC1]},
		0xD1: Opcode{name: CMP, mode: PostIndexedIndirect, cycle: cycles[0xD1]},
		0x29: Opcode{name: AND, mode: Immediate, cycle: cycles[0x29]},
		0x25: Opcode{name: AND, mode: ZeroPage, cycle: cycles[0x25]},
		0x2D: Opcode{name: AND, mode: Absolute, cycle: cycles[0x2D]},
		0x35: Opcode{name: AND, mode: ZeroPageX, cycle: cycles[0x35]},
		0x3D: Opcode{name: AND, mode: AbsoluteX, cycle: cycles[0x3D]},
		0x39: Opcode{name: AND, mode: AbsoluteY, cycle: cycles[0x39]},
		0x21: Opcode{name: AND, mode: PreIndexedIndirect, cycle: cycles[0x21]},
		0x31: Opcode{name: AND, mode: PostIndexedIndirect, cycle: cycles[0x31]},
		0x49: Opcode{name: EOR, mode: Immediate, cycle: cycles[0x49]},
		0x45: Opcode{name: EOR, mode: ZeroPage, cycle: cycles[0x45]},
		0x4D: Opcode{name: EOR, mode: Absolute, cycle: cycles[0x4D]},
		0x55: Opcode{name: EOR, mode: ZeroPageX, cycle: cycles[0x55]},
		0x5D: Opcode{name: EOR, mode: AbsoluteX, cycle: cycles[0x5D]},
		0x59: Opcode{name: EOR, mode: AbsoluteY, cycle: cycles[0x59]},
		0x41: Opcode{name: EOR, mode: PreIndexedIndirect, cycle: cycles[0x41]},
		0x51: Opcode{name: EOR, mode: PostIndexedIndirect, cycle: cycles[0x51]},
		0x09: Opcode{name: ORA, mode: Immediate, cycle: cycles[0x09]},
		0x05: Opcode{name: ORA, mode: ZeroPage, cycle: cycles[0x05]},
		0x0D: Opcode{name: ORA, mode: Absolute, cycle: cycles[0x0D]},
		0x15: Opcode{name: ORA, mode: ZeroPageX, cycle: cycles[0x15]},
		0x1D: Opcode{name: ORA, mode: AbsoluteX, cycle: cycles[0x1D]},
		0x19: Opcode{name: ORA, mode: AbsoluteY, cycle: cycles[0x19]},
		0x01: Opcode{name: ORA, mode: PreIndexedIndirect, cycle: cycles[0x01]},
		0x11: Opcode{name: ORA, mode: PostIndexedIndirect, cycle: cycles[0x11]},
		0x24: Opcode{name: BIT, mode: ZeroPage, cycle: cycles[0x24]},
		0x2C: Opcode{name: BIT, mode: Absolute, cycle: cycles[0x2C]},
		0x0A: Opcode{name: ASL, mode: Accumulator, cycle: cycles[0x0A]},
		0x06: Opcode{name: ASL, mode: ZeroPage, cycle: cycles[0x06]},
		0x0E: Opcode{name: ASL, mode: Absolute, cycle: cycles[0x0E]},
		0x16: Opcode{name: ASL, mode: ZeroPageX, cycle: cycles[0x16]},
		0x1E: Opcode{name: ASL, mode: AbsoluteX, cycle: cycles[0x1E]},
		0x4A: Opcode{name: LSR, mode: Accumulator, cycle: cycles[0x4A]},
		0x46: Opcode{name: LSR, mode: ZeroPage, cycle: cycles[0x46]},
		0x4E: Opcode{name: LSR, mode: Absolute, cycle: cycles[0x4E]},
		0x56: Opcode{name: LSR, mode: ZeroPageX, cycle: cycles[0x56]},
		0x5E: Opcode{name: LSR, mode: AbsoluteX, cycle: cycles[0x5E]},
		0x2A: Opcode{name: ROL, mode: Accumulator, cycle: cycles[0x2A]},
		0x26: Opcode{name: ROL, mode: ZeroPage, cycle: cycles[0x26]},
		0x2E: Opcode{name: ROL, mode: Absolute, cycle: cycles[0x2E]},
		0x36: Opcode{name: ROL, mode: ZeroPageX, cycle: cycles[0x36]},
		0x3E: Opcode{name: ROL, mode: AbsoluteX, cycle: cycles[0x3E]},
		0x6A: Opcode{name: ROR, mode: Accumulator, cycle: cycles[0x6A]},
		0x66: Opcode{name: ROR, mode: ZeroPage, cycle: cycles[0x66]},
		0x6E: Opcode{name: ROR, mode: Absolute, cycle: cycles[0x6E]},
		0x76: Opcode{name: ROR, mode: ZeroPageX, cycle: cycles[0x76]},
		0x7E: Opcode{name: ROR, mode: AbsoluteX, cycle: cycles[0x7E]},
		0xE8: Opcode{name: INX, mode: Implied, cycle: cycles[0xE8]},
		0xC8: Opcode{name: INY, mode: Implied, cycle: cycles[0xC8]},
		0xE6: Opcode{name: INC, mode: ZeroPage, cycle: cycles[0xE6]},
		0xEE: Opcode{name: INC, mode: Absolute, cycle: cycles[0xEE]},
		0xF6: Opcode{name: INC, mode: ZeroPageX, cycle: cycles[0xF6]},
		0xFE: Opcode{name: INC, mode: AbsoluteX, cycle: cycles[0xFE]},
		0xCA: Opcode{name: DEX, mode: Implied, cycle: cycles[0xCA]},
		0x88: Opcode{name: DEY, mode: Implied, cycle: cycles[0x88]},
		0xC6: Opcode{name: DEC, mode: ZeroPage, cycle: cycles[0xC6]},
		0xCE: Opcode{name: DEC, mode: Absolute, cycle: cycles[0xCE]},
		0xD6: Opcode{name: DEC, mode: ZeroPageX, cycle: cycles[0xD6]},
		0xDE: Opcode{name: DEC, mode: AbsoluteX, cycle: cycles[0xDE]},
		0x18: Opcode{name: CLC, mode: Implied, cycle: cycles[0x18]},
		0x58: Opcode{name: CLI, mode: Implied, cycle: cycles[0x58]},
		0xB8: Opcode{name: CLV, mode: Implied, cycle: cycles[0xB8]},
		0x38: Opcode{name: SEC, mode: Implied, cycle: cycles[0x38]},
		0x78: Opcode{name: SEI, mode: Implied, cycle: cycles[0x78]},
		0xEA: Opcode{name: NOP, mode: Implied, cycle: cycles[0xEA]},
		0x00: Opcode{name: BRK, mode: Implied, cycle: cycles[0x00]},
		0x20: Opcode{name: JSR, mode: Absolute, cycle: cycles[0x20]},
		0x4C: Opcode{name: JMP, mode: Absolute, cycle: cycles[0x4C]},
		0x6C: Opcode{name: JMP, mode: IndirectAbsolute, cycle: cycles[0x6C]},
		0x40: Opcode{name: RTI, mode: Implied, cycle: cycles[0x40]},
		0x60: Opcode{name: RTS, mode: Implied, cycle: cycles[0x60]},
		0x10: Opcode{name: BPL, mode: Relative, cycle: cycles[0x10]},
		0x30: Opcode{name: BMI, mode: Relative, cycle: cycles[0x30]},
		0x50: Opcode{name: BVC, mode: Relative, cycle: cycles[0x50]},
		0x70: Opcode{name: BVS, mode: Relative, cycle: cycles[0x70]},
		0x90: Opcode{name: BCC, mode: Relative, cycle: cycles[0x90]},
		0xB0: Opcode{name: BCS, mode: Relative, cycle: cycles[0xB0]},
		0xD0: Opcode{name: BNE, mode: Relative, cycle: cycles[0xD0]},
		0xF0: Opcode{name: BEQ, mode: Relative, cycle: cycles[0xF0]},
		0xF8: Opcode{name: SED, mode: Implied, cycle: cycles[0xF8]},
		0xD8: Opcode{name: CLD, mode: Implied, cycle: cycles[0xD8]},
		0x1A: Opcode{name: NOP, mode: Implied, cycle: cycles[0x1A]},
		0x3A: Opcode{name: NOP, mode: Implied, cycle: cycles[0x3A]},
		0x5A: Opcode{name: NOP, mode: Implied, cycle: cycles[0x5A]},
		0x7A: Opcode{name: NOP, mode: Implied, cycle: cycles[0x7A]},
		0xDA: Opcode{name: NOP, mode: Implied, cycle: cycles[0xDA]},
		0xFA: Opcode{name: NOP, mode: Implied, cycle: cycles[0xFA]},
		0x02: Opcode{name: NOP, mode: Implied, cycle: cycles[0x02]},
		0x12: Opcode{name: NOP, mode: Implied, cycle: cycles[0x12]},
		0x22: Opcode{name: NOP, mode: Implied, cycle: cycles[0x22]},
		0x32: Opcode{name: NOP, mode: Implied, cycle: cycles[0x32]},
		0x42: Opcode{name: NOP, mode: Implied, cycle: cycles[0x42]},
		0x52: Opcode{name: NOP, mode: Implied, cycle: cycles[0x52]},
		0x62: Opcode{name: NOP, mode: Implied, cycle: cycles[0x62]},
		0x72: Opcode{name: NOP, mode: Implied, cycle: cycles[0x72]},
		0x92: Opcode{name: NOP, mode: Implied, cycle: cycles[0x92]},
		0xB2: Opcode{name: NOP, mode: Implied, cycle: cycles[0xB2]},
		0xD2: Opcode{name: NOP, mode: Implied, cycle: cycles[0xD2]},
		0xF2: Opcode{name: NOP, mode: Implied, cycle: cycles[0xF2]},
		0x80: Opcode{name: NOP, mode: Implied, cycle: cycles[0x80]},
		0x82: Opcode{name: NOP, mode: Implied, cycle: cycles[0x82]},
		0x89: Opcode{name: NOP, mode: Implied, cycle: cycles[0x89]},
		0xC2: Opcode{name: NOP, mode: Implied, cycle: cycles[0xC2]},
		0xE2: Opcode{name: NOP, mode: Implied, cycle: cycles[0xE2]},
		0x04: Opcode{name: NOP, mode: Implied, cycle: cycles[0x04]},
		0x44: Opcode{name: NOP, mode: Implied, cycle: cycles[0x44]},
		0x64: Opcode{name: NOP, mode: Implied, cycle: cycles[0x64]},
		0x14: Opcode{name: NOP, mode: Implied, cycle: cycles[0x14]},
		0x34: Opcode{name: NOP, mode: Implied, cycle: cycles[0x34]},
		0x54: Opcode{name: NOP, mode: Implied, cycle: cycles[0x54]},
		0x74: Opcode{name: NOP, mode: Implied, cycle: cycles[0x74]},
		0xD4: Opcode{name: NOP, mode: Implied, cycle: cycles[0xD4]},
		0xF4: Opcode{name: NOP, mode: Implied, cycle: cycles[0xF4]},
		0x0C: Opcode{name: NOP, mode: Implied, cycle: cycles[0x0C]},
		0x1C: Opcode{name: NOP, mode: Implied, cycle: cycles[0x1C]},
		0x3C: Opcode{name: NOP, mode: Implied, cycle: cycles[0x3C]},
		0x5C: Opcode{name: NOP, mode: Implied, cycle: cycles[0x5C]},
		0x7C: Opcode{name: NOP, mode: Implied, cycle: cycles[0x7C]},
		0xDC: Opcode{name: NOP, mode: Implied, cycle: cycles[0xDC]},
		0xFC: Opcode{name: NOP, mode: Implied, cycle: cycles[0xFC]},
		0xA7: Opcode{name: LAX, mode: ZeroPage, cycle: cycles[0xA7]},
		0xB7: Opcode{name: LAX, mode: ZeroPageY, cycle: cycles[0xB7]},
		0xAF: Opcode{name: LAX, mode: Absolute, cycle: cycles[0xAF]},
		0xBF: Opcode{name: LAX, mode: AbsoluteY, cycle: cycles[0xBF]},
		0xA3: Opcode{name: LAX, mode: PreIndexedIndirect, cycle: cycles[0xA3]},
		0xB3: Opcode{name: LAX, mode: PostIndexedIndirect, cycle: cycles[0xB3]},
		0x87: Opcode{name: SAX, mode: ZeroPage, cycle: cycles[0x87]},
		0x97: Opcode{name: SAX, mode: ZeroPageY, cycle: cycles[0x97]},
		0x8F: Opcode{name: SAX, mode: Absolute, cycle: cycles[0x8F]},
		0x83: Opcode{name: SAX, mode: PreIndexedIndirect, cycle: cycles[0x83]},
		0xEB: Opcode{name: SBC, mode: Immediate, cycle: cycles[0xEB]},
		0xC7: Opcode{name: DCP, mode: ZeroPage, cycle: cycles[0xC7]},
		0xD7: Opcode{name: DCP, mode: ZeroPageX, cycle: cycles[0xD7]},
		0xCF: Opcode{name: DCP, mode: Absolute, cycle: cycles[0xCF]},
		0xDF: Opcode{name: DCP, mode: AbsoluteX, cycle: cycles[0xDF]},
		0xDB: Opcode{name: DCP, mode: AbsoluteY, cycle: cycles[0xD8]},
		0xC3: Opcode{name: DCP, mode: PreIndexedIndirect, cycle: cycles[0xC3]},
		0xD3: Opcode{name: DCP, mode: PostIndexedIndirect, cycle: cycles[0xD3]},
		0xE7: Opcode{name: ISB, mode: ZeroPage, cycle: cycles[0xE7]},
		0xF7: Opcode{name: ISB, mode: ZeroPageX, cycle: cycles[0xF7]},
		0xEF: Opcode{name: ISB, mode: Absolute, cycle: cycles[0xEF]},
		0xFF: Opcode{name: ISB, mode: AbsoluteX, cycle: cycles[0xFF]},
		0xFB: Opcode{name: ISB, mode: AbsoluteY, cycle: cycles[0xF8]},
		0xE3: Opcode{name: ISB, mode: PreIndexedIndirect, cycle: cycles[0xE3]},
		0xF3: Opcode{name: ISB, mode: PostIndexedIndirect, cycle: cycles[0xF3]},
		0x07: Opcode{name: SLO, mode: ZeroPage, cycle: cycles[0x07]},
		0x17: Opcode{name: SLO, mode: ZeroPageX, cycle: cycles[0x17]},
		0x0F: Opcode{name: SLO, mode: Absolute, cycle: cycles[0x0F]},
		0x1F: Opcode{name: SLO, mode: AbsoluteX, cycle: cycles[0x1F]},
		0x1B: Opcode{name: SLO, mode: AbsoluteY, cycle: cycles[0x1B]},
		0x03: Opcode{name: SLO, mode: PreIndexedIndirect, cycle: cycles[0x03]},
		0x13: Opcode{name: SLO, mode: PostIndexedIndirect, cycle: cycles[0x13]},
		0x27: Opcode{name: RLA, mode: ZeroPage, cycle: cycles[0x27]},
		0x37: Opcode{name: RLA, mode: ZeroPageX, cycle: cycles[0x37]},
		0x2F: Opcode{name: RLA, mode: Absolute, cycle: cycles[0x2F]},
		0x3F: Opcode{name: RLA, mode: AbsoluteX, cycle: cycles[0x3F]},
		0x3B: Opcode{name: RLA, mode: AbsoluteY, cycle: cycles[0x3B]},
		0x23: Opcode{name: RLA, mode: PreIndexedIndirect, cycle: cycles[0x23]},
		0x33: Opcode{name: RLA, mode: PostIndexedIndirect, cycle: cycles[0x33]},
		0x47: Opcode{name: SRE, mode: ZeroPage, cycle: cycles[0x47]},
		0x57: Opcode{name: SRE, mode: ZeroPageX, cycle: cycles[0x57]},
		0x4F: Opcode{name: SRE, mode: Absolute, cycle: cycles[0x4F]},
		0x5F: Opcode{name: SRE, mode: AbsoluteX, cycle: cycles[0x5F]},
		0x5B: Opcode{name: SRE, mode: AbsoluteY, cycle: cycles[0x5B]},
		0x43: Opcode{name: SRE, mode: PreIndexedIndirect, cycle: cycles[0x43]},
		0x53: Opcode{name: SRE, mode: PostIndexedIndirect, cycle: cycles[0x53]},
		0x67: Opcode{name: RRA, mode: ZeroPage, cycle: cycles[0x67]},
		0x77: Opcode{name: RRA, mode: ZeroPageX, cycle: cycles[0x77]},
		0x6F: Opcode{name: RRA, mode: Absolute, cycle: cycles[0x6F]},
		0x7F: Opcode{name: RRA, mode: AbsoluteX, cycle: cycles[0x7F]},
		0x7B: Opcode{name: RRA, mode: AbsoluteY, cycle: cycles[0x7B]},
		0x63: Opcode{name: RRA, mode: PreIndexedIndirect, cycle: cycles[0x63]},
		0x73: Opcode{name: RRA, mode: PostIndexedIndirect, cycle: cycles[0x73]},
	}

	return m
}
