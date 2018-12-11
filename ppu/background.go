package ppu

type BackgroundCtx struct {
	tile       []int //仮実装  Tileクラスのスライスを本当は管理する
	scroll_x   byte
	scroll_y   byte
	is_enabled bool
}

type Background struct {
	Vec []BackgroundCtx
}

func NewBackground() Background {
	return Background{}
}

func (this *Background) Clear() {
	this.Vec = []BackgroundCtx{}
}
