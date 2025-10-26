package game

import (
	"adagrad/internal/ui"
	"time"
)

type Game struct {
	tiles     [][]int
	selectedX int
	selectedZ int
}

func NewGame(width, height int) *Game {
	return &Game{
		tiles:     makeMap(width, height),
		selectedX: -1,
		selectedZ: -1,
	}
}

func (g *Game) Tiles() [][]int {
	return g.tiles
}

func (g *Game) SelectTile(x, z int) {
	if g.InBounds(x, z) {
		g.selectedX = x
		g.selectedZ = z
		return
	}

	g.selectedX = -1
	g.selectedZ = -1
}

func (g *Game) SelectedTile() (int, int) {
	return g.selectedX, g.selectedZ
}

func (g *Game) InBounds(x, z int) bool {
	if len(g.tiles) == 0 {
		return false
	}

	if z < 0 || z >= len(g.tiles) {
		return false
	}

	row := g.tiles[z]
	return x >= 0 && x < len(row)
}

func makeMap(w, h int) [][]int {
	m := make([][]int, h)
	for z := range h {
		m[z] = make([]int, w)
		for x := range w {
			if x == 0 || z == 0 || x == w-1 || z == h-1 {
				m[z][x] = 1
			}
		}
	}
	for x := 3; x < w-3; x++ {
		m[h/2][x] = 1
	}
	for z := 3; z < h-3; z++ {
		m[z][w/3] = 1
	}
	return m
}

func (gm *Game) Run() {
	ui.Prepare()

	prev := time.Now()
	for !ui.ShouldClose() {
		now := time.Now()
		dt := float32(now.Sub(prev).Seconds())
		prev = now

		tiles := gm.Tiles()
		ui.Render(dt, tiles)
	}
}
