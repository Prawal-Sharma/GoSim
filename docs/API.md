# GoSim API Documentation

## Overview
GoSim provides both REST API endpoints and WebSocket connections for real-time gameplay. The server runs on port 8080 by default.

## Base URL
```
http://localhost:8080
```

## REST API Endpoints

### 1. Health Check
**GET** `/api/health`

Check if the server is running and healthy.

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. Get AI Move
**POST** `/api/ai-move`

Request an AI move for the current board state.

**Request Body:**
```json
{
  "board": [[0,0,1], [2,0,0], [0,1,2]],  // 2D array representing board state
  "boardSize": 9,                         // Board size (9, 13, or 19)
  "color": "Black",                       // "Black" or "White"
  "difficulty": "easy"                    // "random", "easy", "medium", or "hard"
}
```

**Response:**
```json
{
  "x": 4,
  "y": 5
}
```

Or for pass:
```json
{
  "pass": true
}
```

### 3. Get Puzzles
**GET** `/api/puzzles`

Retrieve all available puzzles.

**Response:**
```json
[
  {
    "id": 1,
    "title": "Basic Capture",
    "description": "Capture the white stone",
    "difficulty": "beginner",
    "board": [[0,0,1], [1,2,1], [0,1,0]],
    "solution": {
      "moves": [{"x": 1, "y": 0, "color": 1}],
      "explanation": "Place black stone to capture"
    },
    "hints": ["Look for the last liberty"]
  }
]
```

### 4. Get Lessons
**GET** `/api/lessons`

Retrieve all available lessons.

**Response:**
```json
[
  {
    "id": 1,
    "title": "Introduction to Go",
    "level": "beginner",
    "description": "Learn the basics of Go",
    "content": "Go is an ancient board game...",
    "exercises": [
      {
        "type": "quiz",
        "question": "What is the objective of Go?",
        "answer": "Control territory"
      }
    ]
  }
]
```

## WebSocket API

### Connection
**WebSocket URL:** `ws://localhost:8080/ws`

### Message Format
All WebSocket messages use JSON format:

```json
{
  "type": "message_type",
  "data": {
    // message-specific data
  }
}
```

### Client to Server Messages

#### 1. Create Game
```json
{
  "type": "create_game",
  "data": {
    "boardSize": 19
  }
}
```

#### 2. Join Game
```json
{
  "type": "join_game",
  "data": {
    "roomId": "ABC123"
  }
}
```

#### 3. Make Move
```json
{
  "type": "make_move",
  "data": {
    "x": 3,
    "y": 3
  }
}
```

#### 4. Pass
```json
{
  "type": "pass",
  "data": {}
}
```

#### 5. Resign
```json
{
  "type": "resign",
  "data": {}
}
```

#### 6. Undo
```json
{
  "type": "undo",
  "data": {}
}
```

#### 7. Get Valid Moves
```json
{
  "type": "get_valid_moves",
  "data": {}
}
```

### Server to Client Messages

#### 1. Game Created
```json
{
  "type": "game_created",
  "data": {
    "roomId": "ABC123",
    "boardSize": 19,
    "color": "Black"
  }
}
```

#### 2. Game Joined
```json
{
  "type": "game_joined",
  "data": {
    "roomId": "ABC123",
    "boardSize": 19,
    "color": "White"
  }
}
```

#### 3. Game Started
```json
{
  "type": "game_started",
  "data": {
    "board": [[0,0,0]...],
    "info": {
      "boardSize": 19,
      "currentTurn": "Black",
      "moveCount": 0,
      "blackCaptures": 0,
      "whiteCaptures": 0
    }
  }
}
```

#### 4. Move Made
```json
{
  "type": "move_made",
  "data": {
    "x": 3,
    "y": 3,
    "color": "Black",
    "board": [[0,0,1]...],
    "info": {
      "currentTurn": "White",
      "moveCount": 1,
      "blackCaptures": 0,
      "whiteCaptures": 0
    }
  }
}
```

#### 5. Game Over
```json
{
  "type": "game_over",
  "data": {
    "winner": "Black",
    "scores": {
      "Black": 45,
      "White": 40
    }
  }
}
```

#### 6. Error
```json
{
  "type": "error",
  "data": {
    "message": "Invalid move"
  }
}
```

#### 7. Valid Moves Response
```json
{
  "type": "valid_moves",
  "data": {
    "moves": [
      {"x": 0, "y": 0},
      {"x": 0, "y": 1},
      {"x": 1, "y": 0}
    ]
  }
}
```

## Board State Representation

The board is represented as a 2D array where:
- `0` = Empty intersection
- `1` = Black stone
- `2` = White stone

Example 3x3 board:
```json
[
  [0, 1, 0],  // Top row
  [2, 0, 1],  // Middle row
  [0, 2, 0]   // Bottom row
]
```

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Invalid move | The move violates game rules |
| 400 | Position already occupied | Stone already exists at this position |
| 400 | Suicide move not allowed | Move would result in self-capture |
| 400 | Ko rule violation | Move would recreate previous board state |
| 400 | Not your turn | Player attempted move out of turn |
| 404 | Room not found | Game room doesn't exist |
| 400 | Game is full | Room already has two players |

## Rate Limiting

Currently no rate limiting is implemented, but for production use, consider:
- API calls: 100 requests per minute per IP
- WebSocket messages: 10 messages per second per connection

## CORS Policy

CORS is enabled for all origins in development. For production, update the allowed origins in `cmd/server/main.go`:

```go
cors.Handler(cors.Options{
    AllowedOrigins: []string{"https://yourdomain.com"},
    // ... other options
})
```