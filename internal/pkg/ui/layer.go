package ui

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
