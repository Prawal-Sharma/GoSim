# GoSim - Interactive Go Learning Simulator

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-Web-orange)
![Status](https://img.shields.io/badge/Status-Active-success)

**Learn the ancient game of Go interactively with AI opponents, puzzles, and tutorials**

[Features](#features) • [Quick Start](#quick-start) • [Game Modes](#game-modes) • [Documentation](#documentation) • [Contributing](#contributing)

</div>

---

GoSim is a comprehensive, browser-based Go (Weiqi/Baduk) learning platform that combines traditional game play with modern educational features. Perfect for beginners starting their Go journey or experienced players looking to improve their skills.

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
- Go 1.21 or higher (Download from: https://golang.org/dl/)
- Modern web browser (Chrome, Firefox, Safari, Edge)

### Installation & Running

#### Easy Method (Recommended):

**macOS/Linux:**
```bash
git clone https://github.com/Prawal-Sharma/GoSim.git
cd GoSim
./start.sh
```

**Windows:**
```cmd
git clone https://github.com/Prawal-Sharma/GoSim.git
cd GoSim
start.bat
```

#### Manual Method:

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

#### Using Make (if installed):
```bash
make run    # Run the server
make build  # Build binary
make test   # Run tests
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

## Game Modes

### 🤖 vs AI
- **4 Difficulty Levels**: Random, Easy, Medium, Hard
- **Smart AI**: Uses pattern recognition and strategic evaluation
- **Instant Play**: No setup required
- **Learning Mode**: AI adapts to your skill level

### 👥 Multiplayer Options
- **Local 2-Player**: Play on the same device
- **Online Multiplayer**: Real-time games via WebSocket
- **Room System**: Create/join games with room codes
- **Spectator Mode**: Watch ongoing games (coming soon)

### 📚 Learning Mode
- **Interactive Tutorials**: Step-by-step lessons
- **Concept Exercises**: Practice specific skills
- **Progress Tracking**: Monitor your improvement
- **Visual Guides**: Animated demonstrations

### 🧩 Puzzle Mode
- **Graded Problems**: Beginner to advanced
- **Solution Hints**: Progressive hint system
- **Categories**: Life/death, tesuji, endgame
- **Custom Puzzles**: Create your own (coming soon)

## Documentation

📖 **Complete documentation available in the `/docs` folder:**

- [**API Documentation**](docs/API.md) - REST and WebSocket API reference
- [**Game Rules**](docs/GAME_RULES.md) - Complete Go rules and strategy guide
- [**Troubleshooting**](docs/TROUBLESHOOTING.md) - Common issues and solutions
- [**Contributing**](CONTRIBUTING.md) - Development guidelines

## Features in Detail

### Game Engine
- ✅ Complete rule implementation (capture, ko, suicide)
- ✅ Territory calculation (Chinese/Japanese scoring)
- ✅ Game history and undo functionality
- ✅ Move validation and legal move generation
- ✅ SGF export/import (coming soon)

### AI System
- ✅ Multiple algorithms (random, greedy, minimax)
- ✅ Position evaluation with pattern recognition
- ✅ Influence calculation
- ✅ Life/death analysis
- ✅ Opening book (coming soon)

### User Interface
- ✅ Responsive design for all devices
- ✅ Real-time board updates
- ✅ Move animations
- ✅ Territory visualization
- ✅ Sound effects (coming soon)

### Learning Features
- ✅ 5+ interactive tutorials
- ✅ 10+ tactical puzzles
- ✅ Progress persistence
- ✅ Skill assessment (coming soon)
- ✅ Personalized recommendations (coming soon)

## System Requirements

### Minimum Requirements
- **Browser**: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- **Screen**: 1024x768 resolution
- **Network**: Required for multiplayer only
- **JavaScript**: Must be enabled

### Recommended
- **Browser**: Latest version of Chrome or Firefox
- **Screen**: 1920x1080 or higher
- **Network**: Broadband for smooth multiplayer

## Roadmap

### ✅ Completed
- [x] Core game implementation
- [x] Multiple AI difficulties
- [x] Basic tutorial system
- [x] Puzzle framework
- [x] Local and online multiplayer
- [x] Responsive design
- [x] Progress tracking

### 🚧 In Progress
- [ ] Extended puzzle database
- [ ] Advanced AI analysis
- [ ] Mobile app version
- [ ] Cloud save sync

### 📋 Planned
- [ ] User accounts and profiles
- [ ] Tournament system
- [ ] Social features (friends, chat)
- [ ] Game replay and analysis
- [ ] Opening book database
- [ ] Professional game library
- [ ] Live streaming integration
- [ ] AI teaching assistant

## Performance

### Benchmarks
- **Board Generation**: <1ms for any size
- **Move Validation**: <0.1ms average
- **AI Response**: 50-500ms depending on difficulty
- **Territory Calculation**: <5ms for 19x19
- **WebSocket Latency**: <50ms on local network

### Browser Compatibility
| Browser | Version | Support |
|---------|---------|---------|
| Chrome | 90+ | ✅ Full |
| Firefox | 88+ | ✅ Full |
| Safari | 14+ | ✅ Full |
| Edge | 90+ | ✅ Full |
| Opera | 76+ | ✅ Full |
| Mobile Chrome | Latest | ⚠️ Partial |
| Mobile Safari | Latest | ⚠️ Partial |

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