// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi
// based on https://github.com/hajimehoshi/ebiten/blob/master/examples/paint/main.go

package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
)

const (
	WIDTH  = 320
	HEIGHT = 240
)

var (
	tick  int
	uvtex *ebiten.Image
)

// uv map.
var imageRect = []image.Rectangle{
	image.Rect(0, 0, 32, 32),
	image.Rect(32, 0, 64, 32),
	image.Rect(0, 32, 32, 64),
	image.Rect(32, 32, 64, 64),
}

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	tick++

	w := 32
	h := 32

	// set up L2W.
	op := &ebiten.DrawImageOptions{}

	// anchoring to center.
	op.GeoM.Translate(-float64(w)*.5, -float64(h)*.5)
	op.GeoM.Scale(2.0, 2.0)

	// move.
	x := WIDTH >> 1
	y := HEIGHT >> 1
	op.GeoM.Translate(float64(x), float64(y))

	// set uv.
	op.SourceRect = &imageRect[(tick%160)/40]

	// queue the command.
	screen.DrawImage(uvtex, op)

	// display UV.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("UV: %d:%d:%d:%d",
		op.SourceRect.Min.X, op.SourceRect.Min.Y,
		op.SourceRect.Max.X, op.SourceRect.Max.Y))

	return nil
}

func main() {
	tick = 0

	var err error
	uvtex, _, err = ebitenutil.NewImageFromFile("../assets/uvtest.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(update, WIDTH, HEIGHT, 2, "Sprite animation"); err != nil {
		panic(err)
	}
}
