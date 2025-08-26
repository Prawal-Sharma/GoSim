package game

import (
	"math"
	"math/rand"
	"time"
)

type AI struct {
	Color      Color
	Difficulty string
	Game       *Game
}

func NewAI(color Color, difficulty string) *AI {
	rand.Seed(time.Now().UnixNano())
	return &AI{
		Color:      color,
		Difficulty: difficulty,
	}
}

func (ai *AI) GetMove(game *Game) *Point {
	ai.Game = game
	
	switch ai.Difficulty {
	case "random":
		return ai.getRandomMove()
	case "easy":
		return ai.getEasyMove()
	case "medium":
		return ai.getMediumMove()
	case "hard":
		return ai.getHardMove()
	default:
		return ai.getEasyMove()
	}
}

func (ai *AI) getRandomMove() *Point {
	validMoves := ai.Game.GetValidMoves(ai.Color)
	if len(validMoves) == 0 {
		return nil
	}
	
	randomIndex := rand.Intn(len(validMoves))
	return &validMoves[randomIndex]
}

func (ai *AI) getEasyMove() *Point {
	validMoves := ai.Game.GetValidMoves(ai.Color)
	if len(validMoves) == 0 {
		return nil
	}
	
	var bestMove *Point
	bestScore := -1000
	
	for _, move := range validMoves {
		score := ai.evaluateMove(move)
		if score > bestScore {
			bestScore = score
			bestMove = &move
		}
	}
	
	return bestMove
}

func (ai *AI) getMediumMove() *Point {
	validMoves := ai.Game.GetValidMoves(ai.Color)
	if len(validMoves) == 0 {
		return nil
	}
	
	var candidates []ScoredMove
	
	for _, move := range validMoves {
		score := ai.evaluateMoveAdvanced(move)
		candidates = append(candidates, ScoredMove{Move: move, Score: score})
	}
	
	sortMovesByScore(candidates)
	
	topCount := minInt(5, len(candidates))
	topMoves := candidates[:topCount]
	
	if len(topMoves) > 0 {
		selected := topMoves[rand.Intn(len(topMoves))]
		return &selected.Move
	}
	
	return &validMoves[0]
}

func (ai *AI) getHardMove() *Point {
	validMoves := ai.Game.GetValidMoves(ai.Color)
	if len(validMoves) == 0 {
		return nil
	}
	
	var bestMove *Point
	bestScore := math.Inf(-1)
	
	for _, move := range validMoves {
		score := ai.minimax(move, 3, math.Inf(-1), math.Inf(1), true)
		if score > bestScore {
			bestScore = score
			bestMove = &move
		}
	}
	
	return bestMove
}

func (ai *AI) evaluateMove(move Point) int {
	score := 0
	
	tempGame := &Game{
		Board:       ai.Game.Board.Clone(),
		Rules:       ai.Game.Rules,
		CurrentTurn: ai.Color,
		Passed:      make(map[Color]bool),
	}
	
	for k, v := range ai.Game.Board.Captures {
		tempGame.Board.Captures[k] = v
	}
	
	tempGame.Board.SetStone(move, ai.Color)
	captured := tempGame.Board.CaptureDeadGroups(ai.Color)
	score += captured * 10
	
	group := tempGame.Board.GetGroup(move)
	liberties := tempGame.Board.GetLiberties(group)
	score += len(liberties) * 2
	
	score += ai.getPositionScore(move)
	
	for _, neighbor := range tempGame.Board.GetNeighbors(move) {
		if tempGame.Board.GetColor(neighbor) == ai.Color {
			score += 3
		}
	}
	
	opponent := OpponentColor(ai.Color)
	for x := 0; x < tempGame.Board.Size; x++ {
		for y := 0; y < tempGame.Board.Size; y++ {
			p := Point{x, y}
			if tempGame.Board.GetColor(p) == opponent {
				oppGroup := tempGame.Board.GetGroup(p)
				oppLiberties := tempGame.Board.GetLiberties(oppGroup)
				if len(oppLiberties) == 1 {
					score += 5
				}
			}
		}
	}
	
	return score
}

func (ai *AI) evaluateMoveAdvanced(move Point) float64 {
	score := float64(ai.evaluateMove(move))
	
	tempGame := &Game{
		Board:       ai.Game.Board.Clone(),
		Rules:       ai.Game.Rules,
		CurrentTurn: ai.Color,
		Passed:      make(map[Color]bool),
	}
	
	for k, v := range ai.Game.Board.Captures {
		tempGame.Board.Captures[k] = v
	}
	
	tempGame.MakeMove(move, ai.Color)
	
	territory := tempGame.Board.CountTerritory()
	territoryScore := float64(territory[ai.Color] - territory[OpponentColor(ai.Color)])
	score += territoryScore * 0.5
	
	influenceScore := ai.calculateInfluence(tempGame.Board, move)
	score += influenceScore * 0.3
	
	if ai.isEye(tempGame.Board, move, ai.Color) {
		score -= 20
	}
	
	if ai.makesEye(tempGame.Board, move, ai.Color) {
		score += 15
	}
	
	if ai.connectsGroups(tempGame.Board, move, ai.Color) {
		score += 12
	}
	
	if ai.cutsOpponentGroups(tempGame.Board, move, ai.Color) {
		score += 18
	}
	
	return score
}

func (ai *AI) minimax(move Point, depth int, alpha, beta float64, maximizing bool) float64 {
	if depth == 0 {
		return ai.evaluateMoveAdvanced(move)
	}
	
	tempGame := &Game{
		Board:       ai.Game.Board.Clone(),
		Rules:       ai.Game.Rules,
		CurrentTurn: ai.Color,
		Passed:      make(map[Color]bool),
	}
	
	color := ai.Color
	if !maximizing {
		color = OpponentColor(ai.Color)
	}
	
	err := tempGame.MakeMove(move, color)
	if err != nil {
		return 0
	}
	
	validMoves := tempGame.GetValidMoves(OpponentColor(color))
	
	if len(validMoves) == 0 {
		return ai.evaluateBoardState(tempGame.Board)
	}
	
	if maximizing {
		maxScore := math.Inf(-1)
		for _, nextMove := range validMoves {
			score := ai.minimax(nextMove, depth-1, alpha, beta, false)
			maxScore = math.Max(maxScore, score)
			alpha = math.Max(alpha, score)
			if beta <= alpha {
				break
			}
		}
		return maxScore
	} else {
		minScore := math.Inf(1)
		for _, nextMove := range validMoves {
			score := ai.minimax(nextMove, depth-1, alpha, beta, true)
			minScore = math.Min(minScore, score)
			beta = math.Min(beta, score)
			if beta <= alpha {
				break
			}
		}
		return minScore
	}
}

func (ai *AI) getPositionScore(move Point) int {
	size := ai.Game.Board.Size
	score := 0
	
	distToEdge := minInt(move.X, minInt(move.Y, minInt(size-1-move.X, size-1-move.Y)))
	
	if size == 19 {
		if distToEdge <= 2 {
			if (move.X <= 2 || move.X >= 16) && (move.Y <= 2 || move.Y >= 16) {
				score += 8
			} else {
				score += 4
			}
		} else if distToEdge == 3 {
			if (move.X == 3 || move.X == 15) && (move.Y == 3 || move.Y == 15) {
				score += 10
			} else {
				score += 5
			}
		}
	} else if size == 13 {
		if distToEdge <= 2 {
			score += 6
		} else if distToEdge == 3 {
			score += 8
		}
	} else if size == 9 {
		if distToEdge <= 1 {
			score += 5
		} else if distToEdge == 2 {
			score += 7
		}
	}
	
	return score
}

func (ai *AI) calculateInfluence(board *Board, move Point) float64 {
	influence := 0.0
	maxDistance := 5
	
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			if board.Grid[x][y] == Empty {
				distance := absInt(x-move.X) + absInt(y-move.Y)
				if distance <= maxDistance {
					influence += 1.0 / float64(distance+1)
				}
			}
		}
	}
	
	return influence
}

func (ai *AI) isEye(board *Board, point Point, color Color) bool {
	if board.GetColor(point) != Empty {
		return false
	}
	
	neighbors := board.GetNeighbors(point)
	for _, n := range neighbors {
		if board.GetColor(n) != color {
			return false
		}
	}
	
	diagonals := []Point{
		{point.X - 1, point.Y - 1},
		{point.X + 1, point.Y - 1},
		{point.X - 1, point.Y + 1},
		{point.X + 1, point.Y + 1},
	}
	
	opponentCount := 0
	for _, d := range diagonals {
		if board.IsValidPoint(d) && board.GetColor(d) == OpponentColor(color) {
			opponentCount++
		}
	}
	
	return opponentCount <= 1
}

func (ai *AI) makesEye(board *Board, move Point, color Color) bool {
	tempBoard := board.Clone()
	tempBoard.SetStone(move, color)
	
	for _, neighbor := range tempBoard.GetNeighbors(move) {
		if ai.isEye(tempBoard, neighbor, color) {
			return true
		}
	}
	
	return false
}

func (ai *AI) connectsGroups(board *Board, move Point, color Color) bool {
	adjacentGroups := make(map[*[]Point]bool)
	
	for _, neighbor := range board.GetNeighbors(move) {
		if board.GetColor(neighbor) == color {
			group := board.GetGroup(neighbor)
			adjacentGroups[&group] = true
		}
	}
	
	return len(adjacentGroups) >= 2
}

func (ai *AI) cutsOpponentGroups(board *Board, move Point, color Color) bool {
	opponent := OpponentColor(color)
	adjacentOpponentGroups := make(map[*[]Point]bool)
	
	for _, neighbor := range board.GetNeighbors(move) {
		if board.GetColor(neighbor) == opponent {
			group := board.GetGroup(neighbor)
			adjacentOpponentGroups[&group] = true
		}
	}
	
	return len(adjacentOpponentGroups) >= 2
}

func (ai *AI) evaluateBoardState(board *Board) float64 {
	territory := board.CountTerritory()
	captures := board.Captures
	
	aiScore := float64(territory[ai.Color] + captures[ai.Color])
	oppScore := float64(territory[OpponentColor(ai.Color)] + captures[OpponentColor(ai.Color)])
	
	if ai.Color == White {
		aiScore += 6.5
	} else {
		oppScore += 6.5
	}
	
	return aiScore - oppScore
}

type ScoredMove struct {
	Move  Point
	Score float64
}

func sortMovesByScore(moves []ScoredMove) {
	for i := 0; i < len(moves); i++ {
		for j := i + 1; j < len(moves); j++ {
			if moves[j].Score > moves[i].Score {
				moves[i], moves[j] = moves[j], moves[i]
			}
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}