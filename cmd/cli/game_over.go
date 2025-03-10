package main

import (
	"fmt"

	"regenwormen/internal"
)

func handleGameOver(game *internal.Game) {
	clearScreen()
	fmt.Println("=== GAME OVER ===")
	fmt.Println()

	printWinner(game)

	game.Restart()
}
