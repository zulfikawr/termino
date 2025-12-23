package tetromino

import (
	"log"

	"github.com/charmbracelet/lipgloss"
)

type Bitmask uint16

// Tetromino represents a Tetris piece with rotation states, color, and SRS kick table.
type Tetromino struct {
	Name     string
	Masks    [4][4]Bitmask // 4 rotations × 4 rows
	Color    lipgloss.Color
	KickData [4][4][5][2]int // SRS wall kick offsets [from][to][test][x,y]
}

// Standard Tetris guideline colors
var (
	ColorI = lipgloss.Color("#00FFFF")
	ColorJ = lipgloss.Color("#0000FF")
	ColorL = lipgloss.Color("#FF8000")
	ColorO = lipgloss.Color("#FFFF00")
	ColorS = lipgloss.Color("#00FF00")
	ColorT = lipgloss.Color("#800080")
	ColorZ = lipgloss.Color("#FF0000")
)

// Common kick table for J, L, S, T, Z
var kickDataCommon = [4][4][5][2]int{
	// 0 -> R (1), 0 -> L (3)
	0: {
		1: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
		3: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	},
	// R -> 0 (0), R -> 2 (2)
	1: {
		0: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
		2: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	},
	// 2 -> R (1), 2 -> L (3)
	2: {
		1: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
		3: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	},
	// L -> 2 (2), L -> 0 (0)
	3: {
		2: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
		0: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	},
}

// Kick table for I
var kickDataI = [4][4][5][2]int{
	// 0 -> R (1), 0 -> L (3)
	0: {
		1: {{0, 0}, {-2, 0}, {1, 0}, {-2, -1}, {1, 2}},
		3: {{0, 0}, {-1, 0}, {2, 0}, {-1, 2}, {2, -1}},
	},
	// R -> 0 (0), R -> 2 (2)
	1: {
		0: {{0, 0}, {2, 0}, {-1, 0}, {2, 1}, {-1, -2}},
		2: {{0, 0}, {-1, 0}, {2, 0}, {-1, 2}, {2, -1}},
	},
	// 2 -> R (1), 2 -> L (3)
	2: {
		1: {{0, 0}, {1, 0}, {-2, 0}, {1, -2}, {-2, 1}},
		3: {{0, 0}, {2, 0}, {-1, 0}, {2, 1}, {-1, -2}},
	},
	// L -> 2 (2), L -> 0 (0)
	3: {
		2: {{0, 0}, {-2, 0}, {1, 0}, {-2, -1}, {1, 2}},
		0: {{0, 0}, {1, 0}, {-2, 0}, {1, -2}, {-2, 1}},
	},
}

// NewTetromino creates a tetromino piece with the given name, initializing rotation masks and kick data.
func NewTetromino(name string) Tetromino {
	t := Tetromino{Name: name}

	switch name {
	case "I":
		t.Color = ColorI
		t.KickData = kickDataI
		t.Masks[0] = [4]Bitmask{0x0000, 0x000F, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0004, 0x0004, 0x0004, 0x0004}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0000, 0x000F, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0002, 0x0002, 0x0002, 0x0002}

	case "J":
		t.Color = ColorJ
		t.KickData = kickDataCommon
		t.Masks[0] = [4]Bitmask{0x0001, 0x0007, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0006, 0x0002, 0x0002, 0x0000}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0007, 0x0004, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0002, 0x0002, 0x0003, 0x0000}

	case "L":
		t.Color = ColorL
		t.KickData = kickDataCommon
		t.Masks[0] = [4]Bitmask{0x0004, 0x0007, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0002, 0x0002, 0x0006, 0x0000}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0007, 0x0001, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0003, 0x0002, 0x0002, 0x0000}

	case "O":
		t.Color = ColorO
		t.Masks[0] = [4]Bitmask{0x0006, 0x0006, 0x0000, 0x0000}
		t.Masks[1] = t.Masks[0]
		t.Masks[2] = t.Masks[0]
		t.Masks[3] = t.Masks[0]

	case "S":
		t.Color = ColorS
		t.KickData = kickDataCommon
		t.Masks[0] = [4]Bitmask{0x0006, 0x0003, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0002, 0x0006, 0x0004, 0x0000}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0006, 0x0003, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0001, 0x0003, 0x0002, 0x0000}

	case "T":
		t.Color = ColorT
		t.KickData = kickDataCommon
		t.Masks[0] = [4]Bitmask{0x0002, 0x0007, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0002, 0x0006, 0x0002, 0x0000}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0007, 0x0002, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0002, 0x0003, 0x0002, 0x0000}

	case "Z":
		t.Color = ColorZ
		t.KickData = kickDataCommon
		t.Masks[0] = [4]Bitmask{0x0003, 0x0006, 0x0000, 0x0000}
		t.Masks[1] = [4]Bitmask{0x0004, 0x0006, 0x0002, 0x0000}
		t.Masks[2] = [4]Bitmask{0x0000, 0x0003, 0x0006, 0x0000}
		t.Masks[3] = [4]Bitmask{0x0002, 0x0003, 0x0001, 0x0000}
	default:
		log.Fatalf("Unknown tetromino: %s", name)
	}

	return t
}
