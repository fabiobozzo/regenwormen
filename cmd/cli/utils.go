package main

import (
	"fmt"
	"log"
	"strings"

	"regenwormen/internal"
	"regenwormen/pkg/utils"
)

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

		fmt.Print(utils.MustEncloseCharAtIndex(utils.RemoveEmojis(strings.ToLower(rollSymbol.String())), firstCharIndex))
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
