# GoSim - Interactive Go Learning Simulator

An interactive, browser-based Go (Weiqi/Baduk) simulator designed to teach the ancient game of Go from beginner to expert level through progressive lessons, interactive puzzles, and AI opponents.

## Features

### 🎯 Progressive Learning System
- **Beginner Mode**: Start with 9x9 boards and learn basic rules
- **Intermediate Mode**: Progress to 13x13 boards with tactical puzzles
- **Advanced Mode**: Master 19x19 boards with joseki and pro game analysis

### 🎮 Interactive Gameplay
- Real-time multiplayer via WebSocket
- Multiple AI difficulty levels
- Move suggestions and hints
- Territory visualization
- Game history and variations

### 📚 Comprehensive Tutorials
- Step-by-step rule explanations
- Interactive capture exercises
- Life and death problems
- Opening theory (joseki)
- Endgame techniques

### 🧩 Puzzle System
- Hundreds of graded problems
- Tactical puzzles (ladders, nets, snapbacks)
- Life and death challenges
- Whole-board problems
- Progress tracking

### 📊 Analysis Tools
- Move evaluation
- Territory estimation
- Variation explorer
- Game review mode
- SGF import/export

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Modern web browser (Chrome, Firefox, Safari, Edge)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Prawal-Sharma/GoSim.git
cd GoSim
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Project Structure

```
GoSim/
├── cmd/server/        # Server application
├── pkg/              
│   ├── game/         # Core game logic
│   ├── learning/     # Tutorial and puzzle system
│   └── websocket/    # Real-time communication
├── web/              # Frontend assets
│   ├── js/           # JavaScript files
│   ├── css/          # Stylesheets
│   └── assets/       # Images and resources
├── data/             # Game data
│   ├── puzzles/      # Puzzle database
│   ├── lessons/      # Tutorial content
│   └── joseki/       # Opening patterns
└── test/             # Test files
```

## Learning Path

### 1. Complete Beginner
- What is Go?
- Placing stones
- Basic objective

### 2. Fundamental Rules
- Capturing stones
- Ko rule
- Suicide rule
- Passing and ending the game

### 3. Basic Tactics
- Ladders
- Nets
- Snapbacks
- Basic connections

### 4. Life and Death
- Two eyes principle
- False eyes
- Seki (mutual life)
- Common patterns

### 5. Territory Concepts
- Building territory
- Invading
- Reducing
- Influence vs territory

### 6. Opening Principles
- Corner-side-center
- Basic joseki
- Direction of play
- Whole board thinking

### 7. Middle Game
- Fighting techniques
- Attack and defense
- Thickness usage
- Weak groups

### 8. Endgame
- Counting
- Sente and gote
- Ko threats
- Point values

### 9. Advanced Concepts
- Professional games analysis
- Modern AI strategies
- Advanced joseki
- Positional judgment

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o gosim cmd/server/main.go
```

### Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## Technologies Used
- **Backend**: Go, Gorilla WebSocket
- **Frontend**: HTML5 Canvas, Vanilla JavaScript
- **AI**: Monte Carlo Tree Search (MCTS)
- **Storage**: JSON files, SQLite (for user data)

## Roadmap
- [ ] Basic game implementation
- [ ] AI opponents (multiple difficulties)
- [ ] Tutorial system
- [ ] Puzzle database
- [ ] Multiplayer support
- [ ] User accounts and progress tracking
- [ ] Mobile responsive design
- [ ] Advanced AI analysis
- [ ] Tournament system
- [ ] Social features

## License
MIT License - see LICENSE file for details

## Acknowledgments
- Inspired by traditional Go teaching methods
- AI algorithms based on modern computer Go research
- Puzzle collection from classical Go problems

## Support
For issues, questions, or suggestions, please open an issue on GitHub.

---
Start your journey to Go mastery today! 🎯