package bus

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Keypad struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

func NewKeypad() Keypad{
		return Keypad{}
}

func (key *Keypad) Update(buttons [8]bool) {
	key.buttons = buttons
}


func (key *Keypad) Read() byte {
	value := byte(0)
	if key.index < 8 && key.buttons[key.index] {
		value = 1
	}
	key.index++
	if key.strobe&1 == 1 {
		key.index = 0
	}
	return value
}

func (key *Keypad) Write(value byte) {
	key.strobe = value
	if key.strobe&1 == 1 {
		key.index = 0
	}
}
