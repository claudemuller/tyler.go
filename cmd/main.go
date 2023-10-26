package main

import (
	g "github.com/claudemuller/tyler.go/internal/pkg/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	cellSize := int32(20)
	rows := int32(50)
	cols := int32(50)

	game := g.NewGame(cellSize, rows, cols, 0.5)

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Conway's Game of Life")

	rl.SetTargetFPS(60)

	for game.IsRunning && !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}

	rl.CloseWindow()
}
