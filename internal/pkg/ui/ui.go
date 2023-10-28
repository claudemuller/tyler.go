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

// Layer

type layer struct {
	panels []panel
}

func (l *layer) update() {
	for _, p := range l.panels {
		p.update()
	}
}

func (l *layer) cascadeEvent(inputLocked *bool) {
	for _, p := range l.panels {
		p.cascadeEvent(inputLocked)
	}
}

func (l *layer) render() {
	for _, p := range l.panels {
		p.render()
	}
}

// Panel

type panel struct {
	label       string
	padding     int
	isDraggable *bool
	pos         *rl.Rectangle
	elements    []renderer
}

func newPanel(elements []renderer) panel {
	var draggable bool

	return panel{
		label:       "Options",
		padding:     10,
		isDraggable: &draggable,
		pos: &rl.Rectangle{
			X:      float32(rl.GetScreenWidth() - 200 - 10),
			Y:      10,
			Width:  200,
			Height: 200,
		},
		elements: elements,
	}
}

func (p *panel) update() {
	if *p.isDraggable {
		p.pos.X += rl.GetMouseDelta().X
		p.pos.Y += rl.GetMouseDelta().Y
	}

	for _, e := range p.elements {
		e.update()
	}
}

func (p *panel) cascadeEvent(inputLocked *bool) {
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), *p.pos) {
		*inputLocked = true
		panelHeader := rl.Rectangle{X: p.pos.X, Y: p.pos.Y, Width: p.pos.Width, Height: 30}

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), panelHeader) {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				print("mouse pressed in panel header\n")
				*p.isDraggable = true
			}
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				print("mouse released in panel header\n")
				*p.isDraggable = false
			}
			return
		}

		for _, e := range p.elements {
			e.cascadeEvent(p)
		}

		return
	}
	*inputLocked = false
}

func (p *panel) render() {
	gui.Panel(*p.pos, p.label)
	for i, e := range p.elements {
		e.render(p, float32(i)+0.8)
	}
}

// Element

type renderer interface {
	render(p *panel, extraYPadding float32)
	update()
	cascadeEvent(p *panel)
}

type element struct {
	pos      rl.Rectangle
	colour   rl.Color
	uiType   int
	editable bool
	value    int32
	label    string
	visible  bool
	trigger  func()
}

func (e *element) cascadeEvent(p *panel) {
	posInWorld := rl.Rectangle{
		X:      p.pos.X + e.pos.X + float32(p.padding),
		Y:      p.pos.Y + e.pos.Y + float32(p.padding),
		Width:  e.pos.Width,
		Height: e.pos.Height,
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), posInWorld) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			e.colour = rl.Red
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			e.colour = rl.Black
		}
	}
}

func (e *element) update() {
}

func (e *element) render(p *panel, row float32) {
	switch e.uiType {
	case uiTypeValueBox:
		if gui.ValueBox(
			rl.Rectangle{
				X:      float32(p.pos.ToInt32().X + e.pos.ToInt32().X + int32(p.padding)),
				Y:      float32(p.pos.Y) + float32(p.padding) + row*e.pos.Height,
				Width:  float32(p.pos.ToInt32().Width - int32(p.padding*2) - e.pos.ToInt32().X),
				Height: float32(e.pos.ToInt32().Height),
			},
			e.label,
			&e.value,
			0,
			100,
			e.editable,
		) {
			e.editable = true
		}

	case uiTypeButton:
		if gui.Button(
			rl.Rectangle{
				X:      float32(p.pos.ToInt32().X + e.pos.ToInt32().X + int32(p.padding)),
				Y:      float32(p.pos.Y) + float32(p.padding) + row*e.pos.Height,
				Width:  float32(p.pos.ToInt32().Width - int32(p.padding*2) - e.pos.ToInt32().X),
				Height: float32(e.pos.ToInt32().Height),
			},
			e.label,
		) {
			e.trigger()
		}

	case uiTypeInputBox:
		var value string

		if e.visible {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.8))
			var btnRes int32 = gui.TextInputBox(
				rl.Rectangle{
					X:      float32(rl.GetScreenWidth())/2 - 120,
					Y:      float32(rl.GetScreenHeight())/2 - 60,
					Width:  240,
					Height: 140,
				},
				"Open",
				e.label,
				"Ok;Cancel",
				&value,
				255,
				nil,
			)

			if btnRes == 1 {
				// if err := u.loadTilemap("jungle.png"); err != nil {
				// 	p.statusText = err.Error()
				// }
			}
			if btnRes == 0 || btnRes == 1 || btnRes == 2 {
				e.visible = false
				value = ""
			}
		}
	}
}
