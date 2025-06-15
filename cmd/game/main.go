package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/island-merge/pkg/core"
)

func main() {
	game := core.NewGame()
	
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Island Merge")
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}