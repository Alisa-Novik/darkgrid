package main

import (
	"testing"

	"adagrad/internal/game"
)

func TestGameSelection(t *testing.T) {
	g := game.NewGame(6, 4)

	g.SelectTile(2, 1)
	gotX, gotZ := g.SelectedTile()
	if gotX != 2 || gotZ != 1 {
		t.Fatalf("SelectTile kept (%d,%d), want (2,1)", gotX, gotZ)
	}

	g.SelectTile(10, 10)
	gotX, gotZ = g.SelectedTile()
	if gotX != -1 || gotZ != -1 {
		t.Fatalf("SelectTile cleared to (%d,%d), want (-1,-1)", gotX, gotZ)
	}
}

func TestGameBounds(t *testing.T) {
	g := game.NewGame(3, 2)

	if !g.InBounds(0, 0) {
		t.Fatal("expected (0,0) to be in bounds")
	}
	if g.InBounds(3, 1) {
		t.Fatal("expected (3,1) to be out of bounds")
	}
}
