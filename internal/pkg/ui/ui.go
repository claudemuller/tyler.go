package ui

import (
	"fmt"
	"math"
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

type UI struct {
	Pos                rl.Rectangle
	IsActive           bool
	tileMapCols        int32
	tileMapRows        int32
	debugTexts         map[string]string
	showTilemapDiaglog bool
	tileMapFilename    string
	tilemap            rl.Texture2D
	style              int32
	statusText         string

	drag          bool
	styleDropOpen bool
}

func New() (*UI, error) {
	g := UI{
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

func (u *UI) ProcessInput() {
	if !rl.CheckCollisionPointRec(rl.GetMousePosition(), u.Pos) {
		u.IsActive = false
		return
	}
	u.IsActive = true

	panelHeader := rl.Rectangle{X: u.Pos.X, Y: u.Pos.Y, Width: u.Pos.Width, Height: 30}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		u.drag = rl.CheckCollisionPointRec(rl.GetMousePosition(), panelHeader)
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		u.drag = false
	}
}

func (u *UI) Update() {
	if u.drag {
		u.Pos.X += rl.GetMouseDelta().X
		u.Pos.Y += rl.GetMouseDelta().Y
	}
}

func (u *UI) Render() {
	u.drawOptions()
	u.drawDebug()
}

var editBox1 bool
var editBox2 bool

func (u *UI) drawOptions() {
	var elemHeight float32 = 35
	var elemPadding float32 = 40

	rg.Panel(u.Pos, "Options")

	if rg.ValueBox(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight, Width: 125, Height: 30},
		"cols",
		&u.tileMapCols,
		0,
		100,
		editBox1,
	) {
		editBox1 = true
	}

	if rg.ValueBox(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight*2, Width: 125, Height: 30},
		"rows",
		&u.tileMapRows,
		0,
		100,
		editBox2,
	) {
		editBox2 = true
	}

	if rg.Button(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight*3, Width: 125, Height: 30},
		rg.IconText(rg.ICON_FILE_SAVE, "Open Tilemap"),
	) {
		u.showTilemapDiaglog = true
	}

	oldStyle := u.style
	rg.SetStyle(rg.DROPDOWNBOX, rg.TEXT_ALIGNMENT, rg.TEXT_ALIGN_CENTER)
	if rg.DropdownBox(
		rl.Rectangle{X: u.Pos.X + elemPadding, Y: u.Pos.Y + elemHeight*4, Width: 125, Height: 30},
		strings.Join(styleNames, ";"),
		&u.style,
		u.styleDropOpen,
	) {
		u.styleDropOpen = !u.styleDropOpen
		if u.style != oldStyle {
			if err := loadStyle(styleFilenames[u.style]); err != nil {
				u.statusText = err.Error()
			}
		}
	}

	// rg.SliderBar(rl.Rectangle{600, 40, 120, 20}, "StartAngle", "", startAngle, -450, 450)
	// drawRing = rg.CheckBox(rl.Rectangle{600, 320, 20, 20}, "Draw Ring", drawRing)

	// TODO: check this
	if u.tilemap.ID > 0 {
		var tileSize int32 = 32

		numCols := u.tilemap.Width / tileSize
		numRows := u.tilemap.Height / tileSize

		spaceToDraw := float64(u.Pos.Width - 40)
		destMaxCols := int32(math.Floor(spaceToDraw / float64(tileSize)))

		var destRow int32 = 0
		var destCol int32 = 0

		for i := 0; i < int(numRows); i++ {
			for j := 0; j < int(numCols); j++ {
				rl.DrawTextureRec(
					u.tilemap,
					rl.Rectangle{
						X:      float32(j * int(tileSize)),
						Y:      float32(i * int(tileSize)),
						Width:  float32(tileSize),
						Height: float32(tileSize),
					},
					rl.Vector2{
						X: u.Pos.X + elemPadding + float32(tileSize*destCol),
						Y: u.Pos.Y + elemHeight*5 + float32(tileSize*destRow),
					},
					rl.White,
				)
				destCol++
				if destCol+1 > destMaxCols {
					destRow++
					destCol = 0
					continue
				}
			}
		}
	}

	rg.StatusBar(rl.Rectangle{X: 0, Y: float32(rl.GetScreenHeight()) - 20, Width: float32(rl.GetScreenWidth()), Height: 20}, u.statusText)

	if u.showTilemapDiaglog {
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.8))
		var btnRes int32 = rg.TextInputBox(
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

		fmt.Printf("%s\n", u.tileMapFilename)
		if btnRes == 1 {
			// TODO: Validate textInput value and save
			// strcpy(textInputFileName, textInput)
			if err := u.loadTilemap("jungle.png"); err != nil {
				u.statusText = err.Error()
			}
		}
		if btnRes == 0 || btnRes == 1 || btnRes == 2 {
			u.showTilemapDiaglog = false
			//strcpy(textInput, "\0");
			u.tileMapFilename = ""
		}
	}
}

func (u *UI) AddDebugText(key, text string) {
	u.debugTexts[key] = text
}

func (u *UI) UpdateDebugText(key, text string) {
	if _, ok := u.debugTexts[key]; !ok {
		fmt.Printf("no debug text with that name")
		return
	}
	u.debugTexts[key] = text
}

func (u *UI) loadTilemap(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("error loading theme file: %v\n", err)
	}
	image := rl.LoadImage(filename)
	u.tilemap = rl.LoadTextureFromImage(image)

	rl.UnloadImage(image)

	return nil
}

func (u *UI) drawDebug() {
	// rl.DrawFPS(int32(rl.GetScreenWidth())-90, 10)

	var i int
	for _, v := range u.debugTexts {
		rg.Label(
			rl.NewRectangle(
				float32(10),
				float32(10*i+10),
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
