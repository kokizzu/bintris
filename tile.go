package main

import (
	"strconv"
	"time"
)

type TileSet struct {
	id          int
	Size        int // Number of tiles
	Number      int
	NumberStr   string // Textual representation of the bin number.
	X           float32
	Y           float32
	Sprites     []*Sprite
	numberSlots []*Sprite
	Speed       float64
	ClickFunc   func(x, y float32)
	sizex       float32
	sizey       float32
	scale       float32
	clicked     time.Time
	objectType  ObjectType
	gh          *Game
	hidden      bool
	tileWidth   float32
	tileHeight  float32
}

func (t *TileSet) Init(scale float32, size int, number int, x, y float32, g *Game) {
	t.gh = g
	t.id = g.NewID()
	t.Speed = 200
	t.Number = number
	t.Size = size
	t.Y = g.Y(y)
	t.X = g.X(x)
	t.sizex = t.tileWidth
	t.sizey = t.tileHeight
	t.scale = scale
	t.tileWidth = g.X(241)
	t.tileHeight = g.Y(30)

	t.NumberStr = strconv.Itoa(number)
	c := &Sprite{hidden: false}
	t.Sprites = append(t.Sprites, c)

	c.Init(g.X(13), t.Y, 0.6, 6.0, 5.0, "tile", g)
	g.AddObjects(c)

	for x := 241 / 4; x < 241; x += 241 / 4 {
		s := g.tex.AddText("0", g.X(float32(x)-(241/8)+8), t.Y+30, 0.7, 4.5, 3.0, EffectMetaballsBlue)

		t.numberSlots = append(t.numberSlots, s...)
		s = g.tex.AddText("1", g.X(float32(x)-(241/8)+8), t.Y+30, 0.7, 4.5, 3.0, EffectMetaballsBlue)
		s[0].Hide()
		t.numberSlots = append(t.numberSlots, s...)
	}

	n := g.tex.AddText(t.NumberStr, g.X(270), g.Y(y+7), 0.7, 4.5, 3.0, EffectMetaballsBlue)
	t.Sprites = append(t.Sprites, n...)

	//t.ClickFunc = func(x, y float32) {
	//	c.Click(x, y)
	//}

}

func (s *TileSet) Hidden() bool {
	return s.hidden
}

func (t *TileSet) GetID() int {
	return t.id
}

func (t *TileSet) Click(x, y float32) {
	if time.Since(t.clicked) < time.Duration(100*time.Millisecond) {
		return
	}

	t.clicked = time.Now()
	slot0 := 0
	slot1 := 0
	if x > t.X && x < t.X+t.tileWidth/float32(t.Size) {
		slot0 = 0
		slot1 = 1
	} else if x > t.X+t.tileWidth/float32(t.Size) && x < t.X+t.tileWidth/float32(t.Size)*2 {
		slot0 = 2
		slot1 = 3
	} else if x > t.X+t.tileWidth/float32(t.Size)*2 && x < t.X+t.tileWidth/float32(t.Size)*3 {
		slot0 = 4
		slot1 = 5
	} else if x > t.X+t.tileWidth/float32(t.Size)*3 && x < t.X+t.tileWidth/float32(t.Size)*4 {
		slot0 = 6
		slot1 = 7
	}
	c0 := t.numberSlots[slot0]
	c1 := t.numberSlots[slot1]
	if c0.hidden {
		c0.hidden = false
		c1.hidden = true
	} else {
		c0.hidden = true
		c1.hidden = false
	}

	// Verify number
	t.VerifyNumber()
}

func (t *TileSet) VerifyNumber() {
	num := 0
	for i := 0; i < t.Size*2; i += 2 {
		if t.numberSlots[i].hidden {
			num |= (1 << (3 - (i / 2)))
		}
	}

	if num == t.Number {
		for i := range t.Sprites {
			t.gh.DeleteObject(*t.Sprites[i])
		}
		for i := range t.numberSlots {
			t.gh.DeleteObject(*t.numberSlots[i])
		}
		t.hidden = true
	}
}

func (t *TileSet) GetObjectType() ObjectType {
	return ObjectTypeTileSet
}

func (t *TileSet) Draw(dt float64) {
}

func (c *TileSet) GetX() float32 {
	return c.X
}

func (c *TileSet) GetY() float32 {
	return c.Y
}

func (t *TileSet) Update(dt float64) {
	if t.hidden {
		return
	}

	// Check for collissions
	if t.Y < t.gh.Y(9) {
		return
	}

	for _, o := range t.gh.objects {
		if o.GetObjectType() == ObjectTypeTileSet {
			if o.hidden {
				continue
			}
			if o.GetID() == t.id {
				continue
			}

			if int(t.Y-t.tileHeight-4) < int(o.GetY()) && int(o.GetY()) < int(t.Y) {
				return
			}
		}
	}

	t.Y -= float32(t.Speed * dt)
	for i := range t.Sprites {
		t.Sprites[i].ChangeY(-float32(t.Speed * dt))
	}
	for i := range t.numberSlots {
		t.numberSlots[i].ChangeY(-float32(t.Speed * dt))
	}
}

func (t *TileSet) Resize() {
	//a := float32(s.gh.size.WidthPx) / float32(s.gh.size.HeightPx)
	t.X *= float32(t.gh.size.WidthPx) / float32(t.gh.sizePrev.WidthPx)
	t.Y *= float32(t.gh.size.HeightPx) / float32(t.gh.sizePrev.HeightPx)
	t.Speed *= float64(float32(t.gh.size.HeightPx) / float32(t.gh.sizePrev.HeightPx))
}
