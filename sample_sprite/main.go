// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi
//
// very simple Sprite class implementaion.
// The next step would be implement node chain.

package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"math"
)

const (
	WIDTH  = 320
	HEIGHT = 240
)

const (
	STATE_INIT = iota + 1
	STATE_IDLE
)

var (
	state int
	tick  int

	bg     *Sprite
	gopher *Sprite
)

type Vector2 struct {
	x float64
	y float64
}

type RGBA8 struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func (p *RGBA8) Reset() {
	p.r = 255
	p.g = 255
	p.b = 255
	p.a = 255
}

func (p *RGBA8) FSetAlpha(a float64) {
	p.a = uint8(255.0 * a)
}

func (p *RGBA8) FGetColor() (r, g, b, a float64) {
	return float64(p.r) / 255.0, float64(p.g) / 255.0, float64(p.b) / 255.0, float64(p.a) / 255.0
}

// primitive Sprite.
type Sprite struct {
	anchor Vector2
	pos    Vector2
	rot    float64
	scale  Vector2

	rgba    RGBA8
	texture *ebiten.Image
}

func NewSprite() *Sprite {
	p := &Sprite{}
	p.anchor.x = .5
	p.anchor.y = .5
	p.scale.x = 1
	p.scale.y = 1
	p.rgba.Reset()
	return p
}

func (p *Sprite) SetTexture(tex *ebiten.Image) {
	p.texture = tex
}

func (p *Sprite) FSetPosition(x, y float64) {
	p.pos.x = x
	p.pos.y = y
}

func (p *Sprite) FSetScale(x, y float64) {
	p.scale.x = x
	p.scale.y = y
}

func (p *Sprite) FSetRotation(r float64) {
	p.rot = r
}

func (p *Sprite) FSetAlpha(a float64) {
	p.rgba.FSetAlpha(a)
}

func (p *Sprite) TextureSize() (w, h int) {
	if p.texture != nil {
		return p.texture.Size()
	}
	return 0, 0
}

func (p *Sprite) CalcL2W() ebiten.GeoM {
	var l2w ebiten.GeoM

	w, h := p.texture.Size()
	l2w.Translate(-(float64(w) * p.anchor.x), -(float64(h) * p.anchor.y))
	l2w.Scale(p.scale.x, p.scale.y)
	if p.rot != 0 {
		l2w.Rotate(p.rot)
	}
	l2w.Translate(p.pos.x, p.pos.y)
	return l2w
}

func (p *Sprite) Render(target *ebiten.Image) {
	if p.texture != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM = p.CalcL2W()
		op.ColorM.Scale(p.rgba.FGetColor())

		target.DrawImage(p.texture, op)
	}
}

func initialize() error {
	gopherTex, _, err := ebitenutil.NewImageFromFile("../assets/gophercolor.png", ebiten.FilterNearest)
	if err != nil {
		return err
	}

	gopher = NewSprite()
	gopher.SetTexture(gopherTex)

	bgTex, _, err := ebitenutil.NewImageFromFile("../assets/cal003.png", ebiten.FilterNearest)
	if err != nil {
		return err
	}
	bg = NewSprite()
	bg.SetTexture(bgTex)

	state = STATE_IDLE
	return nil
}

func update(screen *ebiten.Image) error {
	switch state {
	case STATE_INIT:
		initialize()
	case STATE_IDLE:
		if ebiten.IsRunningSlowly() {
			return nil
		}
		bg.FSetPosition(float64(WIDTH*0.5), float64(HEIGHT*0.5))
		bg.Render(screen)

		w, _ := gopher.TextureSize()
		gopher.FSetPosition(float64(tick%(WIDTH+w)-(w>>1)), HEIGHT*0.5)
		gopher.FSetScale(1.0, math.Sin(float64(tick%60)/60*math.Pi)*0.25+1) // 25% to 125%.
		gopher.FSetRotation((math.Sin(float64(tick%120)/120*math.Pi) - 0.5) * (math.Pi * 0.4))

		gopher.FSetAlpha(math.Sin(float64(tick%60)/60*math.Pi)*0.5 + 0.5)
		gopher.Render(screen)
	}
	tick++
	return nil
}

func main() {
	state = STATE_INIT
	tick = 0

	var err error
	if err = ebiten.Run(update, WIDTH, HEIGHT, 2, "Sprite"); err != nil {
		panic(err)
	}
}
