package game

import (
	"time"

	"termino/internal/tetromino"
	"termino/pkg/consts"

	"github.com/charmbracelet/lipgloss"
)

type GameState struct {
	Board              [consts.BoardHeight]tetromino.Bitmask                 // Bitboard for collision
	BoardColors        [consts.BoardHeight][consts.BoardWidth]lipgloss.Color // Color board for rendering
	CurrentPiece       tetromino.Tetromino
	CurrentX, CurrentY int
	CurrentRotation    int
	HoldPiece          *tetromino.Tetromino
	HoldUsed           bool
	NextQueue          []tetromino.Tetromino // Circular buffer or just a slice from Randomizer logic
	Randomizer         *Randomizer

	Score        int
	Level        int
	LinesCleared int
	BackToBack   bool
	Combo        int

	LockDelay    time.Duration
	LockTimer    time.Duration
	LockResets   int
	LastLockTime time.Time

	DasTimer           time.Duration // Delayed Auto Shift
	ArrTimer           time.Duration // Auto Repeat Rate
	GravityAccumulator float64

	GhostY int

	GameOver bool
	Paused   bool
}

func NewGameState() GameState {
	r := NewRandomizer()
	queue := make([]tetromino.Tetromino, 0, consts.PreviewCount)

	for range consts.PreviewCount {
		queue = append(queue, r.Next())
	}

	g := GameState{
		Randomizer: r,
		NextQueue:  queue,
		Level:      1,
		LockDelay:  time.Millisecond * 500,
	}
	g.SpawnNewPiece()
	g.UpdateGhost()
	return g
}

// SpawnNewPiece retrieves the next piece from the queue, updates the next queue,
// and initializes it at the spawn position. Returns false if the spawn position collides (game over).
func (g *GameState) SpawnNewPiece() bool {
	g.CurrentPiece = g.NextQueue[0]
	g.NextQueue = g.NextQueue[1:]
	g.NextQueue = append(g.NextQueue, g.Randomizer.Next())

	g.CurrentX = consts.BoardWidth/2 - 2
	g.CurrentY = 20
	g.CurrentRotation = 0
	g.HoldUsed = false
	g.LockResets = 0
	g.LockTimer = 0

	if !g.canPlace(g.CurrentPiece.Name, g.CurrentX, g.CurrentY, g.CurrentRotation) {
		g.GameOver = true
		return false
	}
	return true
}
