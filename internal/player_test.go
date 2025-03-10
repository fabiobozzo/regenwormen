package internal

import (
	"testing"

	"regenwormen/pkg/utils"
)

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name string
		mode PlayerMode
	}{
		{"human player", Human},
		{"AI player", AI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(tt.mode)

			if player.mode != tt.mode {
				t.Errorf("NewPlayer(%v).mode = %v, want %v", tt.mode, player.mode, tt.mode)
			}

			if player.tiles == nil {
				t.Errorf("NewPlayer() should initialize tiles")
			}
		})
	}
}

func TestPlayerString(t *testing.T) {
	// Player with no tiles
	emptyPlayer := NewPlayer(Human)
	if s := emptyPlayer.String(); s != "[ ]" {
		t.Errorf("Empty player String() = %q, want \"[ ]\"", s)
	}

	// Player with a tile
	playerWithTile := NewPlayer(Human)
	playerWithTile.tiles = utils.NewStack[Tile]()
	playerWithTile.tiles.Push(Tile{Value: 7, Worms: 2})

	if s := playerWithTile.String(); s != "[7]" {
		t.Errorf("Player with tile String() = %q, want \"[7]\"", s)
	}
}

func TestPlayerScore(t *testing.T) {
	player := NewPlayer(Human)

	// Add some tiles
	player.tiles.Push(Tile{Value: 4, Worms: 1})
	player.tiles.Push(Tile{Value: 6, Worms: 2})
	player.tiles.Push(Tile{Value: 8, Worms: 3})

	worms, values := player.Score()

	if worms != 6 {
		t.Errorf("Score() worms = %d, want 6", worms)
	}

	expectedValues := []int{8, 6, 4}
	if len(values) != len(expectedValues) {
		t.Errorf("Score() returned %d values, want %d", len(values), len(expectedValues))
	} else {
		for i, v := range values {
			if v != expectedValues[i] {
				t.Errorf("Score() values[%d] = %d, want %d", i, v, expectedValues[i])
			}
		}
	}
}
