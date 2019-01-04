package ppu

type RenderingData struct {
	Palette    []byte
	Background Background
	Sprites    []SpriteWithAttribute
}

func (this RenderingData) isSetBackground() bool{
	if len(this.Background.Tiles) > 0 {
		return true
	}
	return false
}

func (this RenderingData) isSetSprites() bool{
	if len(this.Sprites) > 0 {
		return true
	}
	return false
}

/** @var int[] */
//public $palette;
/** @var \Nes\Ppu\Tile[] */
//public $background;
/** @var \Nes\Ppu\SpriteWithAttribute[] */
//public $sprites;

