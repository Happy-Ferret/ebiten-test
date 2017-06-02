// based on https://github.com/hajimehoshi/ebiten/blob/master/examples/paint/main.go

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

var (
	tick   int
	back	*ebiten.Image
	gopher *ebiten.Image
)

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	tick++

	screen.DrawImage(back, nil)

	w, h := gopher.Size()

	// set up L2W.
	op := &ebiten.DrawImageOptions{}

	// anchoring to center.
	op.GeoM.Translate(-float64(w)*.5, -float64(h)*.5)

	coef := float64(tick%ebiten.FPS) / ebiten.FPS

	// rotate.
	op.GeoM.Rotate(2.0 * math.Pi * coef)

	// move.
	x := (tick % (WIDTH + w)) - (w >> 1)
	y := (HEIGHT - (h >> 1)) - int(float64(h>>1)*(math.Sin(math.Pi*coef)))
	op.GeoM.Translate(float64(x), float64(y))

	// queue the command.
	screen.DrawImage(gopher, op)
	return nil
}

func main() {
	tick = 0

	var err error
	back, _, err = ebitenutil.NewImageFromFile("../assets/cal001.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	gopher, _, err = ebitenutil.NewImageFromFile("../assets/gophercolor.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(update, WIDTH, HEIGHT, 2, "Image"); err != nil {
		panic(err)
	}
}
