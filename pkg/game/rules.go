package game

import (
	"errors"
)

var (
	ErrInvalidMove       = errors.New("invalid move")
	ErrPositionOccupied  = errors.New("position already occupied")
	ErrSuicideMove       = errors.New("suicide move not allowed")
	ErrKoViolation       = errors.New("ko rule violation")
	ErrGameOver          = errors.New("game is over")
)

type Rules struct {
	AllowSuicide bool
	KoRule       bool
}

func NewRules() *Rules {
	return &Rules{
		AllowSuicide: false,
		KoRule:       true,
	}
}

type Game struct {
	Board        *Board
	Rules        *Rules
	CurrentTurn  Color
	Passed       map[Color]bool
	IsOver       bool
	Winner       *Color
	MoveCount    int
}

func NewGame(boardSize int) *Game {
	return &Game{
		Board:       NewBoard(boardSize),
		Rules:       NewRules(),
		CurrentTurn: Black,
		Passed: map[Color]bool{
			Black: false,
			White: false,
		},
		IsOver:    false,
		MoveCount: 0,
	}
}

func (g *Game) ValidateMove(p Point, color Color) error {
	if g.IsOver {
		return ErrGameOver
	}

	if color != g.CurrentTurn {
		return errors.New("not your turn")
	}

	if !g.Board.IsValidPoint(p) {
		return ErrInvalidMove
	}

	if g.Board.GetColor(p) != Empty {
		return ErrPositionOccupied
	}

	tempBoard := g.Board.Clone()
	tempBoard.SetStone(p, color)

	capturedOpponent := tempBoard.CaptureDeadGroups(color)

	if capturedOpponent == 0 && !tempBoard.HasLiberties(p) {
		if !g.Rules.AllowSuicide {
			return ErrSuicideMove
		}
		selfGroup := tempBoard.GetGroup(p)
		if len(selfGroup) > 0 && !tempBoard.HasLiberties(selfGroup[0]) {
			return ErrSuicideMove
		}
	}

	if g.Rules.KoRule && g.Board.IsKo(p, color) {
		return ErrKoViolation
	}

	return nil
}

func (g *Game) MakeMove(p Point, color Color) error {
	if err := g.ValidateMove(p, color); err != nil {
		return err
	}

	g.Board.SaveState(&p, color)

	g.Board.SetStone(p, color)

	captured := g.Board.CaptureDeadGroups(color)
	g.Board.Captures[color] += captured

	g.Board.LastMove = &p
	g.CurrentTurn = OpponentColor(color)
	g.MoveCount++

	g.Passed[color] = false

	g.Board.KoPoint = nil
	if captured == 1 {
		capturedGroup := []Point{}
		for x := 0; x < g.Board.Size; x++ {
			for y := 0; y < g.Board.Size; y++ {
				point := Point{x, y}
				if g.Board.GetColor(point) == Empty {
					neighbors := g.Board.GetNeighbors(point)
					surroundedByOpponent := true
					for _, n := range neighbors {
						if g.Board.GetColor(n) != OpponentColor(color) {
							surroundedByOpponent = false
							break
						}
					}
					if surroundedByOpponent && len(neighbors) == len(g.Board.GetNeighbors(point)) {
						capturedGroup = append(capturedGroup, point)
					}
				}
			}
		}
		if len(capturedGroup) == 1 {
			g.Board.KoPoint = &capturedGroup[0]
		}
	}

	return nil
}

func (g *Game) Pass(color Color) error {
	if g.IsOver {
		return ErrGameOver
	}

	if color != g.CurrentTurn {
		return errors.New("not your turn")
	}

	g.Passed[color] = true
	g.CurrentTurn = OpponentColor(color)

	g.Board.SaveState(nil, color)

	if g.Passed[Black] && g.Passed[White] {
		g.EndGame()
	}

	return nil
}

func (g *Game) EndGame() {
	g.IsOver = true
	scores := g.CalculateScore()
	
	if scores[Black] > scores[White] {
		winner := Black
		g.Winner = &winner
	} else if scores[White] > scores[Black] {
		winner := White
		g.Winner = &winner
	}
}

func (g *Game) CalculateScore() map[Color]int {
	territory := g.Board.CountTerritory()
	
	scores := map[Color]int{
		Black: territory[Black] + g.Board.Captures[Black],
		White: territory[White] + g.Board.Captures[White],
	}

	komi := 6.5
	scores[White] = int(float64(scores[White]) + komi)

	return scores
}

func (g *Game) GetValidMoves(color Color) []Point {
	validMoves := []Point{}
	
	for x := 0; x < g.Board.Size; x++ {
		for y := 0; y < g.Board.Size; y++ {
			p := Point{x, y}
			if g.ValidateMove(p, color) == nil {
				validMoves = append(validMoves, p)
			}
		}
	}
	
	return validMoves
}

func (g *Game) Resign(color Color) {
	g.IsOver = true
	winner := OpponentColor(color)
	g.Winner = &winner
}

func (g *Game) GetGameInfo() map[string]interface{} {
	info := map[string]interface{}{
		"boardSize":    g.Board.Size,
		"currentTurn":  g.CurrentTurn.String(),
		"moveCount":    g.MoveCount,
		"isOver":       g.IsOver,
		"blackCaptures": g.Board.Captures[Black],
		"whiteCaptures": g.Board.Captures[White],
	}

	if g.IsOver && g.Winner != nil {
		info["winner"] = g.Winner.String()
		info["scores"] = g.CalculateScore()
	}

	return info
}

func (g *Game) GetBoardState() [][]Color {
	state := make([][]Color, g.Board.Size)
	for i := range state {
		state[i] = make([]Color, g.Board.Size)
		copy(state[i], g.Board.Grid[i])
	}
	return state
}

func (g *Game) Undo() bool {
	if len(g.Board.History) <= 1 {
		return false
	}

	g.Board.History = g.Board.History[:len(g.Board.History)-1]
	
	if len(g.Board.History) > 0 {
		lastState := g.Board.History[len(g.Board.History)-1]
		
		for i := range g.Board.Grid {
			copy(g.Board.Grid[i], lastState.Grid[i])
		}
		
		for k, v := range lastState.Captures {
			g.Board.Captures[k] = v
		}
		
		g.Board.LastMove = lastState.Move
		g.CurrentTurn = OpponentColor(lastState.Player)
		g.MoveCount--
		
		g.Passed[Black] = false
		g.Passed[White] = false
		g.IsOver = false
		g.Winner = nil
		
		return true
	}
	
	return false
}