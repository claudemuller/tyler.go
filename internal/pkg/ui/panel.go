package ui

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
