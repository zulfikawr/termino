package game

import (
	"termino/internal/tetromino"
	"termino/pkg/consts"
	"time"
)

// ApplyGravity updates piece position based on elapsed time and current level.
// Gravity speed increases with level, ranging from 1.25 rows/sec at level 1 to 20+ at higher levels.
func (g *GameState) ApplyGravity(dt float64) {
	speed := g.calculateGravitySpeed()
	g.GravityAccumulator += dt * speed

	for g.GravityAccumulator >= 1.0 {
		if g.canPlace(g.CurrentPiece.Name, g.CurrentX, g.CurrentY+1, g.CurrentRotation) {
			g.CurrentY++
			g.resetLockDelay()
			g.UpdateGhost()
		} else {
			g.handleTouchdown(dt)
		}
		g.GravityAccumulator -= 1.0
	}

	if !g.canPlace(g.CurrentPiece.Name, g.CurrentX, g.CurrentY+1, g.CurrentRotation) {
		g.handleTouchdown(dt)
	}
}

// calculateGravitySpeed returns the gravity speed in rows per second based on current level.
func (g *GameState) calculateGravitySpeed() float64 {
	switch {
	case g.Level < 10:
		return 1.0 + float64(g.Level-1)*0.5
	case g.Level < 20:
		return 5.0 + float64(g.Level-10)*2.0
	default:
		return 20.0
	}
}

// handleTouchdown processes the lock delay timer when piece cannot move down.
func (g *GameState) handleTouchdown(dt float64) {
	g.LockTimer += time.Duration(dt * float64(time.Second))
	if g.LockTimer >= g.LockDelay {
		g.LockPiece()
	}
}

// LockPiece burns the current piece onto the board, clears completed lines, and spawns a new piece.
func (g *GameState) LockPiece() {
	masks := g.CurrentPiece.Masks[g.CurrentRotation]
	for row := 0; row < 4; row++ {
		boardRow := g.CurrentY + row
		if boardRow < 0 || boardRow >= consts.BoardHeight {
			continue
		}
		pieceRowMask := masks[row]

		var shiftedMask uint16
		if g.CurrentX >= 0 {
			shiftedMask = uint16(pieceRowMask) << g.CurrentX
		} else {
			shiftedMask = uint16(pieceRowMask) >> (-g.CurrentX)
		}

		g.Board[boardRow] |= tetromino.Bitmask(shiftedMask)

		for col := 0; col < 4; col++ {
			if (pieceRowMask & tetromino.Bitmask(1<<col)) != 0 {
				boardX := g.CurrentX + col
				if boardX >= 0 && boardX < consts.BoardWidth {
					g.BoardColors[boardRow][boardX] = g.CurrentPiece.Color
				}
			}
		}
	}

	g.ClearLines()
	g.SpawnNewPiece()
	g.UpdateGhost()
}

func (g *GameState) resetLockDelay() {
	g.LockTimer = 0
}

// HoldCurrentPiece swaps the current piece with the held piece (once per lock).
func (g *GameState) HoldCurrentPiece() {
	if g.HoldUsed {
		return
	}

	if g.HoldPiece == nil {
		piece := g.CurrentPiece
		g.HoldPiece = &piece
		g.SpawnNewPiece()
	} else {
		temp := g.CurrentPiece
		g.CurrentPiece = *g.HoldPiece
		g.HoldPiece = &temp

		g.CurrentX = consts.BoardWidth/2 - 2
		g.CurrentY = 18
		g.CurrentRotation = 0
		g.resetLockDelay()
	}

	g.HoldUsed = true
}

// ClearLines removes completed rows and updates score and level accordingly.
func (g *GameState) ClearLines() {
	fullLineMask := uint16(0x03FF)
	linesCleared := 0

	readY := consts.BoardHeight - 1
	writeY := consts.BoardHeight - 1

	for readY >= 0 {
		if (uint16(g.Board[readY]) & fullLineMask) == fullLineMask {
			linesCleared++
			readY--
		} else {
			g.Board[writeY] = g.Board[readY]
			g.BoardColors[writeY] = g.BoardColors[readY]
			writeY--
			readY--
		}
	}

	for writeY >= 0 {
		g.Board[writeY] = 0
		writeY--
	}

	g.updateScore(linesCleared)
	g.updateLevel()
}

// updateScore calculates and applies the score multiplier based on lines cleared.
// Follows Tetris guideline: 100/300/500/800 points per line, multiplied by level.
// Back-to-back tetris (4-line clears) doubles the 800-point base.
func (g *GameState) updateScore(linesCleared int) {
	baseScore := g.calculateLineScore(linesCleared)
	g.Score += baseScore * g.Level
	g.LinesCleared += linesCleared
}

// calculateLineScore returns the base score for the given number of cleared lines.
func (g *GameState) calculateLineScore(lines int) int {
	switch lines {
	case 1:
		g.BackToBack = false
		return 100
	case 2:
		g.BackToBack = false
		return 300
	case 3:
		g.BackToBack = false
		return 500
	case 4:
		if g.BackToBack {
			return 1200
		}
		g.BackToBack = true
		return 800
	}
	return 0
}

// updateLevel increases level every 10 cleared lines.
func (g *GameState) updateLevel() {
	if g.LinesCleared >= g.Level*10 {
		g.Level++
	}
}

// UpdateGhost calculates the lowest valid Y position for the current piece.
func (g *GameState) UpdateGhost() {
	g.GhostY = g.CurrentY
	for i := 0; i < consts.BoardHeight; i++ {
		if g.canPlace(g.CurrentPiece.Name, g.CurrentX, g.GhostY+1, g.CurrentRotation) {
			g.GhostY++
		} else {
			break
		}
	}
}
