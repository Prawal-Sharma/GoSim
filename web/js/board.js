class GoBoard {
    constructor(canvasId, size = 19) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.size = size;
        this.cellSize = 30;
        this.padding = 20;
        this.board = Array(size).fill().map(() => Array(size).fill(0));
        this.lastMove = null;
        this.validMoves = [];
        this.showValidMoves = false;
        this.territoryMarkers = [];
        
        this.setupCanvas();
        this.setupEventListeners();
        this.draw();
    }

    setupCanvas() {
        const totalSize = (this.size - 1) * this.cellSize + 2 * this.padding;
        this.canvas.width = totalSize;
        this.canvas.height = totalSize;
    }

    setupEventListeners() {
        this.canvas.addEventListener('mousemove', (e) => this.handleMouseMove(e));
        this.canvas.addEventListener('click', (e) => this.handleClick(e));
        this.canvas.addEventListener('mouseleave', () => this.handleMouseLeave());
    }

    handleMouseMove(event) {
        const rect = this.canvas.getBoundingClientRect();
        const x = event.clientX - rect.left;
        const y = event.clientY - rect.top;
        
        const gridX = Math.round((x - this.padding) / this.cellSize);
        const gridY = Math.round((y - this.padding) / this.cellSize);
        
        if (gridX >= 0 && gridX < this.size && gridY >= 0 && gridY < this.size) {
            this.hoverPos = { x: gridX, y: gridY };
            this.draw();
        }
    }

    handleMouseLeave() {
        this.hoverPos = null;
        this.draw();
    }

    handleClick(event) {
        const rect = this.canvas.getBoundingClientRect();
        const x = event.clientX - rect.left;
        const y = event.clientY - rect.top;
        
        const gridX = Math.round((x - this.padding) / this.cellSize);
        const gridY = Math.round((y - this.padding) / this.cellSize);
        
        if (gridX >= 0 && gridX < this.size && gridY >= 0 && gridY < this.size) {
            if (this.board[gridX][gridY] === 0) {
                if (this.onMove) {
                    this.onMove(gridX, gridY);
                }
            }
        }
    }

    draw() {
        this.ctx.fillStyle = '#dcb35c';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        
        this.drawGrid();
        this.drawStarPoints();
        this.drawCoordinates();
        
        if (this.showValidMoves) {
            this.drawValidMoves();
        }
        
        this.drawStones();
        
        if (this.territoryMarkers.length > 0) {
            this.drawTerritory();
        }
        
        if (this.lastMove) {
            this.drawLastMoveMarker();
        }
        
        if (this.hoverPos && this.board[this.hoverPos.x][this.hoverPos.y] === 0) {
            this.drawGhostStone();
        }
    }

    drawGrid() {
        this.ctx.strokeStyle = '#000';
        this.ctx.lineWidth = 1;
        
        for (let i = 0; i < this.size; i++) {
            const x = this.padding + i * this.cellSize;
            const y1 = this.padding;
            const y2 = this.padding + (this.size - 1) * this.cellSize;
            
            this.ctx.beginPath();
            this.ctx.moveTo(x, y1);
            this.ctx.lineTo(x, y2);
            this.ctx.stroke();
            
            this.ctx.beginPath();
            this.ctx.moveTo(y1, x);
            this.ctx.lineTo(y2, x);
            this.ctx.stroke();
        }
    }

    drawStarPoints() {
        const starPoints = this.getStarPoints();
        
        this.ctx.fillStyle = '#000';
        starPoints.forEach(point => {
            const x = this.padding + point.x * this.cellSize;
            const y = this.padding + point.y * this.cellSize;
            
            this.ctx.beginPath();
            this.ctx.arc(x, y, 3, 0, 2 * Math.PI);
            this.ctx.fill();
        });
    }

    getStarPoints() {
        if (this.size === 9) {
            return [
                {x: 2, y: 2}, {x: 6, y: 2},
                {x: 4, y: 4},
                {x: 2, y: 6}, {x: 6, y: 6}
            ];
        } else if (this.size === 13) {
            return [
                {x: 3, y: 3}, {x: 9, y: 3},
                {x: 6, y: 6},
                {x: 3, y: 9}, {x: 9, y: 9}
            ];
        } else if (this.size === 19) {
            return [
                {x: 3, y: 3}, {x: 9, y: 3}, {x: 15, y: 3},
                {x: 3, y: 9}, {x: 9, y: 9}, {x: 15, y: 9},
                {x: 3, y: 15}, {x: 9, y: 15}, {x: 15, y: 15}
            ];
        }
        return [];
    }

    drawCoordinates() {
        this.ctx.font = '12px Arial';
        this.ctx.fillStyle = '#666';
        this.ctx.textAlign = 'center';
        this.ctx.textBaseline = 'middle';
        
        for (let i = 0; i < this.size; i++) {
            const letter = String.fromCharCode(65 + i);
            const x = this.padding + i * this.cellSize;
            
            this.ctx.fillText(letter, x, 8);
            this.ctx.fillText(letter, x, this.canvas.height - 8);
            
            const number = this.size - i;
            this.ctx.fillText(number, 8, this.padding + i * this.cellSize);
            this.ctx.fillText(number, this.canvas.width - 8, this.padding + i * this.cellSize);
        }
    }

    drawStones() {
        for (let x = 0; x < this.size; x++) {
            for (let y = 0; y < this.size; y++) {
                if (this.board[x][y] !== 0) {
                    this.drawStone(x, y, this.board[x][y]);
                }
            }
        }
    }

    drawStone(gridX, gridY, color) {
        const x = this.padding + gridX * this.cellSize;
        const y = this.padding + gridY * this.cellSize;
        const radius = this.cellSize * 0.45;
        
        this.ctx.save();
        
        this.ctx.shadowColor = 'rgba(0, 0, 0, 0.5)';
        this.ctx.shadowBlur = 5;
        this.ctx.shadowOffsetX = 2;
        this.ctx.shadowOffsetY = 2;
        
        if (color === 1) {
            const gradient = this.ctx.createRadialGradient(x - radius/3, y - radius/3, 0, x, y, radius);
            gradient.addColorStop(0, '#4a4a4a');
            gradient.addColorStop(1, '#000000');
            this.ctx.fillStyle = gradient;
        } else {
            const gradient = this.ctx.createRadialGradient(x - radius/3, y - radius/3, 0, x, y, radius);
            gradient.addColorStop(0, '#ffffff');
            gradient.addColorStop(1, '#e0e0e0');
            this.ctx.fillStyle = gradient;
        }
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
        this.ctx.fill();
        
        this.ctx.restore();
    }

    drawGhostStone() {
        if (!this.hoverPos) return;
        
        const x = this.padding + this.hoverPos.x * this.cellSize;
        const y = this.padding + this.hoverPos.y * this.cellSize;
        const radius = this.cellSize * 0.45;
        
        this.ctx.globalAlpha = 0.5;
        
        if (window.currentPlayer === 'black') {
            this.ctx.fillStyle = '#000';
        } else {
            this.ctx.fillStyle = '#fff';
        }
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
        this.ctx.fill();
        
        this.ctx.globalAlpha = 1;
    }

    drawLastMoveMarker() {
        if (!this.lastMove) return;
        
        const x = this.padding + this.lastMove.x * this.cellSize;
        const y = this.padding + this.lastMove.y * this.cellSize;
        
        this.ctx.strokeStyle = '#ff0000';
        this.ctx.lineWidth = 2;
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, this.cellSize * 0.2, 0, 2 * Math.PI);
        this.ctx.stroke();
    }

    drawValidMoves() {
        this.ctx.fillStyle = 'rgba(0, 255, 0, 0.3)';
        
        this.validMoves.forEach(move => {
            const x = this.padding + move.x * this.cellSize;
            const y = this.padding + move.y * this.cellSize;
            
            this.ctx.beginPath();
            this.ctx.arc(x, y, this.cellSize * 0.15, 0, 2 * Math.PI);
            this.ctx.fill();
        });
    }

    drawTerritory() {
        this.territoryMarkers.forEach(marker => {
            const x = this.padding + marker.x * this.cellSize;
            const y = this.padding + marker.y * this.cellSize;
            
            if (marker.owner === 1) {
                this.ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
            } else {
                this.ctx.fillStyle = 'rgba(255, 255, 255, 0.7)';
            }
            
            this.ctx.fillRect(x - 5, y - 5, 10, 10);
        });
    }

    updateBoard(boardState) {
        this.board = boardState;
        this.draw();
    }

    setLastMove(x, y) {
        this.lastMove = { x, y };
        this.draw();
    }

    setValidMoves(moves) {
        this.validMoves = moves;
        this.draw();
    }

    toggleValidMoves() {
        this.showValidMoves = !this.showValidMoves;
        this.draw();
    }

    showTerritory(territory) {
        this.territoryMarkers = territory;
        this.draw();
    }

    reset(size) {
        this.size = size;
        this.board = Array(size).fill().map(() => Array(size).fill(0));
        this.lastMove = null;
        this.validMoves = [];
        this.territoryMarkers = [];
        this.setupCanvas();
        this.draw();
    }
}