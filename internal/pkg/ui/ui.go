package ui

import (
	"fmt"
	"os"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var styleNames = []string{
	"ashes", "bluish", "candy", "cherry", "cyber", "dark",
	"enefete", "jungle", "lavanda", "sunny", "terminal",
}
var styleFilenames = []string{
	"assets/styles/ashes/ashes.rgs",
	"assets/styles/bluish/bluish.rgs",
	"assets/styles/candy/candy.rgs",
	"assets/styles/cherry/cherry.rgs",
	"assets/styles/cyber/cyber.rgs",
	"assets/styles/dark/dark.rgs",
	"assets/styles/enefete/enefete.rgs",
	"assets/styles/jungle/jungle.rgs",
	"assets/styles/lavanda/lavanda.rgs",
	"assets/styles/sunny/sunny.rgs",
	"assets/styles/terminal/terminal.rgs",
}

type Ui struct {
	Pos               rl.Rectangle
	tileMapCols       int32
	tileMapRows       int32
	debugTexts        map[string]string
	showTilemapPicker bool
	tileMapFilename   string
	style             int32
	statusText        string

	styleDropOpen bool
}

func New() (*Ui, error) {
	g := Ui{
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

	rg.LoadStyleDefault()

	return &g, nil
}

func (u *Ui) Update() {
	panelHeader := rl.Rectangle{X: u.Pos.X, Y: u.Pos.Y, Width: u.Pos.Width, Height: 30}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		if rl.GetMouseX() < int32(u.Pos.X) {
			rl.SetMousePosition(int(u.Pos.X), int(rl.GetMouseY()))
		}
		if rl.GetMouseX() > int32(u.Pos.X+u.Pos.Width) {
			rl.SetMousePosition(int(u.Pos.X+u.Pos.Width), int(rl.GetMouseY()))
		}
		if rl.GetMouseY() < int32(u.Pos.Y) {
			rl.SetMousePosition(int(rl.GetMouseX()), int(u.Pos.Y))
		}
		if rl.GetMouseY() > int32(u.Pos.Y+u.Pos.Height) {
			rl.SetMousePosition(int(rl.GetMouseX()), int(u.Pos.Y+u.Pos.Height))
		}

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), panelHeader) {
			u.Pos.X += rl.GetMouseDelta().X
			u.Pos.Y += rl.GetMouseDelta().Y
		}
	}
}

func (u *Ui) Render() {
	u.drawOptions()
	u.drawDebug()
}

func (u *Ui) drawOptions() {
	var elemHeight float32 = 35
	var elemPadding float32 = 40

	rg.Panel(u.Pos, "Options")

	rg.ValueBox(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight, Width: 125, Height: 30},
		"cols",
		&u.tileMapCols,
		0,
		100,
		true,
	)
	rg.ValueBox(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight*2, Width: 125, Height: 30},
		"rows",
		&u.tileMapRows,
		0,
		100,
		true,
	)

	if rg.Button(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight*3, Width: 125, Height: 30},
		rg.IconText(rg.ICON_FILE_SAVE, "Open Tilemap"),
	) {
		u.showTilemapPicker = true
	}

	if u.showTilemapPicker {
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.8))
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
			&u.tileMapFilename,
			255,
			nil,
		)

		if result == 1 {
			// TODO: Validate textInput value and save
			// strcpy(textInputFileName, textInput)
			if _, err := os.Stat(u.tileMapFilename); err != nil {
				u.statusText = fmt.Sprintf("tilemap file not found: %v", err)
			}
		}
		if (result == 0) || (result == 1) || (result == 2) {
			u.showTilemapPicker = false
			//strcpy(textInput, "\0");
			u.tileMapFilename = ""
		}
	}

	rg.StatusBar(rl.Rectangle{X: 0, Y: float32(rl.GetScreenHeight()) - 20, Width: float32(rl.GetScreenWidth()), Height: 20}, u.statusText)

	oldStyle := u.style
	rg.SetStyle(rg.DROPDOWNBOX, rg.TEXT_ALIGNMENT, rg.TEXT_ALIGN_CENTER)
	if rg.DropdownBox(
		rl.Rectangle{u.Pos.X + elemPadding, u.Pos.Y + elemHeight*4, 125, 30},
		strings.Join(styleNames, ";"),
		&u.style,
		u.styleDropOpen,
	) {
		u.styleDropOpen = !u.styleDropOpen
		if u.style != oldStyle {
			loadStyle(styleFilenames[u.style])
		}
	}

	// rg.SliderBar(rl.Rectangle{600, 40, 120, 20}, "StartAngle", "", startAngle, -450, 450)
	// drawRing = rg.CheckBox(rl.Rectangle{600, 320, 20, 20}, "Draw Ring", drawRing)

}

func (u *Ui) AddDebugText(key, text string) {
	u.debugTexts[key] = text
}

func (u *Ui) UpdateDebugText(key, text string) {
	if _, ok := u.debugTexts[key]; !ok {
		fmt.Printf("no debug text with that name")
		return
	}
	u.debugTexts[key] = text
}

func (u *Ui) drawDebug() {
	// rl.DrawFPS(int32(rl.GetScreenWidth())-90, 10)

	var i int
	for _, v := range u.debugTexts {
		rg.Label(
			rl.NewRectangle(
				float32(10),
				float32(10),
				100,
				20,
			),
			v,
		)
		i++
	}
}

func loadStyle(style string) error {
	if _, err := os.Stat(style); err != nil {
		return fmt.Errorf("error loading theme file: %v\n", err)
	}
	rg.LoadStyle(style)
	return nil
}
