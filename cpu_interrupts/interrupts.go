package cpu_interrupts

type Interrupts struct {
	Nmi bool
	Irq bool
}

func NewInterrupts() Interrupts {
	return Interrupts{Nmi:false, Irq:false}
}

func (this Interrupts) isNmiAssert() bool {
	return this.Nmi;
}


func (this Interrupts) isIrqAssert() bool {
	return this.Irq;
}

func (this *Interrupts) assertNmi() {
	this.Nmi = true
}

func (this *Interrupts) deassertNmi() {
	this.Nmi = false
}

func (this *Interrupts) assertIrq(){
	this.Nmi = true
}

func (this *Interrupts) deassertIrq() {
	this.Nmi = false
}
