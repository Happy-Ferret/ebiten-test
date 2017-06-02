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
	WIDTH    = 320
	HEIGHT   = 240
	SCALE    = 2
	VELOCITY = 1.25
)

// Vector2D{int, int}
type Vector2Di struct {
	x, y int
}

func (v *Vector2Di) Set(x, y int) *Vector2Di {
	v.x = x
	v.y = y
	return v
}

func (v *Vector2Di) Sub(s *Vector2Di) *Vector2Di {
	v.x -= s.x
	v.y -= s.y
	return v
}

func (v *Vector2Di) IsZero() bool {
	return v.x == 0 && v.y == 0
}

// Vector2D{float, float}
type Vector2D struct {
	x, y float64
}

func (v *Vector2D) Clone() *Vector2D {
	return &Vector2D{v.x, v.y}
}

func (v *Vector2D) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector2D) Normalize() *Vector2D {
	l := v.Length()
	return v.Scale(1 / l)
}

func (v *Vector2D) Scale(s float64) *Vector2D {
	v.x *= s
	v.y *= s
	return v
}

func (v *Vector2D) Add(a *Vector2D) *Vector2D {
	v.x += a.x
	v.y += a.y
	return v
}

var (
	count       int
	lastPos     *Vector2Di
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
	if lastPos == nil {
		lastPos = &Vector2Di{x, y}
	} else {
		dir := &Vector2Di{x, y}
		if dir.Sub(lastPos).IsZero() {
			return
		}
		// line to.
		vel := &Vector2D{float64(dir.x), float64(dir.y)}
		length := vel.Length()
		vel.Normalize().Scale(VELOCITY)

		pos := &Vector2D{float64(lastPos.x), float64(lastPos.y)}
		for count := int(length / VELOCITY); count > 0; count-- {
			putDot(pos.x, pos.y)
			pos.Add(vel)
		}

		lastPos.Set(x, y)
	}

	putDot(float64(lastPos.x), float64(lastPos.y))
}

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		lineTo(mx, my)
	} else {
		lastPos = nil
	}

	// draw canvas to screen.
	screen.DrawImage(canvasImage, nil)

	return nil
}

func main() {
	// https://golang.org/src/image/image.go
	const a0, a1, a2 = 0x00, 0x40, 0x80
	brushImage, _ = ebiten.NewImageFromImage(&image.Alpha{
		Pix: []uint8{
			a0, a1, a1, a0,
			a1, a2, a2, a1,
			a1, a2, a2, a1,
			a0, a1, a1, a0,
		},
		Stride: 4, // 4 x 4 square.
		Rect:   image.Rect(0, 0, 4, 4),
	}, ebiten.FilterNearest)

	// clear canvasImage with Gray16(0xEEEF).
	canvasImage, _ = ebiten.NewImage(WIDTH, HEIGHT, ebiten.FilterNearest)
	canvasImage.Fill(color.Gray16{0xeeef})

	lastPos = nil
	if err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "Ebiten Paint"); err != nil {
		log.Fatal(err)
	}
}
