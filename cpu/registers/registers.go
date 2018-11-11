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
	negative     bool
	overflow     bool
	reserved     bool
	break_mode   bool
	decimal_mode bool
	interrupt    bool
	zero         bool
	carry        bool
}

func GetDefaultRegisters() Registers {
	return Registers{
		A:  0,
		X:  0,
		Y:  0,
		PC: 0x8000,
		SP: 0xFD,
		P: Status{
			negative:     false,
			overflow:     false,
			reserved:     true,
			break_mode:   true,
			decimal_mode: false,
			interrupt:    true,
			zero:         false,
			carry:        false,
		},
	}
}
