package main

// https://github.com/hajimehoshi/ebiten/wiki/Tutorial:Your-first-game-in-Ebiten

import (
    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
    ebitenutil.DebugPrint(screen, "Our first game in Ebiten!")
    return nil
}

func main() {
    if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
        panic(err)
    }
}
