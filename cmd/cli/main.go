package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"regenwormen/internal"
	"regenwormen/pkg/utils"
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
			doStart := utils.MustReadString(in, "Do you want to start the game (yes/no)? ")
			doStart = strings.TrimSpace(strings.ToLower(doStart))
			if doStart == "yes" || doStart == "y" {
				nPlayers, err := utils.MustReadInt(in, "How many players? ")
				if err != nil {
					fmt.Println("Please enter a valid number of players:", err)

					continue
				}

				err = game.Start(nPlayers, 0)
				if err != nil {
					fmt.Println("Cannot start the game:", err)
				}
			} else if doStart == "no" || doStart == "n" {
				return
			}
		case internal.GameLoop:
			currentPlayerNr, turnHasJustStarted, err := game.CurrentTurn()
			if err != nil {
				fmt.Println(err)
				game.Stop()

				continue
			}

			if turnHasJustStarted {
				clearScreen()
				fmt.Printf("=== PLAYER %d TURN ===\n\n", currentPlayerNr)
			}

			fmt.Println(game.String())
			fmt.Print("Press the Enter key to roll the dice! ")
			if _, err = fmt.Scanln(); err != nil {
				log.Fatal(err)
			}

			var turnStopped bool
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
		case internal.GameOver:
			clearScreen()
			fmt.Println("=== GAME OVER ===")
			fmt.Println()

			printWinner(game)

			game.Restart()
		default:
			log.Fatal("shutting down... unknown game state:", game.State)
		}
	}
}

func printSymbolPicker(roll []internal.Symbol, game *internal.Game) {
	printed := map[internal.Symbol]struct{}{}

	fmt.Print("\nPick a symbol or (s)top here: ")
	var i int
	for _, rollSymbol := range roll {
		if !game.Dice.CanPick(rollSymbol) {
			continue
		}

		if _, alreadyPrinted := printed[rollSymbol]; alreadyPrinted {
			continue
		}
		printed[rollSymbol] = struct{}{}

		if i > 0 {
			fmt.Print(", ")
		}

		var firstCharIndex int
		if rollSymbol == internal.Cheese {
			firstCharIndex = 1
		}

		fmt.Print(utils.MustEncloseCharAtIndex(utils.RemoveEmojis(rollSymbol.String()), firstCharIndex))
		i++
	}
}

func printWinner(game *internal.Game) {
	players, err := game.FinalScores()
	if err != nil {
		log.Fatal(err)
	}

	iWinner := -1
	wormsWinner := 0
	tie := false

	for i, p := range players {
		playerWorms, playerTiles := p.Score()
		fmt.Printf("P%d captured %d worms with tiles:", i+1, playerWorms)
		for _, tv := range playerTiles {
			fmt.Printf(" [%d]", tv)
		}
		fmt.Println()

		if playerWorms > wormsWinner {
			wormsWinner = playerWorms
			iWinner = i
			tie = false
		} else if playerWorms == wormsWinner {
			tie = true
		}
	}

	if tie {
		fmt.Println("TIE! 🤝")
	} else {
		fmt.Printf("PLAYER #%d WINS! 🎉\n\n", iWinner+1)
	}
}
