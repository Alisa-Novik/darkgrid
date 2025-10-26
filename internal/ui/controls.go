package ui

import "adagrad/internal/core"

type ControlState struct {
	selectedX int
	selectedZ int
}

func (ct *ControlState) SelectTile(x, z int) {
	if !core.InBounds(x, z) {
		ct.selectedX = -1
		ct.selectedZ = -1
		return
	}
	ct.selectedX = x
	ct.selectedZ = z
}

func (ct *ControlState) SelectedTile() (int, int) {
	return ct.selectedX, ct.selectedZ
}
