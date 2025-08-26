package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Prawal-Sharma/GoSim/pkg/game"
	"github.com/Prawal-Sharma/GoSim/pkg/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	hub := websocket.NewHub()
	go hub.Run()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(hub, w, r)
	})

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	r.Post("/api/ai-move", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Board      [][]int `json:"board"`
			BoardSize  int     `json:"boardSize"`
			Color      string  `json:"color"`
			Difficulty string  `json:"difficulty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		boardGame := game.NewGame(req.BoardSize)
		for x := 0; x < req.BoardSize; x++ {
			for y := 0; y < req.BoardSize; y++ {
				color := game.Color(req.Board[x][y])
				if color != game.Empty {
					boardGame.Board.SetStone(game.Point{X: x, Y: y}, color)
				}
			}
		}

		var playerColor game.Color
		if req.Color == "Black" {
			playerColor = game.Black
		} else {
			playerColor = game.White
		}

		move := getAIMove(boardGame, playerColor, req.Difficulty)

		w.Header().Set("Content-Type", "application/json")
		if move != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"x": move.X,
				"y": move.Y,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"pass": true,
			})
		}
	})

	r.Get("/api/puzzles", func(w http.ResponseWriter, r *http.Request) {
		puzzles := loadPuzzles()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(puzzles)
	})

	r.Get("/api/lessons", func(w http.ResponseWriter, r *http.Request) {
		lessons := loadLessons()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lessons)
	})

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getAIMove(g *game.Game, color game.Color, difficulty string) *game.Point {
	validMoves := g.GetValidMoves(color)
	
	if len(validMoves) == 0 {
		return nil
	}

	switch difficulty {
	case "random":
		return &validMoves[0]
	case "easy":
		bestMove := validMoves[0]
		bestScore := 0
		
		for _, move := range validMoves {
			score := evaluateMove(g, move, color)
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		}
		return &bestMove
	default:
		return &validMoves[len(validMoves)/2]
	}
}

func evaluateMove(g *game.Game, move game.Point, color game.Color) int {
	tempGame := &game.Game{
		Board:       g.Board.Clone(),
		Rules:       g.Rules,
		CurrentTurn: color,
		Passed:      make(map[game.Color]bool),
	}

	tempGame.MakeMove(move, color)
	
	score := 0
	
	if tempGame.Board.Captures[color] > g.Board.Captures[color] {
		score += (tempGame.Board.Captures[color] - g.Board.Captures[color]) * 10
	}
	
	group := tempGame.Board.GetGroup(move)
	liberties := tempGame.Board.GetLiberties(group)
	score += len(liberties) * 2
	
	cornerBonus := 0
	if (move.X < 3 || move.X > g.Board.Size-4) && (move.Y < 3 || move.Y > g.Board.Size-4) {
		cornerBonus = 5
	}
	score += cornerBonus
	
	return score
}

func loadPuzzles() []map[string]interface{} {
	puzzles := []map[string]interface{}{
		{
			"id":          1,
			"title":       "Basic Capture",
			"description": "Capture the white stone",
			"difficulty":  "beginner",
			"board":       createPuzzleBoard(9, "capture1"),
		},
		{
			"id":          2,
			"title":       "Ladder Problem",
			"description": "Can black capture the white stone with a ladder?",
			"difficulty":  "intermediate",
			"board":       createPuzzleBoard(9, "ladder1"),
		},
	}
	return puzzles
}

func loadLessons() []map[string]interface{} {
	lessons := []map[string]interface{}{
		{
			"id":          1,
			"title":       "Introduction to Go",
			"description": "Learn the basics of Go",
			"content":     "Go is an ancient board game...",
			"level":       "beginner",
		},
		{
			"id":          2,
			"title":       "Capturing Stones",
			"description": "Learn how to capture opponent stones",
			"content":     "Stones are captured when they have no liberties...",
			"level":       "beginner",
		},
		{
			"id":          3,
			"title":       "Life and Death",
			"description": "Understanding when groups are alive or dead",
			"content":     "A group is alive if it has two eyes...",
			"level":       "intermediate",
		},
	}
	return lessons
}

func createPuzzleBoard(size int, puzzleType string) [][]int {
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}

	switch puzzleType {
	case "capture1":
		if size >= 5 {
			board[4][4] = 2
			board[3][4] = 1
			board[5][4] = 1
			board[4][3] = 1
		}
	case "ladder1":
		if size >= 7 {
			board[3][3] = 2
			board[2][3] = 1
			board[3][2] = 1
			board[4][3] = 1
		}
	}

	return board
}