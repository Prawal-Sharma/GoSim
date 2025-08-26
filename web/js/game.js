class GoGame {
    constructor() {
        this.board = new GoBoard('go-board', 13);
        this.gameMode = null;
        this.playerColor = 'black';
        this.currentTurn = 'black';
        this.moveHistory = [];
        this.gameStarted = false;
        this.wsConnection = null;
        this.aiDifficulty = 'easy';
        this.roomId = null;
        this.gameTimer = null;
        this.startTime = null;
        
        this.setupEventListeners();
        this.initializeGame();
    }

    setupEventListeners() {
        this.board.onMove = (x, y) => this.handleMove(x, y);
        
        document.querySelectorAll('.size-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.changeBoardSize(e));
        });
        
        document.getElementById('single-player-btn').addEventListener('click', () => this.startSinglePlayer());
        document.getElementById('local-player-btn').addEventListener('click', () => this.startLocalPlayer());
        document.getElementById('multiplayer-btn').addEventListener('click', () => this.startMultiplayer());
        document.getElementById('learn-btn').addEventListener('click', () => this.startLearningMode());
        document.getElementById('puzzle-btn').addEventListener('click', () => this.startPuzzleMode());
        
        document.getElementById('pass-btn').addEventListener('click', () => this.pass());
        document.getElementById('resign-btn').addEventListener('click', () => this.resign());
        document.getElementById('undo-btn').addEventListener('click', () => this.undo());
        document.getElementById('new-game-btn').addEventListener('click', () => this.newGame());
        
        document.getElementById('create-room-btn').addEventListener('click', () => this.createRoom());
        document.getElementById('join-room-btn').addEventListener('click', () => this.joinRoom());
        
        document.getElementById('difficulty-select').addEventListener('change', (e) => {
            this.aiDifficulty = e.target.value;
        });
        
        document.querySelector('.close').addEventListener('click', () => this.closeModal());
        document.getElementById('modal-ok-btn').addEventListener('click', () => this.closeModal());
    }

    initializeGame() {
        this.updateStatus('Welcome! Select a game mode to start playing.');
        this.updateTurnIndicator();
    }

    changeBoardSize(event) {
        if (this.gameStarted) {
            this.showModal('Game in Progress', 'Please finish the current game before changing board size.');
            return;
        }
        
        const size = parseInt(event.target.dataset.size);
        document.querySelectorAll('.size-btn').forEach(btn => btn.classList.remove('active'));
        event.target.classList.add('active');
        
        this.board.reset(size);
    }

    startSinglePlayer() {
        this.gameMode = 'ai';
        this.gameStarted = true;
        this.playerColor = 'black';
        this.currentTurn = 'black';
        this.moveHistory = [];
        
        document.getElementById('menu-screen').style.display = 'none';
        document.getElementById('game-controls').style.display = 'block';
        document.getElementById('ai-difficulty').style.display = 'block';
        
        this.board.reset(this.board.size);
        this.startTimer();
        this.updateStatus('Game started! You are playing Black.');
        this.updateTurnIndicator();
    }

    startLocalPlayer() {
        this.gameMode = 'local';
        this.gameStarted = true;
        this.currentTurn = 'black';
        this.moveHistory = [];
        
        document.getElementById('menu-screen').style.display = 'none';
        document.getElementById('game-controls').style.display = 'block';
        document.getElementById('ai-difficulty').style.display = 'none';
        document.getElementById('room-section').style.display = 'none';
        
        this.board.reset(this.board.size);
        this.startTimer();
        this.updateStatus('Local 2-Player Game Started!');
        this.updateTurnIndicator();
    }

    startMultiplayer() {
        this.gameMode = 'multiplayer';
        document.getElementById('room-section').style.display = 'block';
        document.getElementById('ai-difficulty').style.display = 'none';
        
        this.initializeWebSocket();
    }

    startLearningMode() {
        this.gameMode = 'learn';
        this.gameStarted = true; // Allow interaction in learning mode
        document.getElementById('tutorial-panel').style.display = 'block';
        document.getElementById('puzzle-panel').style.display = 'none';
        document.getElementById('menu-screen').style.display = 'none';
        document.getElementById('game-controls').style.display = 'block';
        
        this.loadLesson(1);
    }

    startPuzzleMode() {
        this.gameMode = 'puzzle';
        this.gameStarted = true; // Allow moves in puzzle mode
        document.getElementById('puzzle-panel').style.display = 'block';
        document.getElementById('tutorial-panel').style.display = 'none';
        document.getElementById('menu-screen').style.display = 'none';
        document.getElementById('game-controls').style.display = 'block';
        
        this.loadPuzzle(1);
    }

    async handleMove(x, y) {
        if (!this.gameStarted) {
            return;
        }
        
        // In AI mode, always allow the human player to move when it's their turn
        if (this.gameMode === 'ai') {
            if (this.currentTurn !== this.playerColor) {
                return;
            }
            
            const moveData = {
                x: x,
                y: y,
                color: this.currentTurn === 'black' ? 1 : 2
            };
            
            this.makeMove(moveData);
            
            // Trigger AI move after a short delay
            setTimeout(() => {
                this.makeAIMove();
            }, 500);
        } else if (this.gameMode === 'multiplayer' && this.wsConnection) {
            if (this.currentTurn !== this.playerColor) {
                return;
            }
            this.wsConnection.sendMove(x, y);
        } else if (this.gameMode === 'puzzle') {
            this.checkPuzzleMove(x, y);
        } else if (this.gameMode === 'learn') {
            // Handle learning exercises
            if (window.learning && window.learning.currentExercise) {
                window.learning.checkExerciseMove(x, y);
            }
        } else if (this.gameMode === 'local') {
            // Local two-player mode
            const moveData = {
                x: x,
                y: y,
                color: this.currentTurn === 'black' ? 1 : 2
            };
            
            this.makeMove(moveData);
        }
    }

    makeMove(moveData) {
        const color = moveData.color;
        const currentBoard = this.board.board.map(row => [...row]);
        currentBoard[moveData.x][moveData.y] = color;
        
        this.board.updateBoard(currentBoard);
        this.board.setLastMove(moveData.x, moveData.y);
        
        this.moveHistory.push({
            move: this.moveHistory.length + 1,
            color: color === 1 ? 'Black' : 'White',
            position: this.getPositionNotation(moveData.x, moveData.y)
        });
        
        this.updateMoveHistory();
        this.currentTurn = this.currentTurn === 'black' ? 'white' : 'black';
        this.updateTurnIndicator();
        
        document.getElementById('move-count').textContent = this.moveHistory.length;
    }

    async makeAIMove() {
        if (this.currentTurn === this.playerColor || !this.gameStarted) {
            console.log('AI not moving - turn:', this.currentTurn, 'player:', this.playerColor, 'started:', this.gameStarted);
            return;
        }
        
        console.log('Requesting AI move for', this.currentTurn, 'with difficulty', this.aiDifficulty);
        
        try {
            const apiUrl = window.location.port === '8081' ? 
                'http://localhost:8081/api/ai-move' : 
                '/api/ai-move';
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    board: this.board.board,
                    boardSize: this.board.size,
                    color: this.currentTurn === 'black' ? 'Black' : 'White',
                    difficulty: this.aiDifficulty
                })
            });
            
            const data = await response.json();
            console.log('AI response:', data);
            
            if (data.pass) {
                console.log('AI is passing');
                this.pass();
            } else if (data.x !== undefined && data.y !== undefined) {
                console.log('AI playing at', data.x, data.y);
                this.makeMove({
                    x: data.x,
                    y: data.y,
                    color: this.currentTurn === 'black' ? 1 : 2
                });
            } else {
                console.error('Invalid AI response:', data);
            }
        } catch (error) {
            console.error('AI move error:', error);
        }
    }

    pass() {
        this.moveHistory.push({
            move: this.moveHistory.length + 1,
            color: this.currentTurn === 'black' ? 'Black' : 'White',
            position: 'Pass'
        });
        
        this.updateMoveHistory();
        this.currentTurn = this.currentTurn === 'black' ? 'white' : 'black';
        this.updateTurnIndicator();
        
        if (this.wsConnection) {
            this.wsConnection.sendPass();
        }
    }

    resign() {
        if (!this.gameStarted) return;
        
        const winner = this.currentTurn === 'black' ? 'White' : 'Black';
        this.showModal('Game Over', `${this.currentTurn === 'black' ? 'Black' : 'White'} resigns. ${winner} wins!`);
        this.gameStarted = false;
        
        if (this.wsConnection) {
            this.wsConnection.sendResign();
        }
    }

    undo() {
        if (this.moveHistory.length === 0) return;
        
        if (this.wsConnection) {
            this.wsConnection.sendUndo();
        }
    }

    newGame() {
        this.gameStarted = false;
        this.moveHistory = [];
        this.currentTurn = 'black';
        this.board.reset(this.board.size);
        
        document.getElementById('menu-screen').style.display = 'block';
        document.getElementById('game-controls').style.display = 'none';
        document.getElementById('room-section').style.display = 'none';
        
        this.stopTimer();
        this.updateStatus('Select a game mode to start playing.');
        this.updateMoveHistory();
    }

    createRoom() {
        if (this.wsConnection) {
            this.wsConnection.createGame(this.board.size);
        }
    }

    joinRoom() {
        const roomId = document.getElementById('room-id-input').value.trim();
        if (roomId && this.wsConnection) {
            this.wsConnection.joinGame(roomId);
        }
    }

    initializeWebSocket() {
        this.wsConnection = new WebSocketConnection(this);
    }

    async loadLesson(lessonId) {
        try {
            const response = await fetch('/api/lessons');
            const lessons = await response.json();
            const lesson = lessons.find(l => l.id === lessonId);
            
            if (lesson) {
                document.getElementById('tutorial-content').innerHTML = `
                    <h3>${lesson.title}</h3>
                    <p>${lesson.description}</p>
                    <div>${lesson.content}</div>
                `;
            }
        } catch (error) {
            console.error('Error loading lesson:', error);
        }
    }

    async loadPuzzle(puzzleId) {
        try {
            const response = await fetch('/api/puzzles');
            const puzzles = await response.json();
            const puzzle = puzzles.find(p => p.id === puzzleId);
            
            if (puzzle) {
                document.getElementById('puzzle-title').textContent = puzzle.title;
                document.getElementById('puzzle-description').textContent = puzzle.description;
                
                // Ensure the board is the right size
                if (puzzle.board && puzzle.board.length > 0) {
                    this.board.reset(puzzle.board.length);
                    this.board.updateBoard(puzzle.board);
                }
                
                this.currentPuzzle = puzzle;
                this.updateStatus(puzzle.description);
            }
        } catch (error) {
            console.error('Error loading puzzle:', error);
        }
    }

    checkPuzzleMove(x, y) {
        if (!this.currentPuzzle) {
            console.log('No puzzle loaded');
            return;
        }
        
        // Check if the move matches the solution
        if (this.currentPuzzle.solution && this.currentPuzzle.solution.moves) {
            const solutionMove = this.currentPuzzle.solution.moves[0];
            if (solutionMove && solutionMove.x === x && solutionMove.y === y) {
                this.showModal('Correct!', 'Well done! You solved the puzzle.');
                // Update the board with the solution move
                const currentBoard = this.board.board.map(row => [...row]);
                currentBoard[x][y] = solutionMove.color || 1;
                this.board.updateBoard(currentBoard);
            } else {
                this.updateStatus('Not quite right. Try again!');
            }
        }
    }

    getPositionNotation(x, y) {
        const letter = String.fromCharCode(65 + x);
        const number = this.board.size - y;
        return `${letter}${number}`;
    }

    updateStatus(message) {
        document.getElementById('game-status').textContent = message;
    }

    updateTurnIndicator() {
        const turnText = this.currentTurn === 'black' ? "Black's Turn" : "White's Turn";
        document.getElementById('current-turn').textContent = turnText;
    }

    updateMoveHistory() {
        const historyDiv = document.getElementById('move-history');
        historyDiv.innerHTML = '';
        
        this.moveHistory.forEach(entry => {
            const moveDiv = document.createElement('div');
            moveDiv.className = 'move-entry';
            moveDiv.textContent = `${entry.move}. ${entry.color} ${entry.position}`;
            historyDiv.appendChild(moveDiv);
        });
        
        historyDiv.scrollTop = historyDiv.scrollHeight;
    }

    startTimer() {
        this.startTime = Date.now();
        this.gameTimer = setInterval(() => {
            const elapsed = Math.floor((Date.now() - this.startTime) / 1000);
            const minutes = Math.floor(elapsed / 60).toString().padStart(2, '0');
            const seconds = (elapsed % 60).toString().padStart(2, '0');
            document.getElementById('game-time').textContent = `${minutes}:${seconds}`;
        }, 1000);
    }

    stopTimer() {
        if (this.gameTimer) {
            clearInterval(this.gameTimer);
            this.gameTimer = null;
        }
        document.getElementById('game-time').textContent = '00:00';
    }

    showModal(title, message) {
        document.getElementById('modal-title').textContent = title;
        document.getElementById('modal-body').textContent = message;
        document.getElementById('modal').style.display = 'flex';
    }

    closeModal() {
        document.getElementById('modal').style.display = 'none';
    }

    onWebSocketMessage(data) {
        if (data.type === 'game_created') {
            this.roomId = data.data.roomId;
            this.playerColor = 'black';
            document.getElementById('room-info').innerHTML = `Room ID: ${this.roomId}<br>Waiting for opponent...`;
        } else if (data.type === 'game_joined') {
            this.roomId = data.data.roomId;
            this.playerColor = 'white';
            document.getElementById('room-info').innerHTML = `Room ID: ${this.roomId}<br>Game started!`;
        } else if (data.type === 'game_started') {
            this.gameStarted = true;
            document.getElementById('menu-screen').style.display = 'none';
            document.getElementById('game-controls').style.display = 'block';
            this.startTimer();
        } else if (data.type === 'move_made') {
            this.board.updateBoard(data.data.board);
            this.board.setLastMove(data.data.x, data.data.y);
            this.currentTurn = data.data.info.currentTurn.toLowerCase();
            this.updateTurnIndicator();
        }
    }
}

window.currentPlayer = 'black';

document.addEventListener('DOMContentLoaded', () => {
    window.game = new GoGame();
});