package game

import (
	"math/rand"
	"time"

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

	dish    []bool
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
		dish:         make([]bool, rows*cols),
		dishBuf:      make([]bool, rows*cols),
		frameTimer:   rl.GetTime(),
		frameTimeout: float64(frameSpeed),
		pause:        true,
	}
}

func (g *Game) Seed() {
	// Random cell fill.
	rand.Seed(time.Now().UnixNano())
	for i := int32(0); i < g.rows*g.cols; i++ {
		if rand.Intn(7) == 0 {
			g.dish[i] = true
		}
	}
}

func (g *Game) Update() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.pause = !g.pause
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		g.Restart()
	}

	if rl.IsKeyPressed(rl.KeyS) {
		g.Seed()
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		g.IsRunning = false
	}

	// Add cell.
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		g.dish[g.getMouseOverCell()] = true
	}
	// Kill cell.
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.dish[g.getMouseOverCell()] = false
	}

	// Control frame update.
	if rl.GetTime()-g.frameTimer < g.frameTimeout {
		return
	}
	g.frameTimer = rl.GetTime()

	if !g.pause {
		copy(g.dishBuf, g.dish)

		for i := int32(0); i < g.rows*g.cols; i++ {
			shouldLive, _ := g.shouldLive(i)
			g.dishBuf[i] = shouldLive
		}

		copy(g.dish, g.dishBuf)
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for i := int32(0); i < g.rows*g.cols; i++ {
		x := int32(g.cellSize * (i % g.cols))
		y := int32(g.cellSize * (i / g.cols))
		if g.dish[i] {
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
	g.dish = make([]bool, g.rows*g.cols)
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
		shouldLive, _ := g.shouldLive(activeCell)
		if shouldLive {
			liveText = "will live"
			liveColour = rl.Green
		}
		rl.DrawText(liveText, rl.GetMouseX()-5, rl.GetMouseY()-10, int32(fontSize-10), liveColour)
	}
}

func (g *Game) shouldLive(i int32) (bool, int) {
	var neighbours int

	dishMin := int32(0)
	dishMax := g.rows * g.cols

	// Left
	if i-1 >= dishMin {
		neighbours += btoi(g.dish[i-1])
	}
	// Right
	if i+1 < dishMax {
		neighbours += btoi(g.dish[i+1])
	}

	// Top
	if i-g.cols >= dishMin {
		neighbours += btoi(g.dish[i-g.cols])
	}
	// Top left
	if i-1-g.cols >= dishMin {
		neighbours += btoi(g.dish[i-1-g.cols])
	}
	// Top right
	if i+1-g.cols >= dishMin {
		neighbours += btoi(g.dish[i+1-g.cols])
	}

	// Bottom
	if i+g.cols < dishMax {
		neighbours += btoi(g.dish[i+g.cols])
	}
	// Bottom left
	if i-1+g.cols < dishMax {
		neighbours += btoi(g.dish[i-1+g.cols])
	}
	// Bottom right
	if i+1+g.cols < dishMax {
		neighbours += btoi(g.dish[i+1+g.cols])
	}

	if !g.dish[i] {
		return neighbours == 3, neighbours
	}

	return neighbours >= 2 && neighbours <= 3, neighbours
}

func (g *Game) getMouseOverCell() int32 {
	x := rl.GetMouseX() / g.cellSize
	y := rl.GetMouseY() / g.cellSize
	return y*g.cols + x
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
