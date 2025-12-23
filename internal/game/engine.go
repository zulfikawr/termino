package game

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type Model struct {
	State            GameState
	Width            int
	Height           int
	lastSpacePressed bool
}

func NewModel() Model {
	return Model{
		State:  NewGameState(),
		Width:  80, // Default fallback
		Height: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global Controls
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "p", "esc":
			m.State.Paused = !m.State.Paused
			// If unpausing, we simply continue. Ticks are always running.
		case "r":
			m.State = NewGameState() // Reset
			m.State.SpawnNewPiece()  // Ensure piece matches fix
			m.State.UpdateGhost()
		}

		if m.State.GameOver || m.State.Paused {
			return m, nil
		}

		// Gameplay Input
		switch msg.String() {
		case "left":
			if m.State.canPlace(m.State.CurrentPiece.Name, m.State.CurrentX-1, m.State.CurrentY, m.State.CurrentRotation) {
				m.State.CurrentX--
				m.State.resetLockDelay()
				m.State.UpdateGhost()
			}
		case "right":
			if m.State.canPlace(m.State.CurrentPiece.Name, m.State.CurrentX+1, m.State.CurrentY, m.State.CurrentRotation) {
				m.State.CurrentX++
				m.State.resetLockDelay()
				m.State.UpdateGhost()
			}
		case "down":
			if m.State.canPlace(m.State.CurrentPiece.Name, m.State.CurrentX, m.State.CurrentY+1, m.State.CurrentRotation) {
				m.State.CurrentY++
				m.State.Score++
			}
		case "up", "x":
			m.State.RotateCW()
			m.State.UpdateGhost()
		case "c":
			m.State.RotateCCW()
			m.State.UpdateGhost()
		case "v":
			// Rotate 180 degrees (rotate twice)
			m.State.RotateCW()
			m.State.RotateCW()
			m.State.UpdateGhost()
		case " ":
			// Hard Drop - only trigger if not already pressed
			if !m.lastSpacePressed {
				m.lastSpacePressed = true
				dropDist := 0
				for m.State.canPlace(m.State.CurrentPiece.Name, m.State.CurrentX, m.State.CurrentY+1, m.State.CurrentRotation) {
					m.State.CurrentY++
					dropDist++
				}
				m.State.Score += dropDist * 2
				m.State.LockPiece()
				m.State.UpdateGhost()
			}
		case "z":
			m.State.HoldCurrentPiece()
			m.State.UpdateGhost()
		}

	case tickMsg:
		// Reset space bar pressed flag each tick to allow next press
		m.lastSpacePressed = false

		// Stop physics if paused/over
		if !m.State.Paused && !m.State.GameOver {
			m.State.ApplyGravity(1.0 / 60.0)
		}

		return m, tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m Model) View() string {
	return RenderGame(&m.State, m.Width, m.Height)
}
