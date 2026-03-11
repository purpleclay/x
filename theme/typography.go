package theme

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

// Base header style with padding, bold, and white foreground.
var header = lipgloss.NewStyle().Padding(0, 3).Bold(true).Foreground(lipgloss.Color("#ffffff"))

// Header styles ranked by importance from H1 (most) to H6 (least).
// Each header uses progressively darker background shades and
// adapts for light and dark terminals.
var (
	// H1 is the most prominent header style.
	H1 = header.Background(compat.AdaptiveColor{
		Light: Purple50,
		Dark:  Purple200,
	})

	// H2 is a prominent header style.
	H2 = header.Background(compat.AdaptiveColor{
		Light: Purple100,
		Dark:  Purple300,
	})

	// H3 is a mid-level header style.
	H3 = header.Background(compat.AdaptiveColor{
		Light: Purple200,
		Dark:  Purple400,
	})

	// H4 is a mid-level header style.
	H4 = header.Background(compat.AdaptiveColor{
		Light: Purple300,
		Dark:  Purple500,
	})

	// H5 is a less prominent header style.
	H5 = header.Background(compat.AdaptiveColor{
		Light: Purple400,
		Dark:  Purple600,
	})

	// H6 is the least prominent header style.
	H6 = header.Background(compat.AdaptiveColor{
		Light: Purple500,
		Dark:  Purple700,
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
		Foreground(compat.AdaptiveColor{
			Light: Purple500,
			Dark:  Purple50,
		})

	// Mark renders highlighted text with a background color.
	Mark = lipgloss.NewStyle().
		Padding(0, 1).
		Background(compat.AdaptiveColor{
			Light: Purple50,
			Dark:  Purple700,
		})
)

// Link styles.
var (
	// Link renders a hyperlink with bold, underline, and themed color.
	Link = lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(compat.AdaptiveColor{
			Light: Purple400,
			Dark:  Purple100,
		})
)
