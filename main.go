package main

import (
	"log"

	// "adagrad/internal/ui"
	"adagrad/internal/core"
	"adagrad/internal/game"
)

func main() {
	// if err := ui.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	//
	log.Fatal("Starting Game")
	game := game.NewGame(core.Width, core.Height)
	log.Fatal("Game")
	game.Run()
}
