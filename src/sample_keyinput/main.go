package main

// https://github.com/hajimehoshi/ebiten/wiki/Tutorial%3AHandle-user-inputs

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

func update(screen *ebiten.Image) error {
	// Fill the screen with #0000FF color
	screen.Fill(color.NRGBA{0x00, 0x00, 0xff, 0xff})

	// Display the text though the debug function
	ebitenutil.DebugPrint(screen, "Press cursor keys!")

	str := "\n"
	// When the "up arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		str += "You're pressing the 'UP' button.\n"
	}
	// When the "down arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		str += "You're pressing the 'DOWN' button.\n"
	}
	// When the "left arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		str += "You're pressing the 'LEFT' button.\n"
	}
	// When the "right arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		str += "You're pressing the 'RIGHT' button."
	}

	ebitenutil.DebugPrint(screen, str)
	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
