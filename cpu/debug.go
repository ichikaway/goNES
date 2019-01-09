package cpu

import (
	"fmt"
)

func printAddressingMode(mode int) {
	list := [...]string{
		"Immediate",
		"ZeroPage",
		"Relative",
		"Implied",
		"Absolute",
		"Accumulator",
		"ZeroPageX",
		"ZeroPageY",
		"AbsoluteX",
		"AbsoluteY",
		"PreIndexedIndirect",
		"PostIndexedIndirect",
		"IndirectAbsolute",
	}

	fmt.Println("debug mode: ", list[mode], " , int: " ,mode)
}

func getOpecodeName(opecode int) string {
		ope := [...]string{
		"LDA",
		"LDX",
		"LDY",
		"STA",
		"STX",
		"STY",
		"TXA",
		"TYA",
		"TXS",
		"TAY",
		"TAX",
		"TSX",
		"PHP",
		"PLP",
		"PHA",
		"PLA",
		"ADC",
		"SBC",
		"CPX",
		"CPY",
		"CMP",
		"AND",
		"EOR",
		"ORA",
		"BIT",
		"ASL",
		"LSR",
		"ROL",
		"ROR",
		"INX",
		"INY",
		"INC",
		"DEX",
		"DEY",
		"DEC",
		"CLC",
		"CLI",
		"CLV",
		"SEC",
		"SEI",
		"NOP",
		"BRK",
		"JSR",
		"JMP",
		"RTI",
		"RTS",
		"BPL",
		"BMI",
		"BVC",
		"BVS",
		"BCC",
		"BCS",
		"BNE",
		"BEQ",
		"SED",
		"CLD",
		"LAX",
		"SAX",
		"DCP",
		"ISB",
		"SLO",
		"RLA",
		"SRE",
		"RRA",
	}
	return ope[opecode]
}

func printOpecode(opecode int) {


	fmt.Println("debug opecode: ", getOpecodeName(opecode), " , int: " ,opecode)
}
