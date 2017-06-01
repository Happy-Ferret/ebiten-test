package main

// https://github.com/hajimehoshi/ebiten/wiki/Tutorial:Your-first-game-in-Ebiten

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
	str := "Our first game in Ebiten!\n"
	str += "https://github.com/hajimehoshi/ebiten\n"
	ebitenutil.DebugPrint(screen, str)
	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
