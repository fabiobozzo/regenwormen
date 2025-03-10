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
	ai    AIStrategy
}

func NewPlayer(mode PlayerMode) Player {
	var ai AIStrategy
	if mode == AI {
		ai = NewSimpleAIStrategy()
	}

	return Player{
		mode:  mode,
		tiles: utils.NewStack[Tile](),
		ai:    ai,
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

func (p Player) IsAI() bool {
	return p.mode == AI
}

func (p Player) AiThink(game *Game) (shouldRoll bool, explanation string) {
	if !p.IsAI() || p.ai == nil {
		return false, ""
	}

	return p.ai.ShouldRoll(game)
}

func (p Player) AiChoosePick(game *Game) (symbol Symbol, explanation string) {
	if !p.IsAI() || p.ai == nil {
		return -1, ""
	}

	return p.ai.ChooseSymbol(game)
}
