package ui

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	statusBarHeight = 20
)

const (
	mousePressed = iota
	mouseReleased
)

const (
	uiTypeValueBox = iota + 1
	uiTypeButton
	uiTypeInputBox
)

type UI struct {
	InputLocked bool
	layers      []layer
	statusText  string
}

func New() (*UI, error) {
	gui.LoadStyleDefault()

	colsInput := element{
		pos: rl.Rectangle{
			X:      40,
			Y:      0,
			Width:  0,
			Height: 30,
		},
		label:  "cols",
		uiType: uiTypeValueBox,
	}

	rowsInput := element{
		pos: rl.Rectangle{
			X:      40,
			Y:      0,
			Width:  0,
			Height: 30,
		},
		label:  "rows",
		uiType: uiTypeValueBox,
	}

	fileDialog := element{
		pos: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  0,
			Height: 30,
		},
		label:  gui.IconText(gui.ICON_FILE_SAVE, "Open file..."),
		uiType: uiTypeInputBox,
	}

	button := element{
		pos: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  0,
			Height: 30,
		},
		label:  gui.IconText(gui.ICON_FILE_SAVE, "Open Tilemap"),
		uiType: uiTypeButton,
		trigger: func() {
			fileDialog.visible = true
		},
	}

	ui := UI{
		layers: []layer{
			{
				[]panel{
					newPanel([]renderer{&colsInput, &rowsInput, &button, &fileDialog}),
				},
			},
		},
		statusText: "Ready...",
	}

	return &ui, nil
}

func (u *UI) ProcessInput() {
	for _, l := range u.layers {
		l.cascadeEvent(&u.InputLocked)
	}
}

func (u *UI) Update() {
	for _, l := range u.layers {
		l.update()
	}
}

func (u *UI) Render() {
	for _, l := range u.layers {
		l.render()
	}

	gui.StatusBar(
		rl.Rectangle{
			X:      0,
			Y:      float32(rl.GetScreenHeight()) - statusBarHeight,
			Width:  float32(rl.GetScreenWidth()),
			Height: statusBarHeight,
		},
		u.statusText,
	)
}
