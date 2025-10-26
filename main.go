package main

import (
	"log"

	// "adagrad/internal/ui"
	"adagrad/internal/core"
	"adagrad/internal/game"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	// if err := ui.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	//
	log.Printf("Starting Game")
	game := game.NewGame(core.Width, core.Height)
	log.Printf("Game")
	defer glfw.Terminate()
	game.Run()
}
