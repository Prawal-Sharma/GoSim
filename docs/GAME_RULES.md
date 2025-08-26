# Go Game Rules - Complete Guide

## Introduction
Go (围棋 wéiqí in Chinese, 囲碁 igo in Japanese, 바둑 baduk in Korean) is one of the oldest board games still played today, with a history spanning over 4,000 years. Despite having simple rules, Go offers infinite strategic depth.

## Basic Concepts

### The Board
- Standard sizes: 19×19, 13×13, or 9×9 intersections
- Stones are placed on intersections, not in squares
- The board starts empty

### The Stones
- Two players: Black and White
- Black plays first (traditionally)
- Once placed, stones don't move

### The Objective
Control more territory than your opponent by the end of the game.

## Fundamental Rules

### Rule 1: Alternating Turns
Players take turns placing one stone on an empty intersection.

### Rule 2: Liberties
- **Liberty**: An empty point directly adjacent to a stone (not diagonally)
- A stone in the center has 4 liberties
- A stone on the edge has 3 liberties
- A stone in the corner has 2 liberties

### Rule 3: Connected Groups
Stones of the same color that are adjacent (not diagonally) form a group and share liberties.

### Rule 4: Capture
When a group has no liberties, it is captured and removed from the board.

**Example:**
```
. . . . .
. ● ○ . .
● ○ . ○ .
. ● ○ . .
. . . . .
```
If Black plays at the center, the White stone is captured.

### Rule 5: Suicide Rule
You cannot place a stone that would have no liberties (suicide), UNLESS that move captures enemy stones.

**Illegal move example:**
```
. ○ . 
○ X ○  (X = illegal for Black)
. ○ .
```

**Legal if it captures:**
```
. ○ ● .
○ X ● ○  (X = legal for Black, captures White)
. ○ ● .
```

### Rule 6: Ko Rule
You cannot immediately recapture a stone that just captured your stone if it would recreate the previous board position.

**Ko situation:**
```
. ● ○ .
● . X ○  (After White captures at X)
. ● ○ .

Next turn: Black cannot immediately recapture at the center
```

## Life and Death

### Two Eyes Live
A group with two separate "eyes" (empty spaces surrounded by the group) can never be captured.

**Living group example:**
```
○ ○ ○ ○ ○
○ . ○ . ○  (Two separate eyes)
○ ○ ○ ○ ○
```

### Dead Groups
Groups without two eyes that can be captured are considered "dead."

### Seki (Mutual Life)
When neither player can capture the other without losing their own group.

## Ending the Game

### Passing
- A player may pass instead of placing a stone
- When both players pass consecutively, the game ends

### Territory
After the game ends:
1. Remove dead stones (by agreement)
2. Count empty intersections surrounded by each color
3. Add captured stones (Japanese rules) or stones on board (Chinese rules)

## Scoring Systems

### Japanese Scoring
- **Score = Territory + Captures**
- Count empty points surrounded by your stones
- Add captured enemy stones
- Komi (compensation for White): Usually 6.5 points

### Chinese Scoring
- **Score = Territory + Stones on Board**
- Count empty points surrounded by your stones
- Add your stones remaining on the board
- Komi: Usually 7.5 points

### Area vs Territory
- Japanese rules count territory and captures
- Chinese rules count area (territory + living stones)
- Usually produce similar results

## Advanced Concepts

### Komi
White receives extra points (typically 6.5 or 7.5) to compensate for Black's first-move advantage.

### Handicap Games
Weaker players can place 2-9 stones before White's first move.

### Common Patterns

#### Ladder (Shicho)
A sequence where stones are repeatedly put in atari in a diagonal pattern.

#### Net (Geta)
A technique to capture stones that cannot escape.

#### Snapback
Allowing a stone to be captured to immediately recapture more stones.

#### Ko Threats
Moves that force a response, used to win ko fights.

## Basic Strategy

### Opening Principles
1. **Corners first** - Easiest to secure territory
2. **Sides second** - Next easiest
3. **Center last** - Hardest to secure

### 4-4 Point (Hoshi)
- Balanced between territory and influence
- Common opening move

### 3-3 Point (San-san)
- Secures corner territory
- Less influence on the center

### Common Proverbs
- "The empty triangle is bad" - Inefficient shape
- "Hane at the head of two stones" - Standard technique
- "Don't approach thickness" - Avoid playing near strong positions
- "Urgent points before big points" - Safety first

## Etiquette

### Traditional Etiquette
- Bow before and after the game
- Hold stones between index and middle finger
- Place stones firmly and decisively
- Don't hover over the board
- Resign when the game is clearly lost

### Online Etiquette
- Greet opponent at start
- Don't disconnect when losing
- Thank opponent after game
- Don't use excessive time

## Ranks and Ratings

### Kyu and Dan System
- **Kyu (級)**: Student ranks (30k-1k, weaker to stronger)
- **Dan (段)**: Master ranks (1d-9d amateur, 1p-9p professional)

### Typical Progress
- 30k-20k: Complete beginner
- 20k-10k: Basic understanding
- 10k-5k: Intermediate player
- 5k-1k: Advanced amateur
- 1d-5d: Strong amateur
- 6d+: Very strong amateur
- Professional: Separate ranking system

## Common Mistakes for Beginners

1. **Playing too close to opponent's strong groups**
2. **Not making eyes for groups**
3. **Ignoring the whole board**
4. **Following opponent around**
5. **Not studying life and death**
6. **Playing moves without purpose**
7. **Creating weak groups**
8. **Not protecting cutting points**

## Tips for Improvement

1. **Solve problems daily** - Life/death and tesuji
2. **Play games** - Both serious and casual
3. **Review your games** - Learn from mistakes
4. **Study professional games** - See proper technique
5. **Learn basic joseki** - Corner patterns
6. **Practice reading** - Calculate variations
7. **Join a Go club** - Learn from stronger players
8. **Use learning resources** - Books, videos, apps

## Glossary

- **Atari**: A group with only one liberty
- **Joseki**: Standard corner sequences
- **Tesuji**: Clever tactical move
- **Sente**: Initiative, forcing move
- **Gote**: Defensive move, losing initiative
- **Moyo**: Large territorial framework
- **Thickness**: Strong, influential position
- **Aji**: Latent possibilities in a position
- **Tenuki**: Playing elsewhere, ignoring local situation

## Resources for Learning

### Books
- "Go: A Complete Introduction to the Game" by Cho Chikun
- "The Second Book of Go" by Richard Bozulich
- "Lessons in the Fundamentals of Go" by Kageyama Toshiro

### Websites
- Online-Go.com (OGS) - Play online
- GoKGS.com - KGS Go Server
- Sensei's Library - Go wiki

### Problem Sites
- 101weiqi.com - Thousands of problems
- goproblems.com - Interactive problems
- tsumego.tasuki.org - Daily problems

## Using GoSim to Learn

GoSim provides several features to help you learn:

1. **Progressive Difficulty**: Start with 9×9 boards
2. **AI Opponents**: Practice against different skill levels
3. **Puzzle Mode**: Solve tactical problems
4. **Tutorial System**: Interactive lessons
5. **Game Analysis**: Review your moves
6. **Territory Visualization**: See controlled areas

Start with the tutorials, practice with puzzles, then play against the AI at increasing difficulties!