package game

import (
	"fmt"
)

type Color int

const (
	Empty Color = iota
	Black
	White
)

type Point struct {
	X, Y int
}

type Board struct {
	Size     int
	Grid     [][]Color
	LastMove *Point
	Captures map[Color]int
	History  []BoardState
	KoPoint  *Point
}

type BoardState struct {
	Grid     [][]Color
	Captures map[Color]int
	Move     *Point
	Player   Color
}

func NewBoard(size int) *Board {
	grid := make([][]Color, size)
	for i := range grid {
		grid[i] = make([]Color, size)
	}

	return &Board{
		Size: size,
		Grid: grid,
		Captures: map[Color]int{
			Black: 0,
			White: 0,
		},
		History: make([]BoardState, 0),
	}
}

func (b *Board) IsValidPoint(p Point) bool {
	return p.X >= 0 && p.X < b.Size && p.Y >= 0 && p.Y < b.Size
}

func (b *Board) GetColor(p Point) Color {
	if !b.IsValidPoint(p) {
		return Empty
	}
	return b.Grid[p.X][p.Y]
}

func (b *Board) SetStone(p Point, color Color) {
	if b.IsValidPoint(p) {
		b.Grid[p.X][p.Y] = color
	}
}

func (b *Board) GetNeighbors(p Point) []Point {
	neighbors := []Point{}
	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	
	for _, d := range directions {
		np := Point{p.X + d.X, p.Y + d.Y}
		if b.IsValidPoint(np) {
			neighbors = append(neighbors, np)
		}
	}
	return neighbors
}

func (b *Board) GetGroup(p Point) []Point {
	color := b.GetColor(p)
	if color == Empty {
		return []Point{}
	}

	group := []Point{}
	visited := make(map[Point]bool)
	queue := []Point{p}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		group = append(group, current)

		for _, neighbor := range b.GetNeighbors(current) {
			if !visited[neighbor] && b.GetColor(neighbor) == color {
				queue = append(queue, neighbor)
			}
		}
	}

	return group
}

func (b *Board) GetLiberties(group []Point) []Point {
	liberties := make(map[Point]bool)
	
	for _, stone := range group {
		for _, neighbor := range b.GetNeighbors(stone) {
			if b.GetColor(neighbor) == Empty {
				liberties[neighbor] = true
			}
		}
	}

	result := []Point{}
	for p := range liberties {
		result = append(result, p)
	}
	return result
}

func (b *Board) HasLiberties(p Point) bool {
	group := b.GetGroup(p)
	liberties := b.GetLiberties(group)
	return len(liberties) > 0
}

func (b *Board) RemoveGroup(group []Point) int {
	captured := len(group)
	for _, p := range group {
		b.SetStone(p, Empty)
	}
	return captured
}

func (b *Board) CaptureDeadGroups(color Color) int {
	opponent := OpponentColor(color)
	totalCaptured := 0

	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			p := Point{x, y}
			if b.GetColor(p) == opponent && !b.HasLiberties(p) {
				group := b.GetGroup(p)
				captured := b.RemoveGroup(group)
				totalCaptured += captured
			}
		}
	}

	return totalCaptured
}

func (b *Board) SaveState(move *Point, player Color) {
	gridCopy := make([][]Color, b.Size)
	for i := range gridCopy {
		gridCopy[i] = make([]Color, b.Size)
		copy(gridCopy[i], b.Grid[i])
	}

	capturesCopy := make(map[Color]int)
	for k, v := range b.Captures {
		capturesCopy[k] = v
	}

	b.History = append(b.History, BoardState{
		Grid:     gridCopy,
		Captures: capturesCopy,
		Move:     move,
		Player:   player,
	})
}

func (b *Board) IsKo(p Point, color Color) bool {
	if len(b.History) < 2 {
		return false
	}

	tempBoard := b.Clone()
	tempBoard.SetStone(p, color)
	tempBoard.CaptureDeadGroups(color)

	previousState := b.History[len(b.History)-2]
	
	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			if tempBoard.Grid[x][y] != previousState.Grid[x][y] {
				return false
			}
		}
	}
	
	return true
}

func (b *Board) Clone() *Board {
	gridCopy := make([][]Color, b.Size)
	for i := range gridCopy {
		gridCopy[i] = make([]Color, b.Size)
		copy(gridCopy[i], b.Grid[i])
	}

	capturesCopy := make(map[Color]int)
	for k, v := range b.Captures {
		capturesCopy[k] = v
	}

	clone := &Board{
		Size:     b.Size,
		Grid:     gridCopy,
		LastMove: b.LastMove,
		Captures: capturesCopy,
		KoPoint:  b.KoPoint,
	}

	return clone
}

func (b *Board) CountTerritory() map[Color]int {
	territory := map[Color]int{
		Black: 0,
		White: 0,
	}

	visited := make(map[Point]bool)

	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			p := Point{x, y}
			if b.GetColor(p) == Empty && !visited[p] {
				region, owner := b.findTerritoryRegion(p, visited)
				if owner != Empty {
					territory[owner] += len(region)
				}
			}
		}
	}

	return territory
}

func (b *Board) findTerritoryRegion(start Point, visited map[Point]bool) ([]Point, Color) {
	region := []Point{}
	owner := Empty
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

		if b.GetColor(p) == Empty {
			region = append(region, p)
			
			for _, neighbor := range b.GetNeighbors(p) {
				neighborColor := b.GetColor(neighbor)
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

	if hasBlack && !hasWhite {
		owner = Black
	} else if hasWhite && !hasBlack {
		owner = White
	}

	return region, owner
}

func (b *Board) String() string {
	result := "  "
	for i := 0; i < b.Size; i++ {
		result += fmt.Sprintf("%c ", 'A'+i)
	}
	result += "\n"

	for y := 0; y < b.Size; y++ {
		result += fmt.Sprintf("%2d ", b.Size-y)
		for x := 0; x < b.Size; x++ {
			switch b.Grid[x][y] {
			case Empty:
				result += ". "
			case Black:
				result += "● "
			case White:
				result += "○ "
			}
		}
		result += fmt.Sprintf("%d\n", b.Size-y)
	}

	result += "  "
	for i := 0; i < b.Size; i++ {
		result += fmt.Sprintf("%c ", 'A'+i)
	}

	return result
}

func OpponentColor(color Color) Color {
	if color == Black {
		return White
	}
	return Black
}

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Empty"
	}
}