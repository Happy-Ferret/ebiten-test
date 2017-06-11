// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi

package main

// https://github.com/hajimehoshi/ebiten/wiki/Tutorial%3AHandle-user-inputs

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

func update(screen *ebiten.Image) error {
	// Fill the screen with pink.
	screen.Fill(color.NRGBA{0xff, 0x00, 0xff, 0xff})

	// display the text using debug output.
	ebitenutil.DebugPrint(screen, "Mouse Click! Click! Click!")

	str := "\n"
	// WM_LBUTTONDOWN.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		str += "You're pressing the 'LEFT' mouse button.\n"
	}
	// WM_RBUTTONDOWN.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		str += "You're pressing the 'RIGHT' mouse button.\n"
	}
	// WM_CBUTTONDOWN.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		str += "You're pressing the 'MIDDLE' mouse button."
	}

	ebitenutil.DebugPrint(screen, str)
	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
