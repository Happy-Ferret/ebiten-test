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
	back   *ebiten.Image
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

	palette := [][]float64{
		{0.8, 0.8, 0.8}, // 15
		{0.5, 0.5, 0.0}, // 14
		{0.0, 0.5, 0.5}, // 13
		{0.0, 0.5, 0.0}, // 12
		{0.5, 0.0, 0.5}, // 11
		{0.5, 0.0, 0.0}, // 10
		{0.0, 0.0, 0.5}, // 9
		{0.5, 0.5, 0.5}, // 8
		{0.0, 0.0, 0.0}, // 0
		{1.0, 1.0, 0.0}, // 6
		{0.0, 1.0, 1.0}, // 5
		{0.0, 1.0, 0.0}, // 4
		{1.0, 0.0, 1.0}, // 3
		{1.0, 0.0, 0.0}, // 2
		{0.0, 0.0, 1.0}, // 1
		{1.0, 1.0, 1.0}, // 7
	}
	num := len(palette)

	for i := 0; i < num; i++ {
		op.GeoM.Reset()
		op.ColorM.Reset()

		// anchoring to center.
		op.GeoM.Translate(-float64(w)*.5, -float64(h)*.5)

		// move.
		total := WIDTH + w
		x := tick + i*(total/num)
		x = (x % total) - (w >> 1)

		theta := math.Pi * float64((tick+i*(ebiten.FPS/num))%ebiten.FPS) / ebiten.FPS
		sin := math.Sin(theta)
		y := (HEIGHT - (h >> 1)) - int(float64(h>>1)*sin)
		op.GeoM.Translate(float64(x), float64(y))

		// color.
		op.ColorM.Scale(palette[i][0], palette[i][1], palette[i][2], 1.0)

		// queue the command.
		screen.DrawImage(gopher, op)
	}
	return nil
}

func main() {
	tick = 0

	var err error
	back, _, err = ebitenutil.NewImageFromFile("../assets/cal003.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	gopher, _, err = ebitenutil.NewImageFromFile("../assets/gopherbw.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(update, WIDTH, HEIGHT, 2, "Color"); err != nil {
		panic(err)
	}
}
