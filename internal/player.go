package internal

import (
	"fmt"

	"regenwormen/pkg/utils"
)

type PlayerMode int

const (
	AI PlayerMode = iota
	Human
)

type Player struct {
	mode  PlayerMode
	tiles *utils.Stack[Tile]
}

func NewPlayer(mode PlayerMode) Player {
	return Player{
		mode:  mode,
		tiles: utils.NewStack[Tile](),
	}
}

func (p Player) String() string {
	var s string
	topTile, exists := p.tiles.Top()
	if !exists {
		s = " "
	} else {
		s = fmt.Sprintf("%d", topTile.Value)
	}

	return fmt.Sprintf("[%s]", s)
}

func (p Player) Score() (worms int, values []int) {
	tile, popped := p.tiles.Pop()
	for popped {
		worms += tile.Worms
		values = append(values, tile.Value)

		tile, popped = p.tiles.Pop()
	}

	return
}
