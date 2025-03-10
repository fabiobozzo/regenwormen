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

	var humanPlayers, aiPlayers int
	var err error

	for {
		humanPlayers, err = utils.MustReadInt(in, "How many human players? ")
		if err != nil {
			fmt.Println("Please enter a valid number:", err)

			continue
		}

		aiPlayers, err = utils.MustReadInt(in, "How many AI players? ")
		if err != nil {
			fmt.Println("Please enter a valid number:", err)

			continue
		}

		break
	}

	if err = game.Start(humanPlayers, aiPlayers); err != nil {
		fmt.Println("Cannot start the game: ", err)
		return
	}

	return
}
