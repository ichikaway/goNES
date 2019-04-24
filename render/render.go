package render

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"goNES/ppu"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

type Renderer struct {
	FrameBuffer [256 * 256 * 4]byte
	Serial      int
}

func NewRenderer() Renderer {
	return Renderer{}
}

func getColors() [64][3]byte {
	colors := [64][3]byte{
		{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96},
		{0xA1, 0x00, 0x5E}, {0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00},
		{0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00}, {0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E},
		{0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05}, {0x05, 0x05, 0x05},
		{0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
		{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00},
		{0xC4, 0x62, 0x00}, {0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55},
		{0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21}, {0x09, 0x09, 0x09}, {0x09, 0x09, 0x09},
		{0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF}, {0xD4, 0x80, 0xFF},
		{0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
		{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4},
		{0x05, 0xFB, 0xFF}, {0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D},
		{0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF}, {0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB},
		{0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0}, {0xFF, 0xEF, 0xA6},
		{0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
		{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},
	}
	return colors
}


func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}

func (this Renderer) drawDot() {

	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault

	termbox.Clear(color, backgroundColor)

	drawLine(0, 0, "Press ESC to exit." + string("aa"))
	drawLine(2, 1, fmt.Sprintf("date: %v", time.Now()))

	width := 256
	height := 224

	runes := []rune(". ")

	for y := 2; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (x + (y * 0x100)) * 4

			runeVal := runes[1]
			for i := 0 ; i < 3 ; i++ {
				if this.FrameBuffer[index + i] >= 128 {
					runeVal = runes[0]
				}
			}
			termbox.SetCell(x, y, runeVal, color, backgroundColor)
		}
	}

	termbox.Flush()
}

func (this Renderer) drawPng() {
	width := 256
	height := 224

	img := image.NewRGBA(image.Rect(0,0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (x + (y * 0x100)) * 4

			color := color.RGBA{
				this.FrameBuffer[index],
				this.FrameBuffer[index + 1],
				this.FrameBuffer[index + 2],
				255,
			}
			img.Set(x, y, color)
		}
	}


	file, _ := os.Create("output.png")
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func (this *Renderer) Render(data ppu.RenderingData) {
	//fmt.Println("tiles num:" , len(data.Background.Tiles))
	//fmt.Println("palette:" , data.Palette)

	//for i := range this.FrameBuffer {
	//	this.FrameBuffer[i] = 200
	//}

	if data.IsSetBackground() {
		this.renderBackground(data.Background, data.Palette)
	}

	if data.IsSetSprites() {
		//fmt.Println(data.Sprites)
		this.renderSprites(data.Sprites, data.Palette, data.Background)
	}

	/*
	fmt.Println("count: ",len(this.FrameBuffer))
	for index, d := range this.FrameBuffer {
		if d != 0 && d != 5 && d != 255{
			fmt.Println("index:", index, "data:",d)
		}
	}
	panic("")
	*/

	//this.drawPng()
	this.drawDot()
}

func (this *Renderer) renderBackground(background ppu.Background, palette []byte) {
	tiles := background.Tiles
	//fmt.Println(tiles)
	for i,tile := range tiles {
		x := (i % 33) * 8
		y := (i / 33) * 8

		/*
		flag := false
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if tile.Sprite[i][j] != 0 {
					flag = true
				}
			}
		}
		if flag {
			fmt.Println("i:",i, " x:",x, " y:", y)
		}
		*/

		this.renderTile(tile, x, y, palette)
	}
	/*
	for i := 0; i < len(tiles); i++ {
		x := (i % 33) * 8
		y := (i / 33) * 8
		this.renderTile(tiles[i], x, y, palette)
	}
	*/
}

func (this *Renderer) renderTile(tile ppu.Tile, tileX int, tileY int, palette []byte) {
	offsetX := int(tile.Scroll_x) % 8
	offsetY := int(tile.Scroll_y) % 8
	colors := getColors()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			//paletteIndex := (tile.PaletteId * 4) + int(tile.Sprite[i][j])
			//colorId := palette[paletteIndex]

			//let color_id = bg.tile.palette[bg.tile.sprite[i][j] as usize];
			colorId := palette[tile.Sprite[i][j]]

			color := colors[colorId]
			x := tileX + j - offsetX
			y := tileY + i - offsetY
			if x >= 0 && 0xFF >= x && y >= 0 && y < 224 {
				index := (x + (y * 0x100)) * 4
				//fmt.Println(color)
				//fmt.Println("index:", index)
				/*
				for _, cl := range color {
					if cl != 5 {
						//fmt.Println("cl:", cl, "")
					}
				}
				*/
				this.FrameBuffer[index] = color[0]
				this.FrameBuffer[index+1] = color[1]
				this.FrameBuffer[index+2] = color[2]
				this.FrameBuffer[index+3] = 0xFF
			}
		}
	}
}

func (this *Renderer) renderSprites(sprites []ppu.SpriteWithAttribute, palette []byte, background ppu.Background) {
	for _, sprite := range sprites {
		//fmt.Println(sprite)
		if sprite.IsSet {
			this.renderSprite(sprite, palette, background)
		}
	}
}

func (this Renderer) shouldPixelHide(x int, y int, background ppu.Background) bool {
	tileX := x / 8
	tileY := y / 8
	backgroundIndex := tileY*33 + tileX
	sprite := background.Tiles[backgroundIndex].Sprite
	// NOTE: If background pixel is not transparent, we need to hide sprite.
	return (sprite[y%8][x%8] % 4) != 0

	/**
		//rust実装
        let tile_x = x / 8;
        let tile_y = y / 8;
        let background_index = tile_y * 33 + tile_x;
        let sprite = &background[background_index];
        // NOTE: If background pixel is not transparent, we need to hide sprite.
        (sprite.tile.sprite[y % 8][x % 8] % 4) != 0
	 */
}

func (this *Renderer) renderSprite(sprite ppu.SpriteWithAttribute, palette []byte, background ppu.Background) {
	isVerticalReverse := (sprite.Attribute & 0x80) == 0x80
	isHrizontalReverse := (sprite.Attribute & 0x40) == 0x40
	isLowPriority := (sprite.Attribute & 0x20) == 0x20

	//paletteId := sprite.Attribute & 0x03
	colors := getColors()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			xHrizontal := j
			if isHrizontalReverse {
				xHrizontal = 7 - j
			}
			yVertical := i
			if isVerticalReverse {
				yVertical = 7 - i
			}
			x := int(sprite.X) + xHrizontal
			y := int(sprite.Y) + yVertical

			if isLowPriority && this.shouldPixelHide(x, y, background) {
				continue
			}

			if sprite.SpriteArry[i][j] != 0 {
				//colorId := palette[(paletteId*4)+sprite.SpriteArry[i][j]+0x10]
				colorId := palette[sprite.SpriteArry[i][j]]
				fmt.Println(colorId)
				color := colors[colorId]
				index := (x + (y * 0x100)) * 4
				this.FrameBuffer[index] = color[0]
				this.FrameBuffer[index+1] = color[1]
				this.FrameBuffer[index+2] = color[2]
			}
		}
	}
}
