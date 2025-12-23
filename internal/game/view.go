package game

import (
	"fmt"
	"termino/internal/render"
	"termino/internal/tetromino"
	"termino/pkg/consts"

	"github.com/charmbracelet/lipgloss"
)

// Global buffer instance
var ScreenBuffer *render.Buffer

func InitScreen() {
	ScreenBuffer = render.NewBuffer(80, 24)
}

func RenderGame(state *GameState, screenW, screenH int) string {
	if screenW == 0 {
		screenW = 80
	}
	if screenH == 0 {
		screenH = 24
	}

	if ScreenBuffer == nil || ScreenBuffer.Width() != screenW || ScreenBuffer.Height() != screenH {
		ScreenBuffer = render.NewBuffer(screenW, screenH)
	}

	ScreenBuffer.Reset()

	boardPixelW := consts.BoardWidth*2 + 2
	boardPixelH := consts.VisibleHeight + 2

	offsetX := (screenW - boardPixelW) / 2
	offsetY := (screenH - boardPixelH) / 2

	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}

	drawBox(ScreenBuffer, offsetX, offsetY, consts.BoardWidth+1, consts.VisibleHeight+2, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))

	visibleStart := consts.BoardHeight - consts.VisibleHeight

	for y := range consts.VisibleHeight {
		boardRowIdx := visibleStart + y
		rowMask := state.Board[boardRowIdx]

		for x := range consts.BoardWidth {
			if (rowMask & tetromino.Bitmask(1<<x)) != 0 {
				col := state.BoardColors[boardRowIdx][x]
				if col == "" {
					col = lipgloss.Color("#888888")
				}
				drawBlock(ScreenBuffer, offsetX+1+x*2, offsetY+1+y, col)
			}
		}
	}

	ghostY := state.GhostY
	drawGhost(ScreenBuffer, state.CurrentPiece, state.CurrentX, ghostY, state.CurrentRotation, offsetX+1, offsetY+1, visibleStart)
	drawTetromino(ScreenBuffer, state.CurrentPiece, state.CurrentX, state.CurrentY, state.CurrentRotation, offsetX+1, offsetY+1, visibleStart)
	drawUI(ScreenBuffer, state, offsetX, offsetY)

	return ScreenBuffer.Render()
}

// drawTetromino renders the current falling piece to the buffer.
func drawTetromino(b *render.Buffer, piece tetromino.Tetromino, px, py, rot, offX, offY, visibleStart int) {
	mask := piece.Masks[rot]

	for r := range 4 {
		boardY := py + r
		if boardY < visibleStart {
			continue
		}

		screenY := offY + (boardY - visibleStart)
		rowBits := mask[r]
		for c := range 4 {
			if (rowBits & tetromino.Bitmask(1<<c)) != 0 {
				boardX := px + c
				if boardX >= 0 && boardX < consts.BoardWidth {
					drawBlock(b, offX+boardX*2, screenY, piece.Color)
				}
			}
		}
	}
}

// drawGhost renders a faint preview of where the piece will land.
func drawGhost(b *render.Buffer, piece tetromino.Tetromino, px, py, rot, offX, offY, visibleStart int) {
	mask := piece.Masks[rot]
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Faint(true)

	for r := range 4 {
		boardY := py + r
		if boardY < visibleStart {
			continue
		}
		screenY := offY + (boardY - visibleStart)
		rowBits := mask[r]
		for c := range 4 {
			if (rowBits & tetromino.Bitmask(1<<c)) != 0 {
				boardX := px + c
				if boardX >= 0 && boardX < consts.BoardWidth {
					b.Set(offX+boardX*2, screenY, '█', style)
					b.Set(offX+boardX*2+1, screenY, '█', style)
				}
			}
		}
	}
}

func drawBlock(b *render.Buffer, x, y int, color lipgloss.Color) {
	style := lipgloss.NewStyle().Foreground(color)
	b.Set(x, y, '█', style)
	b.Set(x+1, y, '█', style)
}

func drawBox(b *render.Buffer, x, y, w, h int, style lipgloss.Style) {
	b.Set(x, y, '┌', style)
	widthChars := w * 2
	b.Set(x+widthChars-1, y, '┐', style)
	b.Set(x, y+h-1, '└', style)
	b.Set(x+widthChars-1, y+h-1, '┘', style)

	for i := 1; i < widthChars-1; i++ {
		b.Set(x+i, y, '─', style)
		b.Set(x+i, y+h-1, '─', style)
	}

	for i := 1; i < h-1; i++ {
		b.Set(x, y+i, '│', style)
		b.Set(x+widthChars-1, y+i, '│', style)
	}
}

// drawUI renders score, level, hold, and next queue displays, plus game status overlays.
func drawUI(b *render.Buffer, state *GameState, x, y int) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	writeString(b, x-10, y, "Hold:", style)
	if state.HoldPiece != nil {
		drawMiniPiece(b, *state.HoldPiece, x-10, y+2)
	}

	writeString(b, x-10, y+8, "Score:", style)
	writeString(b, x-10, y+9, fmt.Sprintf("%d", state.Score), style)

	writeString(b, x-10, y+11, fmt.Sprintf("Lvl: %d", state.Level), style)
	writeString(b, x-10, y+13, fmt.Sprintf("Lns: %d", state.LinesCleared), style)

	writeString(b, x+24, y, "Next:", style)
	for i, piece := range state.NextQueue {
		if i > 2 {
			break
		}
		drawMiniPiece(b, piece, x+24, y+2+i*4)
	}

	if state.GameOver {
		b.DimArea(x+1, y+1, consts.BoardWidth*2, consts.VisibleHeight)
		writeString(b, x+6, y+8, "GAME OVER", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true))
		writeString(b, x+6, y+10, "Press 'r'", style)
		writeString(b, x+7, y+11, "to Retry", style)
	} else if state.Paused {
		b.DimArea(x+1, y+1, consts.BoardWidth*2, consts.VisibleHeight)
		writeString(b, x+8, y+9, "PAUSED", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true))
	}
}

func writeString(b *render.Buffer, x, y int, text string, style lipgloss.Style) {
	for i, r := range text {
		b.Set(x+i, y, r, style)
	}
}

// drawMiniPiece renders a 4x4 tetromino piece for the hold and next queue display.
func drawMiniPiece(b *render.Buffer, piece tetromino.Tetromino, x, y int) {
	mask := piece.Masks[0]

	for r := range 4 {
		rowBits := mask[r]
		for c := range 4 {
			if (rowBits & tetromino.Bitmask(1<<c)) != 0 {
				drawBlock(b, x+c*2, y+r, piece.Color)
			}
		}
	}
}
