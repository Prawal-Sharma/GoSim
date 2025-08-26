package websocket

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Prawal-Sharma/GoSim/pkg/game"
	"github.com/gorilla/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	game   *game.Game
	color  game.Color
	send   chan []byte
	hub    *Hub
	id     string
	roomID string
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]*GameRoom
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

type GameRoom struct {
	ID      string
	Game    *game.Game
	Players map[game.Color]*Client
}

type Message struct {
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
	RoomID  string                 `json:"roomId,omitempty"`
	PlayerID string                `json:"playerId,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*GameRoom),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered: %s", client.id)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered: %s", client.id)
			}

		case message := <-h.broadcast:
			if message.RoomID != "" {
				if room, ok := h.rooms[message.RoomID]; ok {
					for _, player := range room.Players {
						if player != nil {
							data, _ := json.Marshal(message)
							select {
							case player.send <- data:
							default:
								close(player.send)
								delete(h.clients, player)
							}
						}
					}
				}
			} else {
				data, _ := json.Marshal(message)
				for client := range h.clients {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "create_game":
		c.handleCreateGame(msg)
	case "join_game":
		c.handleJoinGame(msg)
	case "make_move":
		c.handleMakeMove(msg)
	case "pass":
		c.handlePass(msg)
	case "resign":
		c.handleResign(msg)
	case "undo":
		c.handleUndo(msg)
	case "get_valid_moves":
		c.handleGetValidMoves(msg)
	}
}

func (c *Client) handleCreateGame(msg Message) {
	boardSize := 19
	if size, ok := msg.Data["boardSize"].(float64); ok {
		boardSize = int(size)
	}

	roomID := generateRoomID()
	gameRoom := &GameRoom{
		ID:      roomID,
		Game:    game.NewGame(boardSize),
		Players: make(map[game.Color]*Client),
	}

	gameRoom.Players[game.Black] = c
	c.game = gameRoom.Game
	c.color = game.Black
	c.roomID = roomID

	c.hub.rooms[roomID] = gameRoom

	response := Message{
		Type: "game_created",
		Data: map[string]interface{}{
			"roomId":    roomID,
			"boardSize": boardSize,
			"color":     "Black",
		},
	}

	data, _ := json.Marshal(response)
	c.send <- data
}

func (c *Client) handleJoinGame(msg Message) {
	roomID, ok := msg.Data["roomId"].(string)
	if !ok {
		c.sendError("Invalid room ID")
		return
	}

	room, exists := c.hub.rooms[roomID]
	if !exists {
		c.sendError("Room not found")
		return
	}

	if room.Players[game.White] != nil {
		c.sendError("Game is full")
		return
	}

	room.Players[game.White] = c
	c.game = room.Game
	c.color = game.White
	c.roomID = roomID

	response := Message{
		Type: "game_joined",
		Data: map[string]interface{}{
			"roomId":    roomID,
			"boardSize": room.Game.Board.Size,
			"color":     "White",
		},
	}

	data, _ := json.Marshal(response)
	c.send <- data

	c.hub.broadcast <- Message{
		Type:   "game_started",
		RoomID: roomID,
		Data: map[string]interface{}{
			"board": c.game.GetBoardState(),
			"info":  c.game.GetGameInfo(),
		},
	}
}

func (c *Client) handleMakeMove(msg Message) {
	if c.game == nil {
		c.sendError("Not in a game")
		return
	}

	x, okX := msg.Data["x"].(float64)
	y, okY := msg.Data["y"].(float64)

	if !okX || !okY {
		c.sendError("Invalid move coordinates")
		return
	}

	point := game.Point{X: int(x), Y: int(y)}

	err := c.game.MakeMove(point, c.color)
	if err != nil {
		c.sendError(err.Error())
		return
	}

	c.hub.broadcast <- Message{
		Type:   "move_made",
		RoomID: c.roomID,
		Data: map[string]interface{}{
			"x":     int(x),
			"y":     int(y),
			"color": c.color.String(),
			"board": c.game.GetBoardState(),
			"info":  c.game.GetGameInfo(),
		},
	}
}

func (c *Client) handlePass(msg Message) {
	if c.game == nil {
		c.sendError("Not in a game")
		return
	}

	err := c.game.Pass(c.color)
	if err != nil {
		c.sendError(err.Error())
		return
	}

	c.hub.broadcast <- Message{
		Type:   "pass",
		RoomID: c.roomID,
		Data: map[string]interface{}{
			"color": c.color.String(),
			"info":  c.game.GetGameInfo(),
		},
	}

	if c.game.IsOver {
		c.hub.broadcast <- Message{
			Type:   "game_over",
			RoomID: c.roomID,
			Data: map[string]interface{}{
				"winner": c.game.Winner,
				"scores": c.game.CalculateScore(),
			},
		}
	}
}

func (c *Client) handleResign(msg Message) {
	if c.game == nil {
		c.sendError("Not in a game")
		return
	}

	c.game.Resign(c.color)

	c.hub.broadcast <- Message{
		Type:   "resign",
		RoomID: c.roomID,
		Data: map[string]interface{}{
			"color":  c.color.String(),
			"winner": game.OpponentColor(c.color).String(),
		},
	}
}

func (c *Client) handleUndo(msg Message) {
	if c.game == nil {
		c.sendError("Not in a game")
		return
	}

	success := c.game.Undo()
	if !success {
		c.sendError("Cannot undo")
		return
	}

	c.hub.broadcast <- Message{
		Type:   "undo",
		RoomID: c.roomID,
		Data: map[string]interface{}{
			"board": c.game.GetBoardState(),
			"info":  c.game.GetGameInfo(),
		},
	}
}

func (c *Client) handleGetValidMoves(msg Message) {
	if c.game == nil {
		c.sendError("Not in a game")
		return
	}

	validMoves := c.game.GetValidMoves(c.color)
	movesData := []map[string]int{}

	for _, move := range validMoves {
		movesData = append(movesData, map[string]int{
			"x": move.X,
			"y": move.Y,
		})
	}

	response := Message{
		Type: "valid_moves",
		Data: map[string]interface{}{
			"moves": movesData,
		},
	}

	data, _ := json.Marshal(response)
	c.send <- data
}

func (c *Client) sendError(errorMsg string) {
	response := Message{
		Type: "error",
		Data: map[string]interface{}{
			"message": errorMsg,
		},
	}

	data, _ := json.Marshal(response)
	c.send <- data
}

func generateRoomID() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	roomID := ""
	for i := 0; i < 6; i++ {
		roomID += string(letters[rand.Intn(len(letters))])
	}
	return roomID
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		id:   generateRoomID(),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}