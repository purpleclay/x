package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/x/theme"
)

func main() {
	out := lipgloss.JoinVertical(
		lipgloss.Top,
		colors(),
		"",
		typography(),
	)
	fmt.Fprint(os.Stdout, lipgloss.NewStyle().Margin(2, 2).Render(out))
}

func colorRow(colors []lipgloss.Color) string {
	colorCell := lipgloss.NewStyle().Height(3).Width(12)
	labelCell := lipgloss.NewStyle().Width(12).AlignHorizontal(lipgloss.Center)

	labels := []string{"50", "100", "200", "300", "400", "500", "600", "700", "800", "900"}

	var colorBlocks []string
	var labelBlocks []string

	for i, c := range colors {
		colorBlocks = append(colorBlocks, colorCell.Background(c).Render())
		labelBlocks = append(labelBlocks, labelCell.Render(labels[i]))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, colorBlocks...),
		lipgloss.JoinHorizontal(lipgloss.Top, labelBlocks...),
	)
}

func colors() string {
	purple := []lipgloss.Color{
		theme.Purple50, theme.Purple100, theme.Purple200, theme.Purple300, theme.Purple400,
		theme.Purple500, theme.Purple600, theme.Purple700, theme.Purple800, theme.Purple900,
	}
	green := []lipgloss.Color{
		theme.Green50, theme.Green100, theme.Green200, theme.Green300, theme.Green400,
		theme.Green500, theme.Green600, theme.Green700, theme.Green800, theme.Green900,
	}
	orange := []lipgloss.Color{
		theme.Orange50, theme.Orange100, theme.Orange200, theme.Orange300, theme.Orange400,
		theme.Orange500, theme.Orange600, theme.Orange700, theme.Orange800, theme.Orange900,
	}
	red := []lipgloss.Color{
		theme.Red50, theme.Red100, theme.Red200, theme.Red300, theme.Red400,
		theme.Red500, theme.Red600, theme.Red700, theme.Red800, theme.Red900,
	}
	blue := []lipgloss.Color{
		theme.Blue50, theme.Blue100, theme.Blue200, theme.Blue300, theme.Blue400,
		theme.Blue500, theme.Blue600, theme.Blue700, theme.Blue800, theme.Blue900,
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		theme.H6.Render("Colors"),
		"",
		colorRow(purple),
		"",
		colorRow(green),
		"",
		colorRow(orange),
		"",
		colorRow(red),
		"",
		colorRow(blue),
	)
}

func typography() string {
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		theme.H1.Render("H1")+"  ",
		theme.H2.Render("H2")+"  ",
		theme.H3.Render("H3")+"  ",
		theme.H4.Render("H4")+"  ",
		theme.H5.Render("H5")+"  ",
		theme.H6.Render("H6"),
	)

	styles := lipgloss.JoinHorizontal(
		lipgloss.Top,
		theme.Bold.Render("Bold")+"  ",
		theme.Italic.Render("Italic")+"  ",
		theme.Underline.Render("Underline")+"  ",
		theme.Strikethrough.Render("Strikethrough")+"  ",
		theme.Code.Render("Code")+"  ",
		theme.Mark.Render("Mark")+"  ",
		theme.Link.Render("Link"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		theme.H6.Render("Typography"),
		"",
		headers,
		"",
		styles,
	)
}
