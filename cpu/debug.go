package cpu

import (
	"fmt"
)

func getAddressingMode(mode int) string {
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

	return list[mode]
}

func printAddressingMode(mode int) {

	fmt.Println("debug mode: ", getAddressingMode(mode), " , int: " ,mode)
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
