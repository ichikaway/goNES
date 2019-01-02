package ppu

type Tile struct {
	tile       []Sprite
	scroll_x   byte
	scroll_y   byte
	paletteId  int
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
