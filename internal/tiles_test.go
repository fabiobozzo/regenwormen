package internal

import (
	"errors"
	"testing"
)

func TestNewDefaultBoard(t *testing.T) {
	board := NewDefaultBoard()

	if board.min > board.max {
		t.Errorf("Board min (%d) should be <= max (%d)", board.min, board.max)
	}

	if len(board.tiles) == 0 {
		t.Errorf("Default board should have tiles")
	}
}

func TestBoardTake(t *testing.T) {
	board := NewDefaultBoard()

	// Remember initial state
	initialTiles := make(map[int]int)
	for val, tiles := range board.tiles {
		initialTiles[val] = len(tiles)
	}

	// Take a tile that exists
	for val := range board.tiles {
		tile, err := board.Take(val)

		if err != nil {
			t.Errorf("Take(%d) returned error: %v", val, err)
		}

		if tile.Value != val {
			t.Errorf("Take(%d) returned tile with value %d", val, tile.Value)
		}

		// Check that the count decreased
		if len(board.tiles[val]) != initialTiles[val]-1 {
			t.Errorf("After Take(%d), count should be %d, got %d",
				val, initialTiles[val]-1, len(board.tiles[val]))
		}

		break // Just test one value
	}

	// Take a tile that doesn't exist
	nonExistentValue := board.max + 1
	_, err := board.Take(nonExistentValue)

	if !errors.Is(err, ErrCannotTakeFromBoard) {
		t.Errorf("Take(%d) should return ErrCannotTakeFromBoard, got %v",
			nonExistentValue, err)
	}
}

func TestBoardIsEmpty(t *testing.T) {
	board := NewDefaultBoard()

	if board.IsEmpty() {
		t.Errorf("New default board should not be empty")
	}

	// Empty the board
	for val := range board.tiles {
		for len(board.tiles[val]) > 0 {
			_, _ = board.Take(val)
		}
	}

	if !board.IsEmpty() {
		t.Errorf("Board should be empty after taking all tiles")
	}
}
