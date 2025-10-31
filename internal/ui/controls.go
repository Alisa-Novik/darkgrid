package ui

import "adagrad/internal/core"

type ControlState struct {
	SelectedTile       core.Tile
	RectStart, RectEnd core.Tile
}

func (ct *ControlState) BeginRect(tile core.Tile) {
	ct.RectStart = tile
}

func (ct *ControlState) EndRect(tile core.Tile) {
	ct.RectEnd = tile
}

func (ct *ControlState) UpdateRect(tile core.Tile) {
	ct.RectEnd = tile
}

func (ct *ControlState) IsInRect(tileX, tileZ int) bool {
	return tileX >= ct.RectStart.X && tileX <= ct.RectEnd.X &&
		tileZ >= ct.RectStart.Z && tileZ <= ct.RectEnd.Z
}
