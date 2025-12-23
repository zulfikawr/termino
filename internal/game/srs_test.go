package game

import (
	"termino/internal/tetromino"
	"testing"
)

func TestRotateCW_Standard(t *testing.T) {
	state := NewGameState()
	state.CurrentPiece = tetromino.NewTetromino("T")
	state.CurrentX = 5
	state.CurrentY = 10
	state.CurrentRotation = 0

	// Test basic rotation (no collision)
	if !state.RotateCW() {
		t.Errorf("Expected RotateCW to succeed in empty space")
	}
	if state.CurrentRotation != 1 {
		t.Errorf("Expected rotation to be 1, got %d", state.CurrentRotation)
	}

	state.RotateCW()
	if state.CurrentRotation != 2 {
		t.Errorf("Expected rotation to be 2, got %d", state.CurrentRotation)
	}

	state.RotateCW()
	if state.CurrentRotation != 3 {
		t.Errorf("Expected rotation to be 3, got %d", state.CurrentRotation)
	}

	state.RotateCW()
	if state.CurrentRotation != 0 {
		t.Errorf("Expected rotation to be 0, got %d", state.CurrentRotation)
	}
}

func TestRotateCW_Kick(t *testing.T) {
	// Implement a test where standard rotation fails but kick succeeds
	state := NewGameState()
	state.CurrentPiece = tetromino.NewTetromino("T")
	// T spawn:
	// . 1 .
	// 1 1 1
	// . . .

	// Rotate CW ->
	// . 1 .
	// . 1 1
	// . 1 .

	// Place a block at (5+1, 10+1) -> (6, 11) relative to board?
	// Piece mask 0: T (3 wide). At X=5: cols 5,6,7.
	// Row 1 (center) is at Y+1.
	// If we block the destination, we force a kick.

	// T Rotate 0->R:
	// 0: (0,0), 1: (-1,0), 2: (-1,1), 3: (0,-2), 4: (-1,-2)
	// Kick 1: Left 1 (dx=-1).

	// To force Kick 1, we must block test 0 (0,0).
	// Test 0: Normal rotation.
	// T-Right normal:
	// . 1 .  (X+1, Y)
	// . 1 1  (X+1, Y+1), (X+2, Y+1)
	// . 1 .  (X+1, Y+2)

	// If we block (X+1, Y+2)?
	state.CurrentX = 5
	state.CurrentY = 10

	// Block (6, 12) (which corresponds to T-Right bottom block at local (1,2) -> (5+1, 10+2))
	// Wait, T-Right mask:
	// 0: 0010 (2)
	// 1: 0110 (6)
	// 2: 0010 (2)
	// 3: 0000
	// Bits 0-3 correspond to cols 0-3.
	// Mask 0010 (2) -> bit 1 set -> col 1.
	// Mask 0110 (6) -> bits 1,2 set -> cols 1,2.
	// So T-Right occupies:
	// (1,0)
	// (1,1), (2,1)
	// (1,2)

	// Absolute positions:
	// (6, 10)
	// (6, 11), (7, 11)
	// (6, 12)

	// Set a block at (6, 12) on the board.
	// state.Board[row] |= (1 << col)
	state.Board[12] |= (1 << 6)

	// Attempt rotate. Normal (0,0) should fail because of (6,12).
	// Next kick: (-1, 0). Left 1.
	// New pos: (4, 10).
	// T-Right at (4,10):
	// (5, 10)
	// (5, 11), (6, 11)
	// (5, 12)
	// Check conflicts. (6, 12) is blocked.
	// (5, 12) is free? Yes.
	// So it should succeed at (-1, 0).

	// Wait, standard T kick 0->1:
	// 0: (0,0)
	// 1: (-1, 0) -> Left 1
	// ...

	if !state.RotateCW() {
		t.Errorf("Expected RotateCW to succeed with kick")
	}

	if state.CurrentX != 4 {
		t.Errorf("Expected kick to shift X to 4, got %d", state.CurrentX)
	}

	if state.CurrentY != 10 {
		t.Errorf("Expected kick Y to be 10, got %d", state.CurrentY)
	}
}

func TestWallKick_I(t *testing.T) {
	// Test I piece wall kick from standard position against right wall
	state := NewGameState()
	state.CurrentPiece = tetromino.NewTetromino("I")
	// I Spawn (0):
	// ....
	// #### (Row 1)
	// ....
	// ....
	// Mask 0x000F (1111) at Row 1. Cols 0,1,2,3.

	// Place against right wall.
	// BoardWidth=10.
	// Cols 0..9.
	// We want I to be at X=6.
	// Cols 6,7,8,9. All valid.
	state.CurrentX = 6
	state.CurrentY = 10

	// Rotate CW (0->1).
	// I-Right (1):
	// ..1.
	// ..1.
	// ..1.
	// ..1.
	// Col 2.
	// At X=6: Col becomes 6+2 = 8.
	// Valid.

	// Wait, let's try to put it at X=7.
	// 7,8,9,10. 10 is OOB.
	// So canPlace should fail initial placement?

	// If we are at X=6 (valid).
	// Rotate to Right.
	// Occupies (8, 10), (8, 11), (8, 12), (8, 13). all valid.

	// Let's force a kick.
	// I kicks 0->1:
	// (0,0), (-2,0), (+1,0), (-2,-1), (+1,+2)

	// Block (8, 10).
	state.Board[10] |= (1 << 8)

	// Try rotate.
	// (0,0) -> fails at (8,10).
	// (-2, 0) -> X becomes 4.
	// I-Right at X=4 -> Col 2 maps to 6.
	// (6, 10) is free? Yes.
	// Should kick left 2.

	state.RotateCW()

	if state.CurrentX != 4 {
		t.Errorf("Expected I-kick X to equal 4, got %d", state.CurrentX)
	}
}
