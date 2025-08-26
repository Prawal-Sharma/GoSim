# GoSim Architecture

## Overview

GoSim follows a client-server architecture with clear separation between game logic (backend) and presentation (frontend). The system is designed for modularity, scalability, and ease of maintenance.

## System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Web Browser                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │                Frontend (JavaScript)             │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────┐   │   │
│  │  │  Board   │ │   Game   │ │   Learning   │   │   │
│  │  │  Canvas  │ │  Logic   │ │    Module    │   │   │
│  │  └──────────┘ └──────────┘ └──────────────┘   │   │
│  │         │           │              │            │   │
│  │         └───────────┴──────────────┘            │   │
│  │                     │                           │   │
│  │            ┌────────▼────────┐                 │   │
│  │            │  WebSocket      │                 │   │
│  │            │  Connection     │                 │   │
│  │            └────────┬────────┘                 │   │
│  └─────────────────────┼───────────────────────────┘   │
└───────────────────────┼─────────────────────────────────┘
                        │
                    WebSocket
                        │
┌───────────────────────▼─────────────────────────────────┐
│                   Go Server (Port 8080)                 │
│  ┌─────────────────────────────────────────────────┐   │
│  │              HTTP Router (Chi)                   │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────┐   │   │
│  │  │   Static │ │   REST   │ │   WebSocket  │   │   │
│  │  │   Files  │ │    API   │ │    Handler   │   │   │
│  │  └──────────┘ └──────────┘ └──────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Core Game Engine                    │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────┐   │   │
│  │  │  Board   │ │   Rules  │ │      AI      │   │   │
│  │  │  State   │ │  Engine  │ │   Algorithms │   │   │
│  │  └──────────┘ └──────────┘ └──────────────┘   │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────┐   │   │
│  │  │ Scoring  │ │  Game    │ │   Learning   │   │   │
│  │  │  System  │ │  Session │ │   Content    │   │   │
│  │  └──────────┘ └──────────┘ └──────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Data Layer                          │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────┐   │   │
│  │  │  Puzzles │ │  Lessons │ │    Game      │   │   │
│  │  │   JSON   │ │   JSON   │ │    State     │   │   │
│  │  └──────────┘ └──────────┘ └──────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## Component Details

### Frontend Components

#### Board Canvas (`web/js/board.js`)
- **Responsibility**: Visual representation of the game board
- **Key Features**:
  - Stone rendering with gradients
  - Grid and star point drawing
  - Ghost stone for hover effects
  - Territory visualization
  - Last move marker
- **Dependencies**: None (vanilla JS)

#### Game Controller (`web/js/game.js`)
- **Responsibility**: Game state management and user interaction
- **Key Features**:
  - Game mode selection
  - Move handling
  - Turn management
  - Timer control
  - History tracking
- **Dependencies**: Board, WebSocket, Learning modules

#### WebSocket Client (`web/js/websocket.js`)
- **Responsibility**: Real-time server communication
- **Key Features**:
  - Auto-reconnection
  - Message queuing
  - State synchronization
  - Error handling
- **Dependencies**: Game controller

#### Learning Module (`web/js/learning.js`)
- **Responsibility**: Educational features
- **Key Features**:
  - Tutorial management
  - Exercise validation
  - Progress tracking
  - Hint system
- **Dependencies**: Game controller, Board

### Backend Components

#### HTTP Server (`cmd/server/main.go`)
- **Responsibility**: Request routing and server initialization
- **Key Features**:
  - Static file serving
  - API endpoint routing
  - WebSocket upgrade
  - CORS handling
- **Dependencies**: All backend packages

#### Game Package (`pkg/game/`)

##### Board (`board.go`)
- **Responsibility**: Board state representation
- **Data Structures**:
  ```go
  type Board struct {
      Size     int
      Grid     [][]Color
      LastMove *Point
      Captures map[Color]int
      History  []BoardState
      KoPoint  *Point
  }
  ```
- **Key Methods**:
  - `GetGroup()`: Find connected stones
  - `GetLiberties()`: Count group liberties
  - `CaptureDeadGroups()`: Remove captured stones
  - `CountTerritory()`: Calculate controlled area

##### Rules (`rules.go`)
- **Responsibility**: Game rule enforcement
- **Key Features**:
  - Move validation
  - Ko detection
  - Suicide prevention
  - Pass handling
  - Game end detection
- **Error Types**:
  - `ErrInvalidMove`
  - `ErrKoViolation`
  - `ErrSuicideMove`
  - `ErrPositionOccupied`

##### AI (`ai.go`)
- **Responsibility**: Computer opponent logic
- **Algorithms**:
  - **Random**: Random legal move selection
  - **Easy**: Basic position evaluation
  - **Medium**: Advanced evaluation with influence
  - **Hard**: Minimax with alpha-beta pruning
- **Evaluation Factors**:
  - Captures
  - Liberties
  - Territory
  - Influence
  - Eye formation
  - Group connections

##### Scoring (`scoring.go`)
- **Responsibility**: Game scoring calculation
- **Scoring Methods**:
  - Chinese: Territory + Stones on board
  - Japanese: Territory + Captures
- **Features**:
  - Dead stone detection
  - Territory marking
  - Komi handling
  - SGF generation

#### WebSocket Package (`pkg/websocket/`)

##### Handler (`handler.go`)
- **Responsibility**: WebSocket connection management
- **Components**:
  - `Hub`: Central message router
  - `Client`: Individual connection handler
  - `GameRoom`: Game session container
- **Message Types**:
  - Game creation/joining
  - Move transmission
  - State synchronization
  - Error handling

### Data Flow

#### Move Execution Flow
```
1. User clicks board
   ↓
2. Frontend validates basic rules
   ↓
3. Send move via WebSocket/API
   ↓
4. Server validates move
   ↓
5. Update game state
   ↓
6. Calculate captures
   ↓
7. Check for ko
   ↓
8. Broadcast update
   ↓
9. Update all clients
   ↓
10. Render new board state
```

#### AI Move Flow
```
1. Player makes move
   ↓
2. Frontend requests AI move
   ↓
3. Server evaluates position
   ↓
4. AI algorithm calculates
   ↓
5. Best move selected
   ↓
6. Move validated
   ↓
7. State updated
   ↓
8. Response sent to client
```

### State Management

#### Client State
```javascript
{
  board: [][], // 2D array of stones
  currentTurn: "black",
  gameMode: "ai",
  gameStarted: true,
  playerColor: "black",
  moveHistory: [],
  captures: {black: 0, white: 0}
}
```

#### Server State
```go
type GameState struct {
    Board       *Board
    CurrentTurn Color
    Passed      map[Color]bool
    IsOver      bool
    Winner      *Color
    MoveCount   int
}
```

### Communication Protocol

#### REST API
- `GET /`: Serve index.html
- `GET /api/health`: Health check
- `POST /api/ai-move`: Get AI move
- `GET /api/puzzles`: Get puzzle list
- `GET /api/lessons`: Get lesson list

#### WebSocket Messages
```json
// Client → Server
{
  "type": "make_move",
  "data": {"x": 3, "y": 3}
}

// Server → Client
{
  "type": "move_made",
  "data": {
    "x": 3,
    "y": 3,
    "board": [...],
    "info": {...}
  }
}
```

## Design Patterns

### Observer Pattern
- WebSocket Hub broadcasts to all clients
- Clients subscribe to game events

### Strategy Pattern
- Different AI algorithms implement same interface
- Scoring methods are interchangeable

### Factory Pattern
- Board creation with different sizes
- Game mode initialization

### Command Pattern
- Move execution and validation
- Undo/redo functionality

## Security Considerations

### Input Validation
- All moves validated server-side
- Bounds checking on coordinates
- Type validation on API inputs

### WebSocket Security
- Origin checking (configurable)
- Message size limits
- Connection rate limiting (planned)

### Data Protection
- No sensitive data stored
- Local storage for preferences only
- No user authentication (yet)

## Performance Optimizations

### Frontend
- Canvas rendering optimizations
- Minimal DOM manipulation
- RequestAnimationFrame for animations
- Event delegation

### Backend
- Efficient board representation
- Cached territory calculations
- Goroutine pool for WebSocket
- Minimal allocations in hot paths

## Scalability

### Current Limitations
- Single server instance
- In-memory game state
- No persistence layer
- No load balancing

### Future Improvements
- Redis for game state
- Horizontal scaling with message queue
- Database for user accounts
- CDN for static assets
- Microservices architecture

## Testing Strategy

### Unit Tests
- Game logic validation
- Rule enforcement
- AI move generation
- Scoring accuracy

### Integration Tests
- API endpoint testing
- WebSocket communication
- Full game scenarios

### Frontend Tests (Planned)
- Component testing
- E2E with Selenium
- Visual regression tests

## Deployment

### Development
```bash
go run cmd/server/main.go
```

### Production
```bash
go build -o gosim cmd/server/main.go
./gosim
```

### Docker (Planned)
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o gosim cmd/server/main.go
CMD ["./gosim"]
```

## Monitoring (Future)

### Metrics to Track
- Active connections
- Game completion rate
- AI response time
- Error rates
- Memory usage

### Logging
- Structured logging with levels
- Error tracking
- Performance profiling
- User analytics (privacy-respecting)

## Dependencies

### Go Dependencies
- `gorilla/websocket`: WebSocket support
- `go-chi/chi`: HTTP routing
- `go-chi/cors`: CORS handling

### Frontend Dependencies
- None (vanilla JavaScript)

### Development Dependencies
- `air`: Hot reload
- `golangci-lint`: Code quality
- `delve`: Debugging