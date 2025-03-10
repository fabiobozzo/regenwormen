package internal

import (
	"errors"
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	if game.State != GameMenu {
		t.Errorf("New game state = %v, want %v", game.State, GameMenu)
	}

	if game.Dice == nil {
		t.Errorf("New game should have dice")
	}

	if game.board == nil {
		t.Errorf("New game should have a board")
	}
}

func TestGameStart(t *testing.T) {
	tests := []struct {
		name         string
		humanPlayers int
		aiPlayers    int
		wantErr      bool
	}{
		{"valid game", 2, 0, false},
		{"valid mixed game", 1, 1, false},
		{"not enough players", 1, 0, true},
		{"zero players", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			err := game.Start(tt.humanPlayers, tt.aiPlayers)

			if (err != nil) != tt.wantErr {
				t.Errorf("Start(%d, %d) error = %v, wantErr %v",
					tt.humanPlayers, tt.aiPlayers, err, tt.wantErr)
			}

			if err == nil {
				expectedPlayers := tt.humanPlayers + tt.aiPlayers
				if len(game.players) != expectedPlayers {
					t.Errorf("After Start(), len(players) = %d, want %d",
						len(game.players), expectedPlayers)
				}

				if game.State != GameLoop {
					t.Errorf("After Start(), State = %v, want %v",
						game.State, GameLoop)
				}
			}
		})
	}
}

func TestGameRestart(t *testing.T) {
	game := NewGame()
	_ = game.Start(2, 0)
	game.turn = 1

	game.Restart()

	if game.State != GameMenu {
		t.Errorf("After Restart(), State = %v, want %v", game.State, GameMenu)
	}

	if len(game.players) != 0 {
		t.Errorf("After Restart(), should have no players, got %d", len(game.players))
	}

	if game.turn != 0 {
		t.Errorf("After Restart(), turn = %d, want 0", game.turn)
	}
}

func TestGameCurrentTurn(t *testing.T) {
	game := NewGame()

	// Test before game starts
	_, _, err := game.CurrentTurn()
	if !errors.Is(err, ErrGameOver) {
		t.Errorf("CurrentTurn() before game starts should return ErrGameOver, got %v", err)
	}

	// Test during game
	_ = game.Start(2, 0)
	playerN, justStarted, err := game.CurrentTurn()

	if err != nil {
		t.Errorf("CurrentTurn() during game returned error: %v", err)
	}

	if playerN != 1 {
		t.Errorf("CurrentTurn() playerN = %d, want 1", playerN)
	}

	if !justStarted {
		t.Errorf("CurrentTurn() justStarted = %v, want true", justStarted)
	}

	// Add some picked dice and check again
	game.Dice.picked = []Symbol{Worm}
	_, justStarted, _ = game.CurrentTurn()

	if justStarted {
		t.Errorf("CurrentTurn() with picked dice, justStarted = %v, want false", justStarted)
	}
}

func TestGameNextTurn(t *testing.T) {
	game := NewGame()
	_ = game.Start(3, 0)

	// First turn should be player 1 (index 0)
	if game.turn != 0 {
		t.Errorf("Initial turn = %d, want 0", game.turn)
	}

	// Move to next turn
	game.NextTurn()

	if game.turn != 1 {
		t.Errorf("After NextTurn(), turn = %d, want 1", game.turn)
	}

	// Test wrapping around
	game.turn = 2
	game.NextTurn()

	if game.turn != 0 {
		t.Errorf("After NextTurn() from last player, turn = %d, want 0", game.turn)
	}
}

func TestGameResolveCurrentTurn(t *testing.T) {
	tests := []struct {
		name          string
		setupGame     func(*Game)
		wantTileValue int
		wantTileCount int
		wantWorms     int
	}{
		{
			name: "no worms - should not get tile",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				g.Dice.picked = []Symbol{Bread, Cucumber, Ketchup}
			},
			wantTileValue: 0,
			wantTileCount: 0,
			wantWorms:     0,
		},
		{
			name: "with worms - should get exact tile",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				g.Dice.picked = []Symbol{Worm, Bread, Bread, Cucumber} // Score: 6
			},
			wantTileValue: 6,
			wantTileCount: 1,
			wantWorms:     2,
		},
		{
			name: "high score but tile not available - should get lower tile",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				// Remove tiles with value 8 from board
				_, _ = g.board.Take(8)
				_, _ = g.board.Take(8)
				// Roll score of 8
				g.Dice.picked = []Symbol{Worm, Worm, Bread, Bread, Bread} // Score: 8
			},
			wantTileValue: 7, // Should get next available lower tile
			wantTileCount: 1,
			wantWorms:     2,
		},
		{
			name: "can steal matching tile from opponent",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				// Remove tiles with value 6 from board
				_, _ = g.board.Take(6)
				_, _ = g.board.Take(6)
				// Give opponent a tile
				g.players[1].tiles.Push(Tile{Value: 6, Worms: 2})
				// Roll matching score
				g.Dice.picked = []Symbol{Worm, Cucumber, Bread, Bread} // Score: 6
			},
			wantTileValue: 6,
			wantTileCount: 1,
			wantWorms:     2,
		},
		{
			name: "prefer board tile over stealing",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				// Give opponent a tile
				g.players[1].tiles.Push(Tile{Value: 5, Worms: 1})
				// Roll matching score, but tile also available on board
				g.Dice.picked = []Symbol{Worm, Bread, Bread} // Score: 5
			},
			wantTileValue: 5,
			wantTileCount: 1,
			wantWorms:     1,
		},
		{
			name: "no matching tile anywhere - get nothing",
			setupGame: func(g *Game) {
				g.Start(2, 0)
				// Empty the board
				for val := g.board.min; val <= g.board.max; val++ {
					for len(g.board.tiles[val]) > 0 {
						_, _ = g.board.Take(val)
					}
				}
				g.Dice.picked = []Symbol{Worm, Bread, Bread} // Score: 5
			},
			wantTileValue: 0,
			wantTileCount: 0,
			wantWorms:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			tt.setupGame(game)

			initialPlayerTiles := game.players[game.turn].tiles.Len()
			game.resolveCurrentTurn()

			player := game.players[game.turn]
			gotTiles := player.tiles.Len()

			if gotTiles != initialPlayerTiles+tt.wantTileCount {
				t.Errorf("Player got %d tiles, want %d",
					gotTiles-initialPlayerTiles, tt.wantTileCount)
			}

			if tt.wantTileCount > 0 {
				if topTile, exists := player.tiles.Top(); !exists {
					t.Error("Expected player to have a tile, but got none")
				} else {
					if topTile.Value != tt.wantTileValue {
						t.Errorf("Player got tile value %d, want %d",
							topTile.Value, tt.wantTileValue)
					}
					if topTile.Worms != tt.wantWorms {
						t.Errorf("Player got tile with %d worms, want %d",
							topTile.Worms, tt.wantWorms)
					}
				}
			}
		})
	}
}
