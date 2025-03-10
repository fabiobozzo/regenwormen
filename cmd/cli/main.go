package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"regenwormen/internal"
)

func init() {
	initClearScreen()
}

func main() {
	game := internal.NewGame()
	in := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("=== WELCOME TO REGENWORMEN ===")
	fmt.Println()

	for {
		switch game.State {
		case internal.GameMenu:
			if exit := handleGameMenu(in, game); exit {
				return
			}
		case internal.GameLoop:
			handleGameLoop(in, game)
		case internal.GameOver:
			handleGameOver(game)
		default:
			log.Fatal("shutting down... unknown game state:", game.State)
		}
	}
}
