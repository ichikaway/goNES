package nes

import (
	"github.com/nsf/termbox-go"
	"goNES/bus"
	"goNES/cpu"
	"goNES/cpu_interrupts"
	"goNES/cpubus"
	"goNES/dma"
	"goNES/ppu"
	"goNES/render"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

type Rom struct {
	ProgramRom         []byte
	CharacterRom       []byte
	isHorizontalMirror bool
	mapper             uint8
}

type Nes struct {
	Rom          Rom
	Ram          bus.Ram
	characterMem bus.Ram
	ProgramRom   bus.Rom
	PpuBus       bus.PpuBus
	Interrupts   cpu_interrupts.Interrupts
	Dma          dma.Dma
	Ppu          ppu.Ppu
	CpuBus       cpubus.CpuBus
	Cpu          cpu.Cpu
}

func New(data []byte) Nes {
	return Nes{Rom: parse(data)}
}

func (nes *Nes) Load() {
	nes.Ram = bus.NewRam(2048)
	nes.characterMem = bus.NewRam(0x4000)
	for i := 0; i < len(nes.Rom.CharacterRom); i++ {
		nes.characterMem.Write(uint16(i), nes.Rom.CharacterRom[i])
	}

	nes.ProgramRom = bus.NewRom(nes.Rom.ProgramRom)
	nes.PpuBus = bus.NewPpuBus(&nes.characterMem)
	nes.Interrupts = cpu_interrupts.NewInterrupts()

	nes.Ppu = ppu.NewPpu(nes.PpuBus, nes.Interrupts, nes.Rom.isHorizontalMirror)

	nes.Dma = dma.NewDma(nes.Ram, nes.Ppu)

	nes.CpuBus = cpubus.NewCpuBus(
		nes.Ram,
		nes.ProgramRom,
		nes.Ppu,
		nes.Dma,
	)

	nes.Cpu = cpu.NewCpu(nes.CpuBus, nes.Interrupts)
	nes.Cpu.Reset()
}

func (nes *Nes) frame(keyCh chan termbox.Key, frameCount *int, startTime time.Time) {
	for {
		select {
		case key := <-keyCh:
			mu.Lock()
			switch key {
			case termbox.KeyEsc, termbox.KeyCtrlC: //終了
				mu.Unlock()
				os.Exit(0)
			}
			mu.Unlock()
			//fmt.Println(key)
			break
		default:
			break
		}

		cycle := 0
		if nes.Dma.IsDmaProcessing() {
			nes.Dma.RunDma()
			cycle = 514
		}
		cycle += nes.Cpu.Run()

		if nes.Ppu.Run(cycle * 3) {
			*frameCount++
			renderer := render.NewRenderer(*frameCount, startTime)
			renderer.Render(nes.Ppu.RenderingData)
			break
		}
	}
}

func keyEvent(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}

func (nes Nes) Start() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	keyCh := make(chan termbox.Key)
	go keyEvent(keyCh)


	startTime := time.Now()
	frameCount := 0
	for {
		nes.frame(keyCh, &frameCount, startTime)
	}
}

