package internal

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrGameOver    = errors.New("the game is over")
	ErrGameNotOver = errors.New("the game is not over yet")
)

type GameState int

const (
	GameMenu GameState = iota
	GameLoop
	GameOver
)

type Game struct {
	State   GameState
	Dice    *Dice
	players []Player
	turn    int
	board   *Board
}

func NewGame() *Game {
	return &Game{
		State: GameMenu,
		Dice:  NewDice(DefaultDiceCount),
		turn:  0,
		board: NewDefaultBoard(),
	}
}

func (g *Game) Start(humanPlayers, aiPlayers int) (err error) {
	if humanPlayers+aiPlayers < 2 {
		return errors.New("a game must have at least 2 players")
	}

	for i := 0; i < humanPlayers; i++ {
		g.players = append(g.players, NewPlayer(Human))
	}
	for i := 0; i < aiPlayers; i++ {
		g.players = append(g.players, NewPlayer(AI))
	}

	g.State = GameLoop
	g.turn = 0

	return
}

func (g *Game) Restart() {
	g.State = GameMenu

	g.players = nil
	g.turn = 0
	g.board = NewDefaultBoard()
	g.Dice.Reset()
}

func (g *Game) Stop() {
	g.State = GameOver
}

func (g *Game) CurrentTurn() (playerN int, justStarted bool, err error) {
	if g.State != GameLoop {
		return 0, false, ErrGameOver
	}

	return g.turn + 1, len(g.Dice.picked) == 0, nil
}

func (g *Game) NextTurn() {
	g.resolveCurrentTurn()
	g.Dice.Reset()

	if g.board.IsEmpty() {
		g.Stop()

		return
	}

	g.turn++
	if len(g.players) == g.turn {
		g.turn = 0
	}

	return
}

func (g *Game) FinalScores() ([]Player, error) {
	if g.State != GameOver {
		return nil, ErrGameNotOver
	}

	return g.players, nil
}

func (g *Game) String() string {
	var sb strings.Builder
	sb.WriteString("Board: ")
	sb.WriteString(g.board.String())
	sb.WriteString("\nPlayers: ")
	for i, p := range g.players {
		sb.WriteString(fmt.Sprintf("P%d:%s ", i+1, p.String()))
	}
	sb.WriteString("\n")

	return sb.String()
}

func (g *Game) resolveCurrentTurn() {
	diceScore, noWorms := g.Dice.PickedScore()
	if diceScore == 0 || noWorms {
		return
	}

	// Get the tile with the scored dice value from the board.
	tile, err := g.board.Take(diceScore)

	if err != nil {
		// If there is no available tile on the board, then try to rob from other players decks.
		var robbed bool
		for _, p := range g.players {
			top, hasTiles := p.tiles.Top()
			if hasTiles && top.Value == diceScore {
				tile, robbed = p.tiles.Pop()

				break
			}
		}

		// If there is nothing to rob then pick a lower value tile from the board, if available.
		if !robbed {
			for i := diceScore - 1; i >= g.board.min; i-- {
				tile, err = g.board.Take(i)
				if err == nil {
					break
				}
			}
		}
	}

	if tile.Value != 0 {
		g.players[g.turn].tiles.Push(tile)
	}

	return
}
