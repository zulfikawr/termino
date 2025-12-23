package game

import (
	"termino/internal/tetromino"
	"termino/pkg/consts"
)

// canPlace checks if the current piece can be placed at the specified position and rotation.
// It validates against board boundaries and existing blocks.
func (g *GameState) canPlace(pieceName string, x, y, rotation int) bool {
	currentMasks := g.CurrentPiece.Masks[rotation]

	for row := range 4 {
		boardRow := y + row

		pieceRowMask := currentMasks[row]
		if pieceRowMask == 0 {
			continue
		}

		if boardRow < 0 || boardRow >= consts.BoardHeight {
			return false
		}

		for col := range 4 {
			if (pieceRowMask & (1 << col)) != 0 {
				boardCol := x + col
				if boardCol < 0 || boardCol >= consts.BoardWidth {
					return false
				}

				if (g.Board[boardRow] & tetromino.Bitmask(1<<boardCol)) != 0 {
					return false
				}
			}
		}
	}
	return true
}

// RotateCW attempts a clockwise rotation using SRS kick table.
func (g *GameState) RotateCW() bool {
	newRotation := (g.CurrentRotation + 1) % 4
	return g.tryRotate(newRotation)
}

// RotateCCW attempts a counter-clockwise rotation using SRS kick table.
func (g *GameState) RotateCCW() bool {
	newRotation := (g.CurrentRotation + 3) % 4
	return g.tryRotate(newRotation)
}

// tryRotate attempts rotation with SRS wall kick tests. Returns true if rotation succeeds.
func (g *GameState) tryRotate(newRotation int) bool {
	kickData := g.CurrentPiece.KickData
	from := g.CurrentRotation
	to := newRotation

	for i := range 5 {
		dx := kickData[from][to][i][0]
		dy := kickData[from][to][i][1]

		testX := g.CurrentX + dx
		testY := g.CurrentY - dy

		if g.canPlace(g.CurrentPiece.Name, testX, testY, newRotation) {
			g.CurrentX = testX
			g.CurrentY = testY
			g.CurrentRotation = newRotation
			g.LockResets++
			g.LockTimer = 0
			return true
		}
	}

	return false
}
