package ppu

type SpriteWithAttribute struct {
	SpriteArry Sprite
	X          byte
	Y          byte
	Attribute  byte
	SpriteId   byte
	IsSet      bool //Spriteがセットされているか判定するgoNes独自パラメータ
}

func NewStripeWithAttribute(sprite Sprite, x byte, y byte, attribute byte, spriteId byte) SpriteWithAttribute {
	return SpriteWithAttribute{
		SpriteArry: sprite,
		X:          x,
		Y:          y,
		Attribute:  attribute,
		SpriteId:   spriteId,
		IsSet:      true,
	}
}
