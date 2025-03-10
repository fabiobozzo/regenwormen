package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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

	currentPlayer := game.CurrentPlayer()

	// AI Turn
	if currentPlayer.IsAI() {
		// Let AI make all its decisions for this turn
		for {
			fmt.Println("🤖 AI thinking whether to roll the dice or not...")
			time.Sleep(500 * time.Millisecond)
			shouldRoll, explanation := currentPlayer.AiThink(game)
			fmt.Println("❗️", explanation)

			if !shouldRoll {
				break // AI decides to stop rolling
			}

			game.Dice.Roll()
			fmt.Println("Rolling the dice.... 🎲")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Roll:", game.Dice.StringRoll())
			fmt.Println()

			// If no valid picks available, turn ends with no points
			if !game.Dice.CanPickAnyFromRoll() {
				fmt.Println("No symbols can be picked from this roll - turn ends with no points! 🤷")

				break
			}

			// AI picks one symbol
			fmt.Println("🤖 AI thinking which symbol to pick...")
			time.Sleep(500 * time.Millisecond)
			symbol, explanation := currentPlayer.AiChoosePick(game)
			fmt.Println("❗️", explanation)

			// No valid choice, turn ends
			if symbol < 0 {
				break
			}

			if err := game.Dice.Pick(symbol); err != nil {
				log.Printf("invalid AI pick: %v\n", err)

				break
			}
			fmt.Println("Picked:", game.Dice.StringPicked())
			fmt.Println()
		}

		fmt.Println()
		_ = utils.MustReadString(in, "Press the Enter ↵ key to continue.")
		game.NextTurn()

		return
	}

	// Human Turn
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

			_ = utils.MustReadString(in, "Press the Enter ↵ key to continue.")
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
