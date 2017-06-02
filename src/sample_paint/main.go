// based on https://github.com/hajimehoshi/ebiten/blob/master/examples/paint/main.go

package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	WIDTH  = 320
	HEIGHT = 240
	SCALE = 2
	VELOCITY = 1.5
)

var (
	count, lx, ly	int
	brushImage  *ebiten.Image
	canvasImage *ebiten.Image
)

func putDot(x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.ColorM.Scale(1.0, 0.50, 0.125, 1.0)
	theta := 2.0 * math.Pi * float64(count%ebiten.FPS) / ebiten.FPS
	op.ColorM.RotateHue(theta)

	canvasImage.DrawImage(brushImage, op)
	count++
}

func lineTo(x, y int) {
	vx := x
	vy := y

	if lx < 0 || ly < 0 {
		vx = 0
		vy = 0
	} else {
		vx -= lx
		vy -= ly
		if vx == 0 || vy == 0 {
			return
		}
	}

	// line to.
	fx := float64(vx);
	fy := float64(vy);
	length := math.Sqrt(fx * fx + fy * fy);
	fx = fx / length * VELOCITY;
	fy = fy / length * VELOCITY;
	sx := float64(lx)
	sy := float64(ly)
	for count := int(length / VELOCITY); count > 0; count-- {
		putDot(sx, sy)
		sx += fx
		sy += fy
	}
	putDot(float64(x), float64(y))
	lx = x
	ly = y
}

func update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		lineTo(mx, my)
	} else {
		lx = -1
		ly = -1
	}

	// draw canvas to screen.
	screen.DrawImage(canvasImage, nil)

	return nil
}

func main() {
	// https://golang.org/src/image/image.go
	const a0, a1, a2 = 0x00, 0x40, 0x80
	brushImage, _ = ebiten.NewImageFromImage(&image.Alpha{
		Pix: []uint8 {
			a0, a1, a1, a0,
			a1, a2, a2, a1,
			a1, a2, a2, a1,
			a0, a1, a1, a0,
		},
		Stride: 4,	// 4 x 4 square.
		Rect: image.Rect(0, 0, 4, 4),
	}, ebiten.FilterNearest)

	// clear canvasImage with Gray16(0xEEEF).
	canvasImage, _ = ebiten.NewImage(WIDTH, HEIGHT, ebiten.FilterNearest)
	canvasImage.Fill(color.Gray16{0xeeef})

	lx = -1
	ly = -1
	if err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "Ebiten Paint"); err != nil {
		log.Fatal(err)
	}
}
