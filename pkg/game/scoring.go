package game

import "fmt"

type ScoringMethod string

const (
	ChineseScoring  ScoringMethod = "chinese"
	JapaneseScoring ScoringMethod = "japanese"
)

type Score struct {
	Black      float64
	White      float64
	Territory  map[Color]int
	Captures   map[Color]int
	Komi       float64
	Method     ScoringMethod
	Winner     *Color
	Difference float64
}

func CalculateScore(board *Board, method ScoringMethod, komi float64) *Score {
	score := &Score{
		Territory: board.CountTerritory(),
		Captures:  board.Captures,
		Komi:      komi,
		Method:    method,
	}

	switch method {
	case ChineseScoring:
		score.calculateChineseScore(board)
	case JapaneseScoring:
		score.calculateJapaneseScore(board)
	default:
		score.calculateChineseScore(board)
	}

	score.determineWinner()
	return score
}

func (s *Score) calculateChineseScore(board *Board) {
	blackStones := 0
	whiteStones := 0

	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			color := board.Grid[x][y]
			if color == Black {
				blackStones++
			} else if color == White {
				whiteStones++
			}
		}
	}

	s.Black = float64(blackStones + s.Territory[Black])
	s.White = float64(whiteStones+s.Territory[White]) + s.Komi
}

func (s *Score) calculateJapaneseScore(board *Board) {
	s.Black = float64(s.Territory[Black] + s.Captures[Black])
	s.White = float64(s.Territory[White]+s.Captures[White]) + s.Komi
}

func (s *Score) determineWinner() {
	s.Difference = s.Black - s.White

	if s.Black > s.White {
		winner := Black
		s.Winner = &winner
	} else if s.White > s.Black {
		winner := White
		s.Winner = &winner
	}
}

func (s *Score) GetResult() string {
	if s.Winner == nil {
		return "Draw"
	}

	winnerStr := "B"
	if *s.Winner == White {
		winnerStr = "W"
	}

	diff := s.Difference
	if diff < 0 {
		diff = -diff
	}

	return winnerStr + "+" + formatFloat(diff)
}

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return fmt.Sprintf("%d", int(f))
	}
	return fmt.Sprintf("%.1f", f)
}

type DeadStoneMarker struct {
	Groups [][]Point
	Stones map[Point]bool
}

func MarkDeadStones(board *Board) *DeadStoneMarker {
	marker := &DeadStoneMarker{
		Groups: make([][]Point, 0),
		Stones: make(map[Point]bool),
	}

	visited := make(map[Point]bool)

	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			p := Point{x, y}
			if !visited[p] && board.GetColor(p) != Empty {
				group := board.GetGroup(p)
				if isGroupDead(board, group) {
					marker.Groups = append(marker.Groups, group)
					for _, stone := range group {
						marker.Stones[stone] = true
						visited[stone] = true
					}
				} else {
					for _, stone := range group {
						visited[stone] = true
					}
				}
			}
		}
	}

	return marker
}

func isGroupDead(board *Board, group []Point) bool {
	if len(group) == 0 {
		return false
	}

	liberties := board.GetLiberties(group)
	if len(liberties) == 0 {
		return true
	}

	color := board.GetColor(group[0])
	eyeCount := 0

	for _, liberty := range liberties {
		if isEye(board, liberty, color) {
			eyeCount++
		}
	}

	return eyeCount < 2
}

func isEye(board *Board, point Point, color Color) bool {
	if board.GetColor(point) != Empty {
		return false
	}

	neighbors := board.GetNeighbors(point)
	for _, n := range neighbors {
		if board.GetColor(n) != color {
			return false
		}
	}

	corners := []Point{
		{point.X - 1, point.Y - 1},
		{point.X + 1, point.Y - 1},
		{point.X - 1, point.Y + 1},
		{point.X + 1, point.Y + 1},
	}

	friendlyCorners := 0
	totalCorners := 0

	for _, c := range corners {
		if board.IsValidPoint(c) {
			totalCorners++
			if board.GetColor(c) == color {
				friendlyCorners++
			}
		}
	}

	if totalCorners == 4 {
		return friendlyCorners >= 3
	}
	return friendlyCorners == totalCorners
}

func EstimateTerritory(board *Board) map[Point]Color {
	territory := make(map[Point]Color)
	visited := make(map[Point]bool)

	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			p := Point{x, y}
			if board.GetColor(p) == Empty && !visited[p] {
				region, owner := findTerritoryRegion(board, p, visited)
				for _, point := range region {
					territory[point] = owner
				}
			}
		}
	}

	return territory
}

func findTerritoryRegion(board *Board, start Point, visited map[Point]bool) ([]Point, Color) {
	region := []Point{}
	hasBlack := false
	hasWhite := false
	queue := []Point{start}
	localVisited := make(map[Point]bool)

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		if localVisited[p] {
			continue
		}
		localVisited[p] = true
		visited[p] = true

		if board.GetColor(p) == Empty {
			region = append(region, p)

			for _, neighbor := range board.GetNeighbors(p) {
				neighborColor := board.GetColor(neighbor)
				if neighborColor == Empty && !localVisited[neighbor] {
					queue = append(queue, neighbor)
				} else if neighborColor == Black {
					hasBlack = true
				} else if neighborColor == White {
					hasWhite = true
				}
			}
		}
	}

	owner := Empty
	if hasBlack && !hasWhite {
		owner = Black
	} else if hasWhite && !hasBlack {
		owner = White
	}

	return region, owner
}

type GameResult struct {
	Score      *Score
	DeadStones *DeadStoneMarker
	Territory  map[Point]Color
	SGF        string
}

func GetGameResult(game *Game, method ScoringMethod, komi float64) *GameResult {
	deadStones := MarkDeadStones(game.Board)
	
	boardCopy := game.Board.Clone()
	for stone := range deadStones.Stones {
		boardCopy.SetStone(stone, Empty)
	}
	
	territory := EstimateTerritory(boardCopy)
	score := CalculateScore(boardCopy, method, komi)
	
	return &GameResult{
		Score:      score,
		DeadStones: deadStones,
		Territory:  territory,
		SGF:        generateSGF(game),
	}
}

func generateSGF(game *Game) string {
	sgf := fmt.Sprintf("(;FF[4]GM[1]SZ[%d]", game.Board.Size)
	
	for i, state := range game.Board.History {
		if state.Move != nil {
			color := "B"
			if state.Player == White {
				color = "W"
			}
			x := string(rune('a' + state.Move.X))
			y := string(rune('a' + state.Move.Y))
			sgf += ";" + color + "[" + x + y + "]"
		} else if i > 0 {
			color := "B"
			if state.Player == White {
				color = "W"
			}
			sgf += ";" + color + "[]"
		}
	}
	
	sgf += ")"
	return sgf
}