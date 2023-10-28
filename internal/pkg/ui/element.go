package ui

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
