package ui

import "adagrad/internal/core"

type ControlState struct {
	Dragging         bool
	selectedX        int
	selectedZ        int
	LeftCtrlPressed  bool
	LeftShiftPressed bool
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
