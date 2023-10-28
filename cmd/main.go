package main

import (
	"fmt"

	"github.com/claudemuller/tyler.go/internal/pkg/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	game, err := engine.New(1024, 768)
	if err != nil {
		fmt.Printf("error starting engine: %v", err)
		return
	}

	for game.IsRunning && !rl.WindowShouldClose() {
		game.Update()
		game.Render()
	}

	rl.CloseWindow()
}
