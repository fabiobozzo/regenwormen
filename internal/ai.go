package internal

import (
	"fmt"
	"slices"
)

// AIStrategy defines the interface for different AI decision-making strategies
type AIStrategy interface {
	ShouldRoll(game *Game) (shouldRoll bool, explanation string)
	ChooseSymbol(game *Game) (symbol Symbol, explanation string)
}

// SimpleAIStrategy implements a basic strategy
type SimpleAIStrategy struct {
	thresholdScore int
}

func NewSimpleAIStrategy() *SimpleAIStrategy {
	return &SimpleAIStrategy{
		thresholdScore: 5, // arbitrary
	}
}

func (s *SimpleAIStrategy) ShouldRoll(game *Game) (bool, string) {
	// If no dice picked yet, always roll
	if len(game.Dice.picked) == 0 {
		return true, "First roll of the turn"
	}

	score, noWorms := game.Dice.PickedScore()

	// If no worms yet, must continue
	if noWorms {
		return true, "Must continue rolling - no worms picked yet"
	}

	// Check if current score matches an available tile
	if game.board.HasTile(score) {
		return false, fmt.Sprintf("Stopping - can get a tile with score %d", score)
	}

	// Check if AI can steal from other players decks
	for i, p := range game.players {
		if i == game.turn { // Skip current player
			continue
		}
		if top, exists := p.tiles.Top(); exists && top.Value == score {
			return false, fmt.Sprintf("Stopping - can steal tile with value %d", score)
		}
	}

	// Check if we can get a lower value tile
	for i := score - 1; i >= game.board.min; i-- {
		if game.board.HasTile(i) {
			return false, fmt.Sprintf("Stopping - can get a lower value tile with score %d", i)
		}
	}

	// Continue rolling if score is too low
	if score < s.thresholdScore {
		return true, fmt.Sprintf("Continue rolling - score %d is too low (threshold: %d)", score, s.thresholdScore)
	}

	// Stop if score is high enough but no tiles available
	return false, fmt.Sprintf("Stopping - scored %d but got no tiles! 🤷", score)
}

func (s *SimpleAIStrategy) ChooseSymbol(game *Game) (Symbol, string) {
	// First priority: Pick worms if we don't have any
	if !slices.Contains(game.Dice.picked, Worm) {
		if slices.Contains(game.Dice.roll, Worm) {
			return Worm, "Picking Worm - need at least one worm to score"
		}
	}

	// Second priority: Pick bread (worth 2 points)
	if slices.Contains(game.Dice.roll, Bread) && game.Dice.CanPick(Bread) {
		return Bread, "Picking Bread - worth 2 points each"
	}

	// Third priority: Pick most frequent symbol
	symbolCounts := make(map[Symbol]int)
	for _, s := range game.Dice.roll {
		if game.Dice.CanPick(s) {
			symbolCounts[s]++
		}
	}

	var bestSymbol Symbol
	maxCount := 0
	for sym, count := range symbolCounts {
		if count > maxCount {
			maxCount = count
			bestSymbol = sym
		}
	}

	return bestSymbol, "Picking most frequent symbol to maximize score"
}
