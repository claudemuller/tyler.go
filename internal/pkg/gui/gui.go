package gui

import (
	"fmt"
	"os"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	theme_blueish = "assets/styles/bluish/bluish.rgs"
)

type Gui struct {
	Pos               rl.Rectangle
	tileMapCols       int32
	tileMapRows       int32
	debugTexts        map[string]string
	showTilemapPicker bool
	tileMapFilename   string
	statusText        string
}

func New() (*Gui, error) {
	if _, err := os.Stat(theme_blueish); err != nil {
		return nil, fmt.Errorf("error loading theme file: %v\n", err)
	}
	rg.LoadStyle(theme_blueish)

	g := Gui{
		Pos: rl.Rectangle{
			X:      float32(rl.GetScreenWidth()-rl.GetScreenWidth()/4) - 20,
			Y:      20,
			Width:  float32(rl.GetScreenWidth() / 4),
			Height: 500,
		},
		tileMapCols: 50,
		tileMapRows: 50,
		debugTexts:  make(map[string]string),
		statusText:  "Ready...",
	}

	return &g, nil
}

func (g *Gui) Update() {
	panelHeader := rl.Rectangle{X: g.Pos.X, Y: g.Pos.Y, Width: g.Pos.Width, Height: 30}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		if rl.GetMouseX() < int32(g.Pos.X) {
			rl.SetMousePosition(int(g.Pos.X), int(rl.GetMouseY()))
		}
		if rl.GetMouseX() > int32(g.Pos.X+g.Pos.Width) {
			rl.SetMousePosition(int(g.Pos.X+g.Pos.Width), int(rl.GetMouseY()))
		}
		if rl.GetMouseY() < int32(g.Pos.Y) {
			rl.SetMousePosition(int(rl.GetMouseX()), int(g.Pos.Y))
		}
		if rl.GetMouseY() > int32(g.Pos.Y+g.Pos.Height) {
			rl.SetMousePosition(int(rl.GetMouseX()), int(g.Pos.Y+g.Pos.Height))
		}

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), panelHeader) {
			g.Pos.X += rl.GetMouseDelta().X
			g.Pos.Y += rl.GetMouseDelta().Y
		}
	}
}

func (g *Gui) Render() {
	g.drawOptions()
	g.drawDebug()
}

func (g *Gui) drawOptions() {
	var elemHeight float32 = 35
	var elemPadding float32 = 40

	rg.Panel(g.Pos, "Options")

	rg.ValueBox(
		rl.Rectangle{X: g.Pos.X + elemPadding, Y: g.Pos.Y + elemHeight, Width: 125, Height: 30},
		"cols",
		&g.tileMapCols,
		0,
		100,
		true,
	)
	rg.ValueBox(
		rl.Rectangle{X: g.Pos.X + elemPadding, Y: g.Pos.Y + elemHeight*2, Width: 125, Height: 30},
		"rows",
		&g.tileMapRows,
		0,
		100,
		true,
	)

	if rg.Button(
		rl.Rectangle{X: g.Pos.X + elemPadding, Y: g.Pos.Y + elemHeight*3, Width: 125, Height: 30},
		rg.IconText(rg.ICON_FILE_SAVE, "Open Tilemap"),
	) {
		g.showTilemapPicker = true
	}

	if g.showTilemapPicker {
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.RayWhite, 0.8))
		var result int32 = rg.TextInputBox(
			rl.Rectangle{
				X:      float32(rl.GetScreenWidth())/2 - 120,
				Y:      float32(rl.GetScreenHeight())/2 - 60,
				Width:  240,
				Height: 140,
			},
			"Open",
			rg.IconText(rg.ICON_FILE_SAVE, "Open file..."),
			"Ok;Cancel",
			&g.tileMapFilename,
			255,
			nil,
		)

		if result == 1 {
			// TODO: Validate textInput value and save
			// strcpy(textInputFileName, textInput)
			if _, err := os.Stat(g.tileMapFilename); err != nil {
				g.statusText = fmt.Sprintf("tilemap file not found: %v", err)
			}
		}
		if (result == 0) || (result == 1) || (result == 2) {
			g.showTilemapPicker = false
			//strcpy(textInput, "\0");
			g.tileMapFilename = ""
		}
	}

	rg.StatusBar(rl.Rectangle{X: 0, Y: float32(rl.GetScreenHeight()) - 20, Width: float32(rl.GetScreenWidth()), Height: 20}, g.statusText)

	// comboText := []string{
	// 	"ashes", "bluish", "candy", "cherry", "cyber", "dark",
	// 	"default", "enefete", "jungle", "lavanda", "sunny", "terminal",
	// }
	// comboList := strings.Join(comboText, ";")
	// comboActive = rg.ComboBox(rl.NewRectangle(500, 280, 200, 20), comboList, comboActive)

	// if comboLastActive != comboActive {
	// 	ch := comboText[comboActive]
	// 	filename := fmt.Sprintf("assets/styles/%s/%s.rgs", ch, ch)
	// 	rg.LoadStyle(filename)
	// 	comboLastActive = comboActive
	// }

	// rg.SliderBar(rl.Rectangle{600, 40, 120, 20}, "StartAngle", "", startAngle, -450, 450)
	// drawRing = rg.CheckBox(rl.Rectangle{600, 320, 20, 20}, "Draw Ring", drawRing)

}

func (g *Gui) AddDebugText(key, text string) {
	g.debugTexts[key] = text
}

func (g *Gui) UpdateDebugText(key, text string) {
	if _, ok := g.debugTexts[key]; !ok {
		fmt.Printf("no debug text with that name")
		return
	}
	g.debugTexts[key] = text
}

func (g *Gui) drawDebug() {
	// rl.DrawFPS(int32(rl.GetScreenWidth())-90, 10)

	bgColour := rl.White
	colour := rl.DarkBlue
	padding := 10
	spacing := 10
	fontSize := 20

	rl.DrawRectangle(10, 10, 204, int32(len(g.debugTexts)*20), bgColour)
	s1 := spacing + padding

	var i int
	for _, v := range g.debugTexts {
		rl.DrawText(
			v,
			int32(10+padding),
			int32(s1*i+padding),
			int32(fontSize),
			colour,
		)
		i++
	}
}
