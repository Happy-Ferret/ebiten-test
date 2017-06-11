// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi
// https://github.com/hajimehoshi/ebiten/wiki/Tutorial:Screen,-colors-and-squares

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var square *ebiten.Image

func update(screen *ebiten.Image) error {
	// Fill the screen with #FF0000 color
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})

	// Display the text though the debug function
	str := "Our first game in Ebiten!\n"
	str += "https://github.com/hajimehoshi/ebiten\n"
	str += "Fill color and put a square block.\n"
	ebitenutil.DebugPrint(screen, str)

	if square == nil {
		// Create an 16x16 image
		square, _ = ebiten.NewImage(16, 16, ebiten.FilterNearest)
	}

	// Fill the square with the white color
	square.Fill(color.White)

	// Create an empty option struct
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(64, 64)

	// Draw the square image to the screen with an empty option
	screen.DrawImage(square, opts)

	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
