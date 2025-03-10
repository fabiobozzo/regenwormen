package internal

import (
	"errors"
	"testing"
)

func TestNewDice(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  int
	}{
		{"default count", 0, DefaultDiceCount},
		{"custom count", 8, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDice(tt.count)
			if d.count != tt.want {
				t.Errorf("NewDice(%d).count = %v, want %v", tt.count, d.count, tt.want)
			}
		})
	}
}

func TestDiceReset(t *testing.T) {
	d := NewDice(DefaultDiceCount)
	d.roll = []Symbol{Worm, Bread}
	d.picked = []Symbol{Cucumber}

	d.Reset()

	if len(d.roll) != 0 || len(d.picked) != 0 {
		t.Errorf("Dice.Reset() didn't clear roll and picked slices")
	}
}

func TestDiceRoll(t *testing.T) {
	d := NewDice(6)
	d.picked = []Symbol{Worm, Bread} // 2 dice already picked

	roll := d.Roll()

	if len(roll) != 4 { // Should roll 4 dice (6 total - 2 picked)
		t.Errorf("Dice.Roll() returned %d dice, want %d", len(roll), 4)
	}

	if len(d.roll) != 4 {
		t.Errorf("Dice.roll has %d dice after Roll(), want %d", len(d.roll), 4)
	}
}

func TestDiceIsDone(t *testing.T) {
	d := NewDice(3)

	if d.IsDone() {
		t.Errorf("New dice should not be done")
	}

	d.picked = []Symbol{Worm, Bread, Cucumber}

	if !d.IsDone() {
		t.Errorf("Dice should be done when all dice are picked")
	}
}

func TestDiceCanPick(t *testing.T) {
	d := NewDice(DefaultDiceCount)
	d.picked = []Symbol{Worm}

	if !d.CanPick(Bread) {
		t.Errorf("Should be able to pick Bread")
	}

	if d.CanPick(Worm) {
		t.Errorf("Should not be able to pick Worm again")
	}
}

func TestDiceCanPickAnyFromRoll(t *testing.T) {
	d := NewDice(DefaultDiceCount)

	// No roll yet
	if d.CanPickAnyFromRoll() {
		t.Errorf("Should not be able to pick from empty roll")
	}

	// Roll with pickable symbols
	d.roll = []Symbol{Worm, Bread}
	if !d.CanPickAnyFromRoll() {
		t.Errorf("Should be able to pick from roll with available symbols")
	}

	// Roll with only already picked symbols
	d.picked = []Symbol{Worm, Bread}
	d.roll = []Symbol{Worm, Bread}
	if d.CanPickAnyFromRoll() {
		t.Errorf("Should not be able to pick when all symbols in roll are already picked")
	}
}

func TestDicePick(t *testing.T) {
	tests := []struct {
		name    string
		dice    *Dice
		pick    Symbol
		wantErr error
		wantLen int
	}{
		{
			name:    "pick valid symbol",
			dice:    &Dice{count: 6, roll: []Symbol{Worm, Worm, Bread}, picked: []Symbol{}},
			pick:    Worm,
			wantErr: nil,
			wantLen: 2, // Should pick both Worms
		},
		{
			name:    "pick already picked symbol",
			dice:    &Dice{count: 6, roll: []Symbol{Worm, Bread}, picked: []Symbol{Worm}},
			pick:    Worm,
			wantErr: ErrDoublePick,
			wantLen: 1,
		},
		{
			name:    "pick symbol not in roll",
			dice:    &Dice{count: 6, roll: []Symbol{Worm, Bread}, picked: []Symbol{}},
			pick:    Cucumber,
			wantErr: ErrPickMustBeInRoll,
			wantLen: 0,
		},
		{
			name:    "pick with no roll",
			dice:    &Dice{count: 6, roll: []Symbol{}, picked: []Symbol{}},
			pick:    Worm,
			wantErr: ErrNoRollYet,
			wantLen: 0,
		},
		{
			name:    "pick when all dice picked",
			dice:    &Dice{count: 3, roll: []Symbol{Worm}, picked: []Symbol{Bread, Cucumber, Cheese}},
			pick:    Worm,
			wantErr: ErrFullyPicked,
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dice.Pick(tt.pick)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Dice.Pick() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.dice.picked) != tt.wantLen {
				t.Errorf("After Pick(), len(picked) = %d, want %d", len(tt.dice.picked), tt.wantLen)
			}

			if err == nil && len(tt.dice.roll) != 0 {
				t.Errorf("After successful Pick(), roll should be empty, got %v", tt.dice.roll)
			}
		})
	}
}

func TestDicePickedScore(t *testing.T) {
	tests := []struct {
		name        string
		picked      []Symbol
		wantScore   int
		wantNoWorms bool
	}{
		{
			name:        "no worms",
			picked:      []Symbol{Bread, Cucumber, Ketchup},
			wantScore:   0,
			wantNoWorms: true,
		},
		{
			name:        "with worms, no bread",
			picked:      []Symbol{Worm, Cucumber, Ketchup},
			wantScore:   3,
			wantNoWorms: false,
		},
		{
			name:        "with worms and bread",
			picked:      []Symbol{Worm, Bread, Bread, Cucumber},
			wantScore:   6, // 4 dice - 2 bread + (2 bread * 2) = 6
			wantNoWorms: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDice(DefaultDiceCount)
			d.picked = tt.picked

			score, noWorms := d.PickedScore()

			if score != tt.wantScore {
				t.Errorf("PickedScore() score = %v, want %v", score, tt.wantScore)
			}

			if noWorms != tt.wantNoWorms {
				t.Errorf("PickedScore() noWorms = %v, want %v", noWorms, tt.wantNoWorms)
			}
		})
	}
}

func TestSymbolString(t *testing.T) {
	tests := []struct {
		symbol Symbol
		want   string
	}{
		{Worm, "Worm 🐛"},
		{Bread, "Bread 🥖"},
		{Cucumber, "Cucumber 🥒"},
		{Ketchup, "Ketchup 🥫"},
		{Cheese, "Cheese 🧀"},
		{Symbol(99), "Unknown 🤷🏻‍"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.symbol.String(); got != tt.want {
				t.Errorf("Symbol.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSymbolFrom(t *testing.T) {
	tests := []struct {
		input    string
		want     Symbol
		wantErr  bool
		errCheck func(error) bool
	}{
		{"w", Worm, false, nil},
		{"W", Worm, false, nil},
		{"worm", Worm, false, nil},
		{"b", Bread, false, nil},
		{"B", Bread, false, nil},
		{"bread", Bread, false, nil},
		{"c", Cucumber, false, nil},
		{"C", Cucumber, false, nil},
		{"cucumber", Cucumber, false, nil},
		{"k", Ketchup, false, nil},
		{"K", Ketchup, false, nil},
		{"ketchup", Ketchup, false, nil},
		{"h", Cheese, false, nil},
		{"H", Cheese, false, nil},
		{"cheese", Cheese, false, nil},
		{"", -1, true, func(err error) bool { return err == ErrInvalidSymbol }},
		{"x", -1, true, func(err error) bool { return err != nil }},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := SymbolFrom(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SymbolFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("SymbolFrom() error = %v, doesn't match expected error condition", err)
				return
			}
			if got != tt.want {
				t.Errorf("SymbolFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
