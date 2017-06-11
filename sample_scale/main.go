// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi

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

	num := 24

	for i := 0; i < num; i++ {
		op.GeoM.Reset()
		op.ColorM.Reset()

		// anchoring to center.
		op.GeoM.Translate(-float64(w)*.5, -float64(h)*.5)

		// scale.
		scale := float64(i+1) / float64(num)
		op.GeoM.Scale(scale, scale)

		// move.
		total := WIDTH + w
		x := tick + i*(total/num)
		x = (x % total) - (w >> 1)

		theta := math.Pi * float64((tick+i*(ebiten.FPS/num))%ebiten.FPS) / ebiten.FPS
		sin := math.Sin(theta)
		y := (HEIGHT - (h >> 1)) - int(float64(h>>1)*sin)
		op.GeoM.Translate(float64(x), float64(y))

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
	gopher, _, err = ebitenutil.NewImageFromFile("../assets/gophercolor.png", ebiten.FilterLinear)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(update, WIDTH, HEIGHT, 2, "Scale"); err != nil {
		panic(err)
	}
}
