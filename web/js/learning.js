class LearningModule {
    constructor(game) {
        this.game = game;
        this.currentLesson = null;
        this.currentPuzzle = null;
        this.progress = {
            rules: 0,
            tactics: 0,
            strategy: 0
        };
        
        this.lessons = [];
        this.puzzles = [];
        
        this.setupEventListeners();
        this.loadProgress();
    }

    setupEventListeners() {
        document.getElementById('hint-btn').addEventListener('click', () => this.showHint());
        document.getElementById('solution-btn').addEventListener('click', () => this.showSolution());
    }

    loadProgress() {
        const savedProgress = localStorage.getItem('goSimProgress');
        if (savedProgress) {
            this.progress = JSON.parse(savedProgress);
            this.updateProgressBars();
        }
    }

    saveProgress() {
        localStorage.setItem('goSimProgress', JSON.stringify(this.progress));
        this.updateProgressBars();
    }

    updateProgressBars() {
        const progressItems = document.querySelectorAll('.progress-item');
        progressItems[0].querySelector('.progress-fill').style.width = `${this.progress.rules}%`;
        progressItems[1].querySelector('.progress-fill').style.width = `${this.progress.tactics}%`;
        progressItems[2].querySelector('.progress-fill').style.width = `${this.progress.strategy}%`;
    }

    async loadLessons() {
        try {
            const response = await fetch('/api/lessons');
            this.lessons = await response.json();
            return this.lessons;
        } catch (error) {
            console.error('Error loading lessons:', error);
            return [];
        }
    }

    async loadPuzzles() {
        try {
            const response = await fetch('/api/puzzles');
            this.puzzles = await response.json();
            return this.puzzles;
        } catch (error) {
            console.error('Error loading puzzles:', error);
            return [];
        }
    }

    displayLesson(lessonId) {
        const lesson = this.lessons.find(l => l.id === lessonId);
        if (!lesson) return;

        this.currentLesson = lesson;
        
        const content = document.getElementById('tutorial-content');
        content.innerHTML = `
            <h3>${lesson.title}</h3>
            <p class="lesson-description">${lesson.description}</p>
            <div class="lesson-content">${this.formatLessonContent(lesson.content)}</div>
            ${this.getLessonExercises(lesson.level)}
        `;

        if (lesson.boardSetup) {
            this.setupLessonBoard(lesson.boardSetup);
        }
    }

    formatLessonContent(content) {
        const sections = content.split('\n\n');
        return sections.map(section => {
            if (section.startsWith('*')) {
                return `<ul>${section.split('\n').map(line => 
                    `<li>${line.replace('*', '').trim()}</li>`
                ).join('')}</ul>`;
            } else if (section.startsWith('#')) {
                return `<h4>${section.replace('#', '').trim()}</h4>`;
            } else {
                return `<p>${section}</p>`;
            }
        }).join('');
    }

    getLessonExercises(level) {
        const exercises = {
            beginner: `
                <div class="lesson-exercises">
                    <h4>Practice Exercise</h4>
                    <button onclick="window.learning.startExercise('capture')">Practice Capturing</button>
                    <button onclick="window.learning.startExercise('territory')">Practice Territory</button>
                </div>
            `,
            intermediate: `
                <div class="lesson-exercises">
                    <h4>Practice Exercise</h4>
                    <button onclick="window.learning.startExercise('ladder')">Practice Ladders</button>
                    <button onclick="window.learning.startExercise('life-death')">Life & Death Problems</button>
                </div>
            `,
            advanced: `
                <div class="lesson-exercises">
                    <h4>Practice Exercise</h4>
                    <button onclick="window.learning.startExercise('joseki')">Study Joseki</button>
                    <button onclick="window.learning.startExercise('endgame')">Endgame Practice</button>
                </div>
            `
        };
        
        return exercises[level] || '';
    }

    setupLessonBoard(setup) {
        const board = Array(this.game.board.size).fill().map(() => 
            Array(this.game.board.size).fill(0)
        );
        
        if (setup.stones) {
            setup.stones.forEach(stone => {
                board[stone.x][stone.y] = stone.color;
            });
        }
        
        this.game.board.updateBoard(board);
        
        if (setup.markers) {
            this.addBoardMarkers(setup.markers);
        }
    }

    displayPuzzle(puzzleId) {
        const puzzle = this.puzzles.find(p => p.id === puzzleId);
        if (!puzzle) return;

        this.currentPuzzle = puzzle;
        
        document.getElementById('puzzle-title').textContent = puzzle.title;
        document.getElementById('puzzle-description').textContent = puzzle.description;
        
        if (puzzle.board) {
            this.game.board.updateBoard(puzzle.board);
        }
        
        this.puzzleSolution = puzzle.solution;
        this.puzzleHints = puzzle.hints || [];
        this.currentHintIndex = 0;
    }

    startExercise(type) {
        switch(type) {
            case 'capture':
                this.setupCaptureExercise();
                break;
            case 'territory':
                this.setupTerritoryExercise();
                break;
            case 'ladder':
                this.setupLadderExercise();
                break;
            case 'life-death':
                this.setupLifeDeathExercise();
                break;
            case 'joseki':
                this.setupJosekiExercise();
                break;
            case 'endgame':
                this.setupEndgameExercise();
                break;
        }
    }

    setupCaptureExercise() {
        const board = Array(9).fill().map(() => Array(9).fill(0));
        
        board[4][4] = 2;
        board[3][4] = 1;
        board[5][4] = 1;
        board[4][3] = 1;
        
        this.game.board.reset(9);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Capture the white stone!');
        this.currentExercise = {
            type: 'capture',
            solution: { x: 4, y: 5 },
            check: (x, y) => x === 4 && y === 5
        };
    }

    setupTerritoryExercise() {
        const board = Array(9).fill().map(() => Array(9).fill(0));
        
        for (let i = 0; i < 9; i++) {
            board[i][4] = i < 4 ? 1 : 0;
        }
        
        this.game.board.reset(9);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Build territory in the corner!');
        this.currentExercise = {
            type: 'territory',
            check: (x, y) => this.checkTerritoryMove(x, y)
        };
    }

    setupLadderExercise() {
        const board = Array(13).fill().map(() => Array(13).fill(0));
        
        board[6][6] = 2;
        board[5][6] = 1;
        board[6][5] = 1;
        
        this.game.board.reset(13);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Can you capture white with a ladder?');
        this.currentExercise = {
            type: 'ladder',
            solution: { x: 7, y: 6 },
            check: (x, y) => this.checkLadderMove(x, y)
        };
    }

    setupLifeDeathExercise() {
        const board = Array(9).fill().map(() => Array(9).fill(0));
        
        for (let i = 2; i <= 6; i++) {
            board[i][2] = 1;
            board[i][3] = 2;
        }
        board[3][2] = 0;
        board[5][2] = 0;
        
        this.game.board.reset(9);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Make the black group alive!');
        this.currentExercise = {
            type: 'life-death',
            check: (x, y) => this.checkLifeDeathMove(x, y)
        };
    }

    setupJosekiExercise() {
        const board = Array(19).fill().map(() => Array(19).fill(0));
        
        board[3][3] = 1;
        board[15][3] = 2;
        
        this.game.board.reset(19);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Play the standard joseki response');
        this.currentExercise = {
            type: 'joseki',
            check: (x, y) => this.checkJosekiMove(x, y)
        };
    }

    setupEndgameExercise() {
        const board = Array(9).fill().map(() => Array(9).fill(0));
        
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                if (i < 4) board[i][j] = 1;
                else if (i > 5) board[i][j] = 2;
            }
        }
        
        board[4][4] = 0;
        board[5][4] = 0;
        
        this.game.board.reset(9);
        this.game.board.updateBoard(board);
        
        this.game.updateStatus('Find the best endgame move');
        this.currentExercise = {
            type: 'endgame',
            check: (x, y) => this.checkEndgameMove(x, y)
        };
    }

    checkExerciseMove(x, y) {
        if (!this.currentExercise) return;
        
        const correct = this.currentExercise.check(x, y);
        
        if (correct) {
            this.game.showModal('Correct!', 'Well done! You solved the exercise.');
            this.updateProgress(this.currentExercise.type);
        } else {
            this.game.updateStatus('Not quite right. Try again!');
        }
        
        return correct;
    }

    checkTerritoryMove(x, y) {
        return (x < 3 && y < 3) || (x < 3 && y > 5);
    }

    checkLadderMove(x, y) {
        const moves = [
            { x: 7, y: 6 },
            { x: 7, y: 7 },
            { x: 8, y: 5 }
        ];
        return moves.some(m => m.x === x && m.y === y);
    }

    checkLifeDeathMove(x, y) {
        return (x === 3 || x === 5) && y === 2;
    }

    checkJosekiMove(x, y) {
        const standardMoves = [
            { x: 16, y: 3 },
            { x: 15, y: 16 },
            { x: 3, y: 15 }
        ];
        return standardMoves.some(m => m.x === x && m.y === y);
    }

    checkEndgameMove(x, y) {
        return (x === 4 || x === 5) && y === 4;
    }

    showHint() {
        if (!this.currentPuzzle || !this.puzzleHints.length) {
            this.game.updateStatus('No hints available for this puzzle.');
            return;
        }
        
        if (this.currentHintIndex < this.puzzleHints.length) {
            const hint = this.puzzleHints[this.currentHintIndex];
            this.game.showModal('Hint', hint);
            this.currentHintIndex++;
        } else {
            this.game.showModal('No More Hints', 'You have seen all available hints.');
        }
    }

    showSolution() {
        if (!this.currentPuzzle || !this.puzzleSolution) {
            this.game.updateStatus('No solution available for this puzzle.');
            return;
        }
        
        const solution = this.puzzleSolution;
        if (solution.moves) {
            this.displaySolutionMoves(solution.moves);
        }
        
        if (solution.explanation) {
            this.game.showModal('Solution', solution.explanation);
        }
    }

    displaySolutionMoves(moves) {
        let index = 0;
        const interval = setInterval(() => {
            if (index < moves.length) {
                const move = moves[index];
                const currentBoard = this.game.board.board;
                currentBoard[move.x][move.y] = move.color;
                this.game.board.updateBoard(currentBoard);
                this.game.board.setLastMove(move.x, move.y);
                index++;
            } else {
                clearInterval(interval);
            }
        }, 1000);
    }

    updateProgress(exerciseType) {
        const progressMap = {
            'capture': 'rules',
            'territory': 'rules',
            'ladder': 'tactics',
            'life-death': 'tactics',
            'joseki': 'strategy',
            'endgame': 'strategy'
        };
        
        const category = progressMap[exerciseType];
        if (category) {
            this.progress[category] = Math.min(100, this.progress[category] + 10);
            this.saveProgress();
        }
    }

    addBoardMarkers(markers) {
        markers.forEach(marker => {
            if (marker.type === 'circle') {
                this.drawCircleMarker(marker.x, marker.y, marker.color);
            } else if (marker.type === 'square') {
                this.drawSquareMarker(marker.x, marker.y, marker.color);
            } else if (marker.type === 'triangle') {
                this.drawTriangleMarker(marker.x, marker.y, marker.color);
            }
        });
    }

    drawCircleMarker(x, y, color) {
        const ctx = this.game.board.ctx;
        const posX = this.game.board.padding + x * this.game.board.cellSize;
        const posY = this.game.board.padding + y * this.game.board.cellSize;
        
        ctx.strokeStyle = color || 'red';
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(posX, posY, 10, 0, 2 * Math.PI);
        ctx.stroke();
    }

    drawSquareMarker(x, y, color) {
        const ctx = this.game.board.ctx;
        const posX = this.game.board.padding + x * this.game.board.cellSize;
        const posY = this.game.board.padding + y * this.game.board.cellSize;
        
        ctx.strokeStyle = color || 'blue';
        ctx.lineWidth = 2;
        ctx.strokeRect(posX - 8, posY - 8, 16, 16);
    }

    drawTriangleMarker(x, y, color) {
        const ctx = this.game.board.ctx;
        const posX = this.game.board.padding + x * this.game.board.cellSize;
        const posY = this.game.board.padding + y * this.game.board.cellSize;
        
        ctx.strokeStyle = color || 'green';
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.moveTo(posX, posY - 10);
        ctx.lineTo(posX - 8, posY + 8);
        ctx.lineTo(posX + 8, posY + 8);
        ctx.closePath();
        ctx.stroke();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (window.game) {
        window.learning = new LearningModule(window.game);
    }
});