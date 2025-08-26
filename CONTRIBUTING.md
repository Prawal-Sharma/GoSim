# Contributing to GoSim

Thank you for your interest in contributing to GoSim! This document provides guidelines and instructions for contributing to the project.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Adding Features](#adding-features)

## Code of Conduct

### Our Pledge
We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity, level of experience, nationality, personal appearance, race, religion, or sexual identity.

### Expected Behavior
- Be respectful and inclusive
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards others

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git
- A modern web browser
- Basic understanding of Go and JavaScript

### Fork and Clone
1. Fork the repository on GitHub
2. Clone your fork:
```bash
git clone https://github.com/YOUR-USERNAME/GoSim.git
cd GoSim
```

3. Add upstream remote:
```bash
git remote add upstream https://github.com/Prawal-Sharma/GoSim.git
```

## Development Setup

### Install Dependencies
```bash
# Go dependencies
go mod download

# Install development tools (optional)
go install github.com/cosmtrek/air@latest  # Hot reload
go install golang.org/x/tools/cmd/goimports@latest  # Import formatting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
```

### Running the Development Server
```bash
# With hot reload (if air is installed)
air

# Standard way
go run cmd/server/main.go

# Using make
make run
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/game

# Using make
make test
```

## Project Structure

```
GoSim/
├── cmd/server/          # Application entry point
│   └── main.go         # Server initialization
├── pkg/                 # Core packages
│   ├── game/           # Game logic
│   │   ├── board.go    # Board representation
│   │   ├── rules.go    # Game rules
│   │   ├── ai.go       # AI implementation
│   │   └── scoring.go  # Scoring system
│   ├── websocket/      # WebSocket handling
│   └── learning/       # Tutorial system (future)
├── web/                # Frontend assets
│   ├── index.html      # Main HTML
│   ├── css/           # Stylesheets
│   ├── js/            # JavaScript files
│   └── assets/        # Images and resources
├── data/              # Game data
│   ├── puzzles/       # Puzzle definitions
│   └── lessons/       # Tutorial content
├── docs/              # Documentation
└── test/              # Test files
```

## How to Contribute

### Reporting Bugs
1. Check if the issue already exists
2. Create a new issue with:
   - Clear, descriptive title
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - System information
   - Screenshots if applicable

### Suggesting Features
1. Check existing issues and discussions
2. Create a feature request with:
   - Use case description
   - Proposed solution
   - Alternative solutions considered
   - Additional context

### Code Contributions

#### 1. Find an Issue
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to claim it
- Ask questions if requirements are unclear

#### 2. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-description
```

#### 3. Make Changes
- Write clean, documented code
- Follow the coding standards
- Add tests for new features
- Update documentation

#### 4. Commit Changes
```bash
# Format your code
go fmt ./...

# Check for issues
go vet ./...

# Commit with descriptive message
git commit -m "feat: add new AI difficulty level

- Implemented expert level AI using advanced MCTS
- Added configuration option in UI
- Updated documentation"
```

Commit message format:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Formatting
- `refactor:` Code restructuring
- `test:` Test additions
- `chore:` Maintenance

## Coding Standards

### Go Code
```go
// Package comment describes the package purpose
package game

// Exported types and functions need comments
// Board represents the Go game board
type Board struct {
    Size int
    Grid [][]Color
}

// NewBoard creates a new game board
func NewBoard(size int) *Board {
    // Implementation
}

// Use meaningful variable names
validMoves := g.GetValidMoves(color)

// Handle errors properly
if err != nil {
    return fmt.Errorf("failed to make move: %w", err)
}
```

### JavaScript Code
```javascript
// Use ES6+ features
class GoBoard {
    constructor(canvasId, size = 19) {
        this.canvas = document.getElementById(canvasId);
        this.size = size;
    }
    
    // Document complex logic
    /**
     * Calculates territory for the given board state
     * @param {Array} board - 2D array of stones
     * @returns {Object} Territory count by color
     */
    calculateTerritory(board) {
        // Implementation
    }
}

// Use const/let, not var
const BOARD_SIZE = 19;
let currentPlayer = 'black';

// Handle async properly
async function fetchPuzzles() {
    try {
        const response = await fetch('/api/puzzles');
        return await response.json();
    } catch (error) {
        console.error('Failed to fetch puzzles:', error);
    }
}
```

### HTML/CSS
- Use semantic HTML5 elements
- Keep CSS organized and commented
- Use CSS classes, not inline styles
- Maintain responsive design

## Testing

### Writing Tests
```go
func TestBoardCreation(t *testing.T) {
    board := NewBoard(19)
    
    if board.Size != 19 {
        t.Errorf("Expected size 19, got %d", board.Size)
    }
}

// Table-driven tests for multiple cases
func TestCapture(t *testing.T) {
    tests := []struct {
        name     string
        moves    []Move
        expected int
    }{
        {"single stone", []Move{{3,3,Black}}, 1},
        {"group capture", []Move{{3,3,Black}, {3,4,Black}}, 2},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Coverage Requirements
- Aim for >80% test coverage for new code
- Critical game logic should have >90% coverage
- Run coverage: `go test -cover ./...`

## Pull Request Process

### Before Submitting
1. **Update from upstream**
```bash
git fetch upstream
git rebase upstream/main
```

2. **Run tests**
```bash
make test
```

3. **Format code**
```bash
go fmt ./...
```

4. **Update documentation**
- Add/update code comments
- Update README if needed
- Add to docs/ if significant feature

### PR Guidelines
1. **Title**: Clear, descriptive title
2. **Description**: 
   - What changes were made
   - Why were they made
   - How to test them
3. **Screenshots**: For UI changes
4. **Link issues**: Reference related issues
5. **Small PRs**: Keep changes focused

### Review Process
1. Automated checks must pass
2. Code review by maintainer
3. Address review feedback
4. Approval and merge

## Adding Features

### New Game Mode
1. Add mode logic in `pkg/game/`
2. Create UI in `web/js/`
3. Add WebSocket handlers if needed
4. Write tests
5. Update documentation

### New AI Algorithm
1. Implement in `pkg/game/ai.go`
2. Add difficulty option
3. Test against existing AIs
4. Document algorithm approach

### New Puzzle Type
1. Define format in `data/puzzles/`
2. Add loader in server
3. Create UI handler
4. Add solution validator

### UI Improvements
1. Maintain responsive design
2. Test on multiple browsers
3. Keep accessibility in mind
4. Follow existing style patterns

## Development Tips

### Debugging Go
```go
// Use log package for debugging
log.Printf("Board state: %+v", board)

// Use debugger
// Install delve: go install github.com/go-delve/delve/cmd/dlv@latest
// Debug: dlv debug cmd/server/main.go
```

### Debugging JavaScript
```javascript
// Use browser DevTools
console.log('Game state:', this.gameState);

// Use debugger statement
debugger; // Pauses execution

// Network tab for API calls
// Console for errors
// Sources for breakpoints
```

### Performance
- Profile Go code: `go test -bench=. -cpuprofile=cpu.prof`
- Use browser Performance tab for JS
- Optimize hot paths
- Cache where appropriate

## Questions?

If you have questions:
1. Check existing documentation
2. Search closed issues
3. Ask in a new issue
4. Contact maintainers

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Given credit in commit messages

Thank you for contributing to GoSim!