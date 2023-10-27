package main

import (
	"github.com/claudemuller/tyler.go/internal/pkg/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	game := engine.New(1024, 768)

	rl.InitWindow(game.WinWidth, game.WinHeight, "Tyler")

	rl.SetTargetFPS(60)

	for game.IsRunning && !rl.WindowShouldClose() {
		game.Update()
		game.Render()
	}

	rl.CloseWindow()
}
