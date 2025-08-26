package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

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

	// Serve static files
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	
	// Serve index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

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
	ai := game.NewAI(color, difficulty)
	return ai.GetMove(g)
}


func loadPuzzles() []map[string]interface{} {
	puzzles := []map[string]interface{}{}
	
	// Try to load from JSON file first
	puzzleFile := filepath.Join("data", "puzzles", "beginner.json")
	data, err := ioutil.ReadFile(puzzleFile)
	if err == nil {
		var loadedPuzzles []map[string]interface{}
		if json.Unmarshal(data, &loadedPuzzles) == nil {
			puzzles = append(puzzles, loadedPuzzles...)
		}
	}
	
	// Add default puzzles if no file found
	if len(puzzles) == 0 {
		puzzles = []map[string]interface{}{
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
	}
	return puzzles
}

func loadLessons() []map[string]interface{} {
	lessons := []map[string]interface{}{}
	
	// Try to load from JSON file first
	lessonFile := filepath.Join("data", "lessons", "basics.json")
	data, err := ioutil.ReadFile(lessonFile)
	if err == nil {
		var loadedLessons []map[string]interface{}
		if json.Unmarshal(data, &loadedLessons) == nil {
			lessons = append(lessons, loadedLessons...)
		}
	}
	
	// Add default lessons if no file found
	if len(lessons) == 0 {
		lessons = []map[string]interface{}{
			{
				"id":          1,
				"title":       "Introduction to Go",
				"description": "Learn the basics of Go",
				"content":     "Go is an ancient board game originated in China over 4000 years ago. The objective is to control more territory than your opponent.",
				"level":       "beginner",
			},
			{
				"id":          2,
				"title":       "Capturing Stones",
				"description": "Learn how to capture opponent stones",
				"content":     "Stones are captured when they have no liberties. Liberties are empty points directly adjacent to a stone.",
				"level":       "beginner",
			},
			{
				"id":          3,
				"title":       "Life and Death",
				"description": "Understanding when groups are alive or dead",
				"content":     "A group is alive if it has two eyes. An eye is an empty space surrounded by your stones that cannot be filled by the opponent.",
				"level":       "intermediate",
			},
		}
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