package main

import (
	"bufio"
	"fmt"

	"regenwormen/internal"
	"regenwormen/pkg/utils"
)

func handleGameMenu(in *bufio.Reader, game *internal.Game) (exit bool) {
	doStart, answered := utils.MustReadBool(in, "Do you want to start the game?", "yes", "no")
	if !answered {
		return
	}
	if !doStart {
		fmt.Println("Player has exited the game")

		return true
	}

	nPlayers, err := utils.MustReadInt(in, "How many players? ")
	if err != nil {
		fmt.Println("Please enter a valid number of players:", err)

		return
	}

	if err = game.Start(nPlayers, 0); err != nil {
		fmt.Println("Cannot start the game:", err)

		return
	}

	return
}
