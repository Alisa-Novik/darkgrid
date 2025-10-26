package core

const (
	Width = 32
	Height = 20
)

func InBounds(x, z int) bool {
	return x >= 0 && z >= 0 && x < Width && z < Height
}
