package game

type Game struct {
	width     int
	height    int
	tiles     [][]int
	selectedX int
	selectedZ int
}

func NewGame(width, height int) *Game {
	return &Game{
		width:     width,
		height:    height,
		tiles:     makeMap(width, height),
		selectedX: -1,
		selectedZ: -1,
	}
}

func (g *Game) Width() int {
	return g.width
}

func (g *Game) Height() int {
	return g.height
}

func (g *Game) Tiles() [][]int {
	return g.tiles
}

func (g *Game) InBounds(x, z int) bool {
	return x >= 0 && z >= 0 && x < g.width && z < g.height
}

func (g *Game) SelectTile(x, z int) {
	if !g.InBounds(x, z) {
		g.selectedX = -1
		g.selectedZ = -1
		return
	}
	g.selectedX = x
	g.selectedZ = z
}

func (g *Game) SelectedTile() (int, int) {
	return g.selectedX, g.selectedZ
}

func makeMap(w, h int) [][]int {
	m := make([][]int, h)
	for z := 0; z < h; z++ {
		m[z] = make([]int, w)
		for x := 0; x < w; x++ {
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
