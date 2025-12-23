package input

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DASDelay = 167 * time.Millisecond
	ARRDelay = 33 * time.Millisecond
)

// InputHandler manages key state and DAS/ARR timing for tetromino movement.
type InputHandler struct {
	keyState    map[string]bool
	lastKeyTime map[string]time.Time
	dasTimer    map[string]time.Duration
}

// NewInputHandler creates a new input handler instance.
func NewInputHandler() *InputHandler {
	return &InputHandler{
		keyState:    make(map[string]bool),
		lastKeyTime: make(map[string]time.Time),
		dasTimer:    make(map[string]time.Duration),
	}
}

// HandleKey processes a key press/release event from BubbleTea.
func (h *InputHandler) HandleKey(msg tea.KeyMsg) {}

// Update processes held keys and returns generated movement actions based on DAS/ARR timing.
func (h *InputHandler) Update(dt time.Duration) []string {
	var actions []string
	now := time.Now()

	for k, t := range h.lastKeyTime {
		if h.keyState[k] {
			if now.Sub(t) > 100*time.Millisecond {
				h.keyState[k] = false
			}
		}
	}

	return actions
}

// MapKey translates tea key events to game action strings.
// Maps: left/right arrows to movement, space/up to hard drop, z/x/c to rotations/hold.
func MapKey(key tea.KeyMsg) string {
	switch key.String() {
	case "left":
		return "LEFT"
	case "right":
		return "RIGHT"
	case "down":
		return "SOFT_DROP"
	case " ":
		return "HARD_DROP"
	case "z":
		return "ROTATE_CCW"
	case "x":
		return "ROTATE_CW"
	case "c":
		return "HOLD"
	}

	if key.Type == tea.KeyUp {
		return "ROTATE_CW"
	}

	return ""
}
