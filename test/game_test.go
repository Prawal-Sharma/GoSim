package test

import (
	"testing"

	"github.com/Prawal-Sharma/GoSim/pkg/game"
)

func TestBoardCreation(t *testing.T) {
	sizes := []int{9, 13, 19}
	
	for _, size := range sizes {
		board := game.NewBoard(size)
		
		if board.Size != size {
			t.Errorf("Expected board size %d, got %d", size, board.Size)
		}
		
		if len(board.Grid) != size {
			t.Errorf("Expected grid length %d, got %d", size, len(board.Grid))
		}
		
		for i, row := range board.Grid {
			if len(row) != size {
				t.Errorf("Expected row %d length %d, got %d", i, size, len(row))
			}
		}
	}
}

func TestBasicCapture(t *testing.T) {
	g := game.NewGame(9)
	
	// Setup a capture scenario
	// Black surrounds white
	moves := []struct {
		x, y  int
		color game.Color
	}{
		{4, 4, game.White},
		{3, 4, game.Black},
		{5, 4, game.Black},
		{4, 3, game.Black},
		{4, 5, game.Black}, // This should capture white
	}
	
	for i, move := range moves[:4] {
		g.Board.SetStone(game.Point{X: move.x, Y: move.y}, move.color)
		if i%2 == 0 {
			g.CurrentTurn = game.Black
		} else {
			g.CurrentTurn = game.White
		}
	}
	
	// Make the capturing move
	err := g.MakeMove(game.Point{X: 4, Y: 5}, game.Black)
	if err != nil {
		t.Errorf("Failed to make capturing move: %v", err)
	}
	
	// Check if white stone was captured
	if g.Board.GetColor(game.Point{X: 4, Y: 4}) != game.Empty {
		t.Error("White stone should have been captured")
	}
	
	if g.Board.Captures[game.Black] != 1 {
		t.Errorf("Black should have 1 capture, got %d", g.Board.Captures[game.Black])
	}
}

func TestKoRule(t *testing.T) {
	g := game.NewGame(9)
	
	// Setup a ko situation
	koSetup := []struct {
		x, y int
		color game.Color
	}{
		{3, 3, game.Black},
		{4, 3, game.White},
		{5, 3, game.Black},
		{3, 4, game.White},
		{5, 4, game.White},
		{3, 5, game.Black},
		{4, 5, game.White},
		{5, 5, game.Black},
		{4, 4, game.Black}, // Black captures at 4,4
	}
	
	for i, move := range koSetup {
		if i == len(koSetup)-1 {
			// Last move should be a proper move
			err := g.MakeMove(game.Point{X: move.x, Y: move.y}, move.color)
			if err != nil {
				t.Errorf("Failed to make move: %v", err)
			}
		} else {
			g.Board.SetStone(game.Point{X: move.x, Y: move.y}, move.color)
		}
	}
	
	// Try to immediately recapture (should violate ko rule)
	err := g.MakeMove(game.Point{X: 4, Y: 4}, game.White)
	if err != game.ErrKoViolation {
		t.Error("Ko rule should prevent immediate recapture")
	}
}

func TestSuicideRule(t *testing.T) {
	g := game.NewGame(9)
	
	// Setup a suicide scenario
	suicideSetup := []struct {
		x, y int
		color game.Color
	}{
		{4, 3, game.White},
		{3, 4, game.White},
		{5, 4, game.White},
		{4, 5, game.White},
	}
	
	for _, move := range suicideSetup {
		g.Board.SetStone(game.Point{X: move.x, Y: move.y}, move.color)
	}
	
	// Try to play a suicide move
	err := g.MakeMove(game.Point{X: 4, Y: 4}, game.Black)
	if err != game.ErrSuicideMove {
		t.Error("Suicide move should not be allowed")
	}
}

func TestPassAndGameEnd(t *testing.T) {
	g := game.NewGame(9)
	
	// Both players pass
	err := g.Pass(game.Black)
	if err != nil {
		t.Errorf("Black pass failed: %v", err)
	}
	
	if g.IsOver {
		t.Error("Game should not be over after one pass")
	}
	
	err = g.Pass(game.White)
	if err != nil {
		t.Errorf("White pass failed: %v", err)
	}
	
	if !g.IsOver {
		t.Error("Game should be over after both players pass")
	}
}

func TestTerritoryCalculation(t *testing.T) {
	board := game.NewBoard(9)
	
	// Create a simple territory scenario
	// Black controls top-left corner
	for x := 0; x < 4; x++ {
		board.SetStone(game.Point{X: x, Y: 3}, game.Black)
		board.SetStone(game.Point{X: 3, Y: x}, game.Black)
	}
	
	territory := board.CountTerritory()
	
	if territory[game.Black] == 0 {
		t.Error("Black should have some territory")
	}
}

func TestAIMove(t *testing.T) {
	g := game.NewGame(9)
	
	difficulties := []string{"random", "easy", "medium", "hard"}
	
	for _, diff := range difficulties {
		ai := game.NewAI(game.Black, diff)
		move := ai.GetMove(g)
		
		if move == nil {
			t.Errorf("AI (%s) should return a valid move on empty board", diff)
		}
		
		if move != nil {
			err := g.ValidateMove(*move, game.Black)
			if err != nil {
				t.Errorf("AI (%s) returned invalid move: %v", diff, err)
			}
		}
	}
}

func TestScoring(t *testing.T) {
	board := game.NewBoard(9)
	
	// Simple scoring scenario
	board.Captures[game.Black] = 5
	board.Captures[game.White] = 3
	
	score := game.CalculateScore(board, game.ChineseScoring, 6.5)
	
	if score.White == score.Black {
		t.Error("Scores should differ with komi")
	}
	
	if score.Komi != 6.5 {
		t.Errorf("Expected komi 6.5, got %f", score.Komi)
	}
}