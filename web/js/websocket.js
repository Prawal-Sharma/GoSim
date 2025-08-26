class WebSocketConnection {
    constructor(game) {
        this.game = game;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        
        this.connect();
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => this.onOpen();
            this.ws.onmessage = (event) => this.onMessage(event);
            this.ws.onclose = () => this.onClose();
            this.ws.onerror = (error) => this.onError(error);
        } catch (error) {
            console.error('WebSocket connection error:', error);
            this.reconnect();
        }
    }

    onOpen() {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        this.game.updateStatus('Connected to server');
    }

    onMessage(event) {
        try {
            const data = JSON.parse(event.data);
            this.handleMessage(data);
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    }

    onClose() {
        console.log('WebSocket disconnected');
        this.game.updateStatus('Disconnected from server');
        this.reconnect();
    }

    onError(error) {
        console.error('WebSocket error:', error);
    }

    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
            
            console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts}) in ${delay}ms...`);
            this.game.updateStatus(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            
            setTimeout(() => {
                this.connect();
            }, delay);
        } else {
            this.game.updateStatus('Connection lost. Please refresh the page.');
            this.game.showModal('Connection Error', 'Unable to connect to the server. Please refresh the page and try again.');
        }
    }

    handleMessage(data) {
        console.log('Received message:', data);
        
        switch (data.type) {
            case 'game_created':
                this.handleGameCreated(data);
                break;
            case 'game_joined':
                this.handleGameJoined(data);
                break;
            case 'game_started':
                this.handleGameStarted(data);
                break;
            case 'move_made':
                this.handleMoveMade(data);
                break;
            case 'pass':
                this.handlePass(data);
                break;
            case 'resign':
                this.handleResign(data);
                break;
            case 'game_over':
                this.handleGameOver(data);
                break;
            case 'undo':
                this.handleUndo(data);
                break;
            case 'valid_moves':
                this.handleValidMoves(data);
                break;
            case 'error':
                this.handleError(data);
                break;
        }
    }

    handleGameCreated(data) {
        this.game.roomId = data.data.roomId;
        this.game.playerColor = 'black';
        document.getElementById('room-info').innerHTML = `
            <strong>Room ID: ${this.game.roomId}</strong><br>
            You are playing: Black<br>
            Waiting for opponent...
        `;
        document.getElementById('player-color').textContent = 'Playing as Black';
    }

    handleGameJoined(data) {
        this.game.roomId = data.data.roomId;
        this.game.playerColor = 'white';
        document.getElementById('room-info').innerHTML = `
            <strong>Room ID: ${this.game.roomId}</strong><br>
            You are playing: White<br>
            Game starting...
        `;
        document.getElementById('player-color').textContent = 'Playing as White';
    }

    handleGameStarted(data) {
        this.game.gameStarted = true;
        this.game.currentTurn = 'black';
        
        document.getElementById('menu-screen').style.display = 'none';
        document.getElementById('game-controls').style.display = 'block';
        
        if (data.data.board) {
            this.game.board.updateBoard(data.data.board);
        }
        
        this.game.startTimer();
        this.game.updateStatus('Game started!');
        this.game.updateTurnIndicator();
    }

    handleMoveMade(data) {
        const boardData = data.data.board;
        const convertedBoard = boardData.map(row => 
            row.map(cell => cell)
        );
        
        this.game.board.updateBoard(convertedBoard);
        this.game.board.setLastMove(data.data.x, data.data.y);
        
        const moveColor = data.data.color;
        this.game.moveHistory.push({
            move: this.game.moveHistory.length + 1,
            color: moveColor,
            position: this.game.getPositionNotation(data.data.x, data.data.y)
        });
        
        this.game.updateMoveHistory();
        
        this.game.currentTurn = data.data.info.currentTurn.toLowerCase();
        this.game.updateTurnIndicator();
        
        document.getElementById('move-count').textContent = data.data.info.moveCount;
        document.getElementById('black-captures').textContent = data.data.info.blackCaptures;
        document.getElementById('white-captures').textContent = data.data.info.whiteCaptures;
    }

    handlePass(data) {
        this.game.moveHistory.push({
            move: this.game.moveHistory.length + 1,
            color: data.data.color,
            position: 'Pass'
        });
        
        this.game.updateMoveHistory();
        this.game.currentTurn = data.data.info.currentTurn.toLowerCase();
        this.game.updateTurnIndicator();
    }

    handleResign(data) {
        this.game.gameStarted = false;
        this.game.showModal('Game Over', `${data.data.color} resigned. ${data.data.winner} wins!`);
    }

    handleGameOver(data) {
        this.game.gameStarted = false;
        
        let message = '';
        if (data.data.winner) {
            message = `${data.data.winner} wins!\n`;
        } else {
            message = 'Game ended in a draw.\n';
        }
        
        if (data.data.scores) {
            message += `\nFinal Score:\n`;
            message += `Black: ${data.data.scores.Black}\n`;
            message += `White: ${data.data.scores.White}`;
        }
        
        this.game.showModal('Game Over', message);
    }

    handleUndo(data) {
        const convertedBoard = data.data.board.map(row => 
            row.map(cell => cell)
        );
        
        this.game.board.updateBoard(convertedBoard);
        
        if (this.game.moveHistory.length > 0) {
            this.game.moveHistory.pop();
            this.game.updateMoveHistory();
        }
        
        this.game.currentTurn = data.data.info.currentTurn.toLowerCase();
        this.game.updateTurnIndicator();
    }

    handleValidMoves(data) {
        const moves = data.data.moves;
        this.game.board.setValidMoves(moves);
    }

    handleError(data) {
        console.error('Server error:', data.data.message);
        this.game.updateStatus(`Error: ${data.data.message}`);
    }

    send(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not connected');
        }
    }

    createGame(boardSize) {
        this.send({
            type: 'create_game',
            data: {
                boardSize: boardSize
            }
        });
    }

    joinGame(roomId) {
        this.send({
            type: 'join_game',
            data: {
                roomId: roomId
            }
        });
    }

    sendMove(x, y) {
        this.send({
            type: 'make_move',
            data: {
                x: x,
                y: y
            }
        });
    }

    sendPass() {
        this.send({
            type: 'pass',
            data: {}
        });
    }

    sendResign() {
        this.send({
            type: 'resign',
            data: {}
        });
    }

    sendUndo() {
        this.send({
            type: 'undo',
            data: {}
        });
    }

    getValidMoves() {
        this.send({
            type: 'get_valid_moves',
            data: {}
        });
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}