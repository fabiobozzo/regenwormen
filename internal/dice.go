package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"

	"regenwormen/pkg/utils"
)

type Symbol int

const (
	Worm Symbol = iota
	Bread
	Cucumber
	Ketchup
	Cheese

	DefaultDiceCount = 6
)

var (
	ErrFullyPicked      = errors.New("you already picked all the Dice for this set")
	ErrNoRollYet        = errors.New("please roll the Dice first")
	ErrPickMustBeInRoll = errors.New("symbol was not rolled")
	ErrDoublePick       = errors.New("symbol was already picked for this set")
	ErrInvalidSymbol    = errors.New("not a valid symbol")
)

func (s Symbol) String() string {
	switch s {
	case Worm:
		return "Worm 🐛"
	case Bread:
		return "Bread 🥖"
	case Cucumber:
		return "Cucumber 🥒"
	case Ketchup:
		return "Ketchup 🥫"
	case Cheese:
		return "Cheese 🧀"
	default:
		return "Unknown 🤷🏻‍"
	}
}

func SymbolFrom(s string) (Symbol, error) {
	if len(s) == 0 {
		return -1, ErrInvalidSymbol
	}

	switch strings.TrimSpace(strings.ToLower(fmt.Sprintf("%c", s[0]))) {
	case "w":
		return Worm, nil
	case "b":
		return Bread, nil
	case "c":
		return Cucumber, nil
	case "k":
		return Ketchup, nil
	case "h":
		return Cheese, nil
	default:
		return -1, fmt.Errorf("%s %w", s, ErrInvalidSymbol)
	}
}

type Dice struct {
	count  int
	roll   []Symbol
	picked []Symbol
}

func NewDice(count int) *Dice {
	if count <= 0 {
		count = DefaultDiceCount
	}

	return &Dice{count: count}
}

func (d *Dice) Reset() {
	d.picked = nil
	d.roll = nil
}

func (d *Dice) Roll() []Symbol {
	d.roll = nil

	for i := 0; i < d.count-len(d.picked); i++ {
		throw := rand.Intn(5) - 1
		if throw < 0 {
			throw = 0
		}

		d.roll = append(d.roll, Symbol(throw))
	}

	return d.roll
}

func (d *Dice) IsDone() bool {
	return len(d.picked) == d.count
}

func (d *Dice) CanPick(s Symbol) bool {
	return !slices.Contains(d.picked, s)
}

func (d *Dice) CanPickAnyFromRoll() (can bool) {
	for _, s := range d.roll {
		if d.CanPick(s) {
			can = true

			break
		}
	}

	return
}

func (d *Dice) Pick(s Symbol) error {
	if d.IsDone() {
		return ErrFullyPicked
	}

	if len(d.roll) == 0 {
		return ErrNoRollYet
	}

	if !slices.Contains(d.roll, s) {
		return ErrPickMustBeInRoll
	}

	if !d.CanPick(s) {
		return ErrDoublePick
	}

	pickedCount := utils.CountOccurrences(d.roll, func(rs Symbol) bool { return rs == s })

	for i := 0; i < pickedCount; i++ {
		d.picked = append(d.picked, s)
	}
	d.roll = nil

	return nil
}

func (d *Dice) PickedScore() (score int, noWorms bool) {
	if utils.CountOccurrences(d.picked, func(ps Symbol) bool { return ps == Worm }) == 0 {
		return 0, true
	}

	breads := utils.CountOccurrences(d.picked, func(ps Symbol) bool { return ps == Bread })
	score = len(d.picked) - breads + (breads * 2)

	return
}

func (d *Dice) String() string {
	var sb strings.Builder

	if len(d.picked) > 0 {
		sb.WriteString("Picked: ")
		sb.WriteString(d.StringPicked())
		sb.WriteString("\n")
	}

	sb.WriteString("Roll: ")
	sb.WriteString(d.StringRoll())

	return sb.String()
}

func (d *Dice) StringPicked() string {
	var sb strings.Builder

	if len(d.picked) == 0 {
		sb.WriteString("[]")
	}
	for _, s := range d.picked {
		sb.WriteString(fmt.Sprintf("[%s] ", s.String()))
	}

	return sb.String()
}

func (d *Dice) StringRoll() string {
	var sb strings.Builder

	for _, s := range d.roll {
		sb.WriteString(fmt.Sprintf("[%s] ", s.String()))
	}

	return sb.String()
}
