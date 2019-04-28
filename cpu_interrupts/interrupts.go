package cpu_interrupts

type Interrupts struct {
	Nmi bool
	Irq bool
}

func NewInterrupts() *Interrupts {
	return &Interrupts{Nmi:false, Irq:false}
}

func (this Interrupts) IsNmiAssert() bool {
	return this.Nmi
}


func (this Interrupts) IsIrqAssert() bool {
	return this.Irq
}

func (this *Interrupts) AssertNmi() {
	this.Nmi = true
}

func (this *Interrupts) DeassertNmi() {
	this.Nmi = false
}

func (this *Interrupts) AssertIrq(){
	this.Irq = true
}

func (this *Interrupts) DeassertIrq() {
	this.Irq = false
}
