package main

import (
	"fmt"

	e "github.com/claudemuller/tyler.go/internal/pkg/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	engine, err := e.New(1024, 768)
	if err != nil {
		fmt.Printf("error starting engine: %v", err)
		return
	}

	for engine.IsRunning && !rl.WindowShouldClose() {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}

	rl.CloseWindow()
}
