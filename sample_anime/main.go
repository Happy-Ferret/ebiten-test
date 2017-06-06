package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
)

const (
	WIDTH  = 320
	HEIGHT = 240
)

var (
	tick int
	back *ebiten.Image
	fish *ebiten.Image
)

// uv map.
var imageRect = []image.Rectangle{
	image.Rect(0, 0, 64, 64),
	image.Rect(64, 0, 128, 64),
	image.Rect(128, 0, 192, 64),
	image.Rect(192, 0, 256, 64),
}

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	tick++

	screen.DrawImage(back, nil)

	w := 64
	h := 64

	// set up L2W.
	op := &ebiten.DrawImageOptions{}

	// anchoring to center.
	op.GeoM.Translate(-float64(w)*.5, -float64(h)*.5)

	// move.
	x := (WIDTH + (w >> 1)) - (tick % (WIDTH + w))
	y := HEIGHT >> 1
	op.GeoM.Translate(float64(x), float64(y))

	// set uv.
	op.SourceRect = &imageRect[(tick%20)/5]

	// queue the command.
	screen.DrawImage(fish, op)
	return nil
}

func main() {
	tick = 0

	var err error
	back, _, err = ebitenutil.NewImageFromFile("../assets/background.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	fish, _, err = ebitenutil.NewImageFromFile("../assets/animefishsheet000.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(update, WIDTH, HEIGHT, 2, "Sprite animation"); err != nil {
		panic(err)
	}
}
