package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Buffer struct {
	width, height int
	current       [][]rune
	previous      [][]rune
	styles        [][]lipgloss.Style
}

func (b *Buffer) Width() int  { return b.width }
func (b *Buffer) Height() int { return b.height }

func (b *Buffer) DimArea(x, y, w, h int) {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	for r := y; r < y+h; r++ {
		for c := x; c < x+w; c++ {
			if r >= 0 && r < b.height && c >= 0 && c < b.width {
				// Check if char is not empty space to retain shape
				if b.current[r][c] != ' ' {
					b.styles[r][c] = dimStyle
				} else {
					// Optionally dim empty space?
					// b.current[r][c] = '·' // Dotted background?
					// b.styles[r][c] = dimStyle
				}
			}
		}
	}
}

func NewBuffer(width, height int) *Buffer {
	b := &Buffer{
		width:    width,
		height:   height,
		current:  make([][]rune, height),
		previous: make([][]rune, height),
		styles:   make([][]lipgloss.Style, height),
	}

	for i := range b.current {
		b.current[i] = make([]rune, width)
		b.previous[i] = make([]rune, width)
		b.styles[i] = make([]lipgloss.Style, width)

		for j := range b.current[i] {
			b.current[i][j] = ' '
			b.previous[i][j] = ' '
		}
	}
	return b
}

func (b *Buffer) Set(x, y int, char rune, style lipgloss.Style) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.current[y][x] = char
	b.styles[y][x] = style
}

func (b *Buffer) Render() string {
	var sb strings.Builder

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			// Render the cell
			sb.WriteString(b.styles[y][x].Render(string(b.current[y][x])))
		}
		if y < b.height-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// Reset clears the current buffer for the next frame
func (b *Buffer) Reset() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.current[y][x] = ' '
			b.styles[y][x] = lipgloss.NewStyle()
		}
	}
}
