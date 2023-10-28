package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
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
}

var debugTexts map[string]string = make(map[string]string)

func addDebugText(key, text string) {
	if _, ok := debugTexts[key]; !ok {
		// debugTexts[]
	}
	debugTexts[key] = text
}

func New(winWidth, winHeight float32) Engine {
	return Engine{
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
	}
}

func (e *Engine) Update() {
	if rl.IsKeyPressed(rl.KeySpace) {
		e.pause = !e.pause
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		e.IsRunning = false
	}

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

	if rl.IsKeyPressed(rl.KeyA) {
		e.camera.Offset.X += 10
	}
	if rl.IsKeyPressed(rl.KeyO) {
		e.camera.Offset.X -= 10
	}
	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		e.camera.Offset.X += rl.GetMouseDelta().X
		e.camera.Offset.Y += rl.GetMouseDelta().Y
	}
	// mousePos := rl.GetMousePosition()
	// if mousePos.X > float32(e.WinWidth-20) && mousePos.X < float32(e.WinWidth) {
	// 	e.camera.Offset.X += 10
	// }
	// if mousePos.X < float32(20) && mousePos.X > 0 {
	// 	e.camera.Offset.X -= 10
	// }

	// Control frame update.
	// if rl.GetTime()-g.frameTimer < g.frameTimeout {
	// 	return
	// }
	// g.frameTimer = rl.GetTime()
}

func (e *Engine) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

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

	e.drawUI()

	rl.EndDrawing()
}

func (e *Engine) drawUI() {
	bgColour := rl.White
	colour := rl.DarkBlue
	if e.pause {
		bgColour = rl.DarkBlue
		colour = rl.White
	}
	padding := 10
	spacing := 10
	fontSize := 20

	rl.DrawRectangle(10, 10, 204, int32(len(debugTexts)*20), bgColour)
	s1 := spacing + padding

	var i int
	for _, v := range debugTexts {
		rl.DrawText(
			v,
			int32(10+padding),
			int32(s1*i+padding),
			int32(fontSize),
			colour,
		)
		i++
	}

	addDebugText("test", "testing")

	if e.pause {
		rl.DrawText("<s> - seed random cells", int32(10+padding), int32(s1*6+spacing-2), int32(fontSize), rl.DarkBlue)
		rl.DrawText("<left-mouse> - add cell", int32(10+padding), int32(s1*8+padding-3), int32(fontSize), rl.DarkBlue)
		rl.DrawText("<right-mouse> - kill cell", int32(10+padding), int32(s1*10+padding-6), int32(fontSize), rl.DarkBlue)

		activeCell := e.getMouseOverCell()
		if activeCell < 0 || activeCell >= e.rows*e.cols {
			return
		}
	}
}

func (e *Engine) getMouseOverCell() int {
	addDebugText("mousepos", fmt.Sprintf("x:%d y:%d", rl.GetMouseX(), rl.GetMouseY()))
	x := float32(rl.GetMouseX()-int32(e.camera.Offset.X)) / (e.tileSize * e.camera.Zoom)
	y := float32(rl.GetMouseY()-int32(e.camera.Offset.Y)) / (e.tileSize * e.camera.Zoom)
	return int(y)*e.cols + int(x)
}
