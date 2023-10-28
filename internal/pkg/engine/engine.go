package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/claudemuller/tyler.go/internal/pkg/ui"
)

type Engine struct {
	WinWidth  int32
	WinHeight int32
	IsRunning bool
	camera    rl.Camera2D
	tileSize  float32
	cols      int
	rows      int
	pause     bool
	tilemap   []int
	ui        *ui.UI
}

func New(winWidth, winHeight int32) (Engine, error) {
	rl.InitWindow(winWidth, winHeight, "Tyler")

	rl.SetTargetFPS(60)

	gui, err := ui.New()
	if err != nil {
		return Engine{}, fmt.Errorf("error creating gui: %v", err)
	}

	e := Engine{
		WinWidth:  int32(winWidth),
		WinHeight: int32(winHeight),
		IsRunning: true,
		camera: rl.Camera2D{
			Target: rl.Vector2{
				X: 0, //64 * 50 / 2.0,
				Y: 0, //64 * 50 / 2.0,
			},
			Offset: rl.Vector2{
				X: 0, //winWidth / 2.0,
				Y: 0, //winHeight / 2.0,
			},
			Rotation: 0.0,
			Zoom:     1.0,
		},
		tileSize: 64,
		cols:     50,
		rows:     50,
		tilemap:  make([]int, 50*50),
		ui:       gui,
	}

	return e, nil
}

func (e *Engine) ProcessInput() {
	// ------------------------------------------------------------------------
	// UI Input

	e.ui.ProcessInput()
	if e.ui.InputLocked {
		return
	}

	// ------------------------------------------------------------------------
	// Mouse

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		e.tilemap[e.getMouseOverCell()] = 1
	}

	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		e.tilemap[e.getMouseOverCell()] = 0
	}

	mouseWheel := rl.GetMouseWheelMoveV()
	if mouseWheel.Y > 0 {
		e.camera.Zoom += 0.1
	}
	if mouseWheel.Y < 0 {
		e.camera.Zoom -= 0.1
	}

	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		e.camera.Offset.X += rl.GetMouseDelta().X
		e.camera.Offset.Y += rl.GetMouseDelta().Y
	}

	// ------------------------------------------------------------------------
	// Keyboard

	if rl.IsKeyPressed(rl.KeySpace) {
		e.pause = !e.pause
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		e.IsRunning = false
	}
}

func (e *Engine) Update() {
	e.ui.Update()
}

func (e *Engine) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	rl.BeginMode2D(e.camera)

	for i := 0; i < e.cols*e.rows; i++ {
		x := int32(int(e.tileSize) * (i % e.cols))
		y := int32(int(e.tileSize) * (i / e.cols))
		if e.tilemap[i] != 0 {
			rl.DrawRectangle(x, y, int32(e.tileSize), int32(e.tileSize), rl.Black)
			continue
		}
		rl.DrawRectangleLines(x, y, int32(e.tileSize), int32(e.tileSize), rl.LightGray)
	}

	rl.EndMode2D()

	e.ui.Render()

	rl.EndDrawing()
}

func (e *Engine) getMouseOverCell() int {
	// e.gui.UpdateDebugText("mousepos", fmt.Sprintf("x:%d y:%d", rl.GetMouseX(), rl.GetMouseY()))

	x := float32(rl.GetMouseX()-int32(e.camera.Offset.X)) / (e.tileSize * e.camera.Zoom)
	y := float32(rl.GetMouseY()-int32(e.camera.Offset.Y)) / (e.tileSize * e.camera.Zoom)
	return int(y)*e.cols + int(x)
}
