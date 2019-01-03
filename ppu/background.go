package ppu

type Tile struct {
	Sprite     Sprite
	Scroll_x   byte
	Scroll_y   byte
	PaletteId  int
}

type Background struct {
	Tiles []Tile
}

func NewBackground() Background {
	return Background{}
}

func (this *Background) Clear() {
	this.Tiles = []Tile{}
}
