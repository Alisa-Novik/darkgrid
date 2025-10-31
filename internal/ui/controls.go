package ui

import "adagrad/internal/core"

type ControlState struct {
	SelectedTile       core.Tile
	RectStart, RectEnd core.Tile
}

func NewControlState() *ControlState {
	return &ControlState{
		SelectedTile: core.NewTile(-1, -1),
		RectStart:    core.NewTile(-1, -1),
		RectEnd:      core.NewTile(-1, -1),
	}
}

func (ct *ControlState) BeginRect(tile core.Tile) {
	if !tile.InBounds() {
		ct.SelectedTile = core.NewTile(-1, -1)
		ct.RectStart = core.NewTile(-1, -1)
		ct.RectEnd = core.NewTile(-1, -1)
		return
	}
	ct.SelectedTile = tile
	ct.RectStart = tile
	ct.RectEnd = tile
}

func (ct *ControlState) EndRect(tile core.Tile) {
	if !ct.SelectedTile.InBounds() {
		return
	}
	if tile.InBounds() {
		ct.RectEnd = tile
	}
	ct.SelectedTile = core.NewTile(-1, -1)
}

func (ct *ControlState) UpdateRect(tile core.Tile) {
	if !ct.SelectedTile.InBounds() || !tile.InBounds() {
		return
	}
	ct.RectEnd = tile
}

func (ct *ControlState) IsInRect(tileX, tileZ int) bool {
	if !ct.RectStart.InBounds() || !ct.RectEnd.InBounds() {
		return false
	}
	minX, maxX := ordered(ct.RectStart.X, ct.RectEnd.X)
	minZ, maxZ := ordered(ct.RectStart.Z, ct.RectEnd.Z)
	return tileX >= minX && tileX <= maxX && tileZ >= minZ && tileZ <= maxZ
}

func ordered(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}
