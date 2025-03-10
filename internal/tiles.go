package internal

import (
	"errors"
	"fmt"
	"strings"
)

var ErrCannotTakeFromBoard = errors.New("cannot take a tile from the board")

var defaultTiles = []Tile{
	{Value: 4, Worms: 1},
	{Value: 4, Worms: 1},
	{Value: 5, Worms: 1},
	{Value: 5, Worms: 1},
	{Value: 6, Worms: 2},
	{Value: 6, Worms: 2},
	{Value: 7, Worms: 2},
	{Value: 7, Worms: 2},
	{Value: 8, Worms: 3},
	{Value: 8, Worms: 3},
	{Value: 9, Worms: 4},
	{Value: 9, Worms: 4},
}

type Tile struct {
	Value int
	Worms int
}

type Board struct {
	tiles map[int][]Tile
	min   int
	max   int
}

func NewDefaultBoard() (board *Board) {
	board = &Board{
		min:   1_000_000,
		max:   0,
		tiles: map[int][]Tile{},
	}
	for _, t := range defaultTiles {
		board.tiles[t.Value] = append(board.tiles[t.Value], t)
		if t.Value > board.max {
			board.max = t.Value
		}
		if t.Value < board.min {
			board.min = t.Value
		}
	}

	return
}

func (b *Board) Take(val int) (t Tile, err error) {
	tilesForValue, exist := b.tiles[val]
	if !exist {
		return t, ErrCannotTakeFromBoard
	}

	if len(tilesForValue) > 0 {
		t = b.tiles[val][len(b.tiles[val])-1]
		b.tiles[val] = b.tiles[val][:len(b.tiles[val])-1]
	}
	if len(b.tiles[val]) == 0 {
		delete(b.tiles, val)
	}

	return
}

func (b *Board) IsEmpty() bool {
	return len(b.tiles) == 0
}

func (b *Board) String() string {
	var sb strings.Builder
	for i := b.min; i <= b.max; i++ {
		for _, t := range b.tiles[i] {
			sb.WriteString(fmt.Sprintf("[%d] ", t.Value))
		}
	}

	return sb.String()
}
