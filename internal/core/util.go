package core

const (
	Width  = 32
	Height = 20
)

type Tile struct {
	X, Z int
}

func NewTile(x, z int) Tile {
	return Tile{X: x, Z: z}
}

func (t *Tile) InBounds() bool {
	return InBounds(t.X, t.Z)
}

func InBounds(x, z int) bool {
	return x >= 0 && z >= 0 && x < Width && z < Height
}
