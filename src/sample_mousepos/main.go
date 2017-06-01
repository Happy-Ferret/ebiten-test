package main

// https://github.com/hajimehoshi/ebiten/wiki/Tutorial%3AHandle-user-inputs

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"fmt"
)

func update(screen *ebiten.Image) error {
	// Fill the screen with pink.
	screen.Fill(color.NRGBA{0xff, 0x00, 0xff, 0xff})

	// display the text using debug output.
	ebitenutil.DebugPrint(screen, "Let's your mouse move around!")

	// ebiten.CursorPosition will give us the mouse cursor position in "int".
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n X: %d, Y: %d", x, y))

	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
