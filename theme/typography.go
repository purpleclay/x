package theme

import "github.com/charmbracelet/lipgloss"

// Base header style with padding, bold, and white foreground.
var header = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("#ffffff"))

// Header styles ranked by importance from H1 (most) to H6 (least).
// Each header uses progressively darker background shades and
// adapts for light and dark terminals.
var (
	// H1 is the most prominent header style.
	H1 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple50),
		Dark:  string(Purple200),
	})

	// H2 is a prominent header style.
	H2 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple100),
		Dark:  string(Purple300),
	})

	// H3 is a mid-level header style.
	H3 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple200),
		Dark:  string(Purple400),
	})

	// H4 is a mid-level header style.
	H4 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple300),
		Dark:  string(Purple500),
	})

	// H5 is a less prominent header style.
	H5 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple400),
		Dark:  string(Purple600),
	})

	// H6 is the least prominent header style.
	H6 = header.Background(lipgloss.AdaptiveColor{
		Light: string(Purple500),
		Dark:  string(Purple700),
	})
)

// Text decoration styles.
var (
	// Bold renders text in bold.
	Bold = lipgloss.NewStyle().Bold(true)

	// Italic renders text in italic.
	Italic = lipgloss.NewStyle().Italic(true)

	// Underline renders text with an underline.
	Underline = lipgloss.NewStyle().Underline(true)

	// Strikethrough renders text with a strikethrough.
	Strikethrough = lipgloss.NewStyle().Strikethrough(true)

	// Code renders inline code with a distinct style.
	Code = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: string(Purple500),
			Dark:  string(Purple50),
		})

	// Mark renders highlighted text with a background color.
	Mark = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: string(Purple50),
			Dark:  string(Purple700),
		})
)

// Link styles.
var (
	// Link renders a hyperlink with bold, underline, and themed color.
	Link = lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.AdaptiveColor{
			Light: string(Purple400),
			Dark:  string(Purple100),
		})
)
