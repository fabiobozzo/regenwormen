package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"regenwormen/internal"
	"regenwormen/pkg/utils"
)

func handleGameLoop(in *bufio.Reader, game *internal.Game) {
	var turnStopped bool

	currentPlayerNr, turnHasJustStarted, err := game.CurrentTurn()
	if err != nil {
		fmt.Println("Cannot determine current turn: ", err)
		game.Stop()

		return
	}

	if turnHasJustStarted {
		clearScreen()
		fmt.Printf("=== PLAYER %d TURN ===\n\n", currentPlayerNr)
	}

	fmt.Println(game.String())
	_ = utils.MustReadString(in, "Press the Enter ↵ key to roll the dice! ")

	for {
		roll := game.Dice.Roll()
		cannotPickFromRoll := !game.Dice.CanPickAnyFromRoll()
		if game.Dice.IsDone() || cannotPickFromRoll || turnStopped {
			if len(roll) > 0 && cannotPickFromRoll {
				fmt.Println("No symbols from last roll could be picked: ", game.Dice.StringRoll())
			}

			score, noWorms := game.Dice.PickedScore()
			fmt.Println()
			if noWorms {
				fmt.Println("You did not pick any worms. Therefore you did not score any points this turn.")
			} else {
				fmt.Printf("Player #%d scored %d points: %s\n", currentPlayerNr, score, game.Dice.StringPicked())
			}

			fmt.Print("Press the Enter key to continue.")
			if _, err = fmt.Scanln(); err != nil {
				log.Fatal(err)
			}

			game.NextTurn()
			turnStopped = false

			break
		}

		fmt.Println("Rolling the dice.... 🎲")
		fmt.Println(game.Dice.String())

		var picked bool
		for !picked {
			printSymbolPicker(roll, game)

			readInput := strings.TrimSpace(utils.MustReadString(in, " "))
			if readInput == "s" || readInput == "stop" {
				turnStopped = true

				break
			}

			var inputSymbol internal.Symbol
			inputSymbol, err = internal.SymbolFrom(readInput)
			if err != nil {
				fmt.Println("Please try again: ", err)

				continue
			}

			err = game.Dice.Pick(inputSymbol)
			if errors.Is(err, internal.ErrFullyPicked) || errors.Is(err, internal.ErrNoRollYet) {
				break
			}
			if err != nil {
				fmt.Println("Invalid pick: ", err)

				continue
			}

			picked = true
		}
	}
}
