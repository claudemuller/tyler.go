package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	IsRunning    bool

	rows     int32
	cols     int32
	cellSize int32

	frameTimer   float64
	frameTimeout float64
	pause        bool

	tilemap []bool
	dishBuf []bool
}

func NewGame(cellSize, cols, rows int32, frameSpeed float64) Game {
	return Game{
		ScreenWidth:  cellSize * cols,
		ScreenHeight: cellSize * rows,
		IsRunning:    true,
		rows:         rows,
		cols:         cols,
		cellSize:     cellSize,
		tilemap:      make([]bool, rows*cols),
		dishBuf:      make([]bool, rows*cols),
		frameTimer:   rl.GetTime(),
		frameTimeout: float64(frameSpeed),
		pause:        true,
	}
}

func (g *Game) Update() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.pause = !g.pause
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		g.Restart()
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		g.IsRunning = false
	}

	// Add cell.
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		g.tilemap[g.getMouseOverCell()] = true
	}
	// Kill cell.
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.tilemap[g.getMouseOverCell()] = false
	}

	// Control frame update.
	if rl.GetTime()-g.frameTimer < g.frameTimeout {
		return
	}
	g.frameTimer = rl.GetTime()
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for i := int32(0); i < g.rows*g.cols; i++ {
		x := int32(g.cellSize * (i % g.cols))
		y := int32(g.cellSize * (i / g.cols))
		if g.tilemap[i] {
			rl.DrawRectangle(x, y, g.cellSize, g.cellSize, rl.Black)
			continue
		}
		rl.DrawRectangleLines(x, y, g.cellSize, g.cellSize, rl.LightGray)
	}

	g.drawUI()

	rl.EndDrawing()
}

func (g *Game) Restart() {
	g.pause = true
	g.tilemap = make([]bool, g.rows*g.cols)
}

func (g *Game) drawUI() {
	bgColour := rl.White
	colour := rl.DarkBlue
	if g.pause {
		bgColour = rl.DarkBlue
		colour = rl.White
	}
	padding := 5
	spacing := 8
	fontSize := 20

	rl.DrawRectangle(10, 10, 204, 28, bgColour)
	s1 := spacing + padding
	rl.DrawText("<space> - edit mode", int32(10+padding), int32(s1), int32(fontSize), colour)
	rl.DrawText("<enter> - restart", int32(10+padding), int32(s1*3), int32(fontSize), rl.DarkBlue)
	rl.DrawText("<q> - quit", int32(10+padding), int32(s1*4+spacing), int32(fontSize), rl.DarkBlue)

	if g.pause {
		rl.DrawText("<s> - seed random cells", int32(10+padding), int32(s1*6+spacing-2), int32(fontSize), rl.DarkBlue)
		rl.DrawText("<left-mouse> - add cell", int32(10+padding), int32(s1*8+padding-3), int32(fontSize), rl.DarkBlue)
		rl.DrawText("<right-mouse> - kill cell", int32(10+padding), int32(s1*10+padding-6), int32(fontSize), rl.DarkBlue)

		activeCell := g.getMouseOverCell()
		if activeCell < 0 || activeCell >= g.rows*g.cols {
			return
		}

		liveText := "will die"
		liveColour := rl.Red
		rl.DrawText(liveText, rl.GetMouseX()-5, rl.GetMouseY()-10, int32(fontSize-10), liveColour)
	}
}

func (g *Game) getMouseOverCell() int32 {
	x := rl.GetMouseX() / g.cellSize
	y := rl.GetMouseY() / g.cellSize
	return y*g.cols + x
}
