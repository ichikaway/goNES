package registers

type Registers struct {
	A  byte
	X  byte
	Y  byte
	SP byte
	PC uint16
	P  Status
}

type Status struct {
	Negative     bool
	Overflow     bool
	Reserved     bool
	Break_mode   bool
	Decimal_mode bool
	Interrupt    bool
	Zero         bool
	Carry        bool
}

func GetDefaultRegisters() Registers {
	return Registers{
		A:  0,
		X:  0,
		Y:  0,
		PC: 0x8000,
		SP: 0xFD,
		P: Status{
			Negative:     false,
			Overflow:     false,
			Reserved:     true,
			Break_mode:   true,
			Decimal_mode: false,
			Interrupt:    true,
			Zero:         false,
			Carry:        false,
		},
	}
}

func (registers *Registers) IncrementPc() {
	registers.PC++
}

func (registers Registers) GetPc() uint16 {
	return registers.PC
}

func UpdateNegativeBy(data byte) bool {
	if (data & 0x80) == 0x80 {
		return true
	}
	return false
}

func UpdateZeroBy(data byte) bool {
	if data == 0 {
		return true
	}
	return false
}
