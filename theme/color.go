package theme

import "github.com/charmbracelet/lipgloss"

// Purple is the primary color palette derived from the PurpleClay brand.
// Shades range from lightest (Purple50) to darkest (Purple900).
var (
	Purple50  = lipgloss.Color("#a980db")
	Purple100 = lipgloss.Color("#906ccf")
	Purple200 = lipgloss.Color("#7958c3")
	Purple300 = lipgloss.Color("#6244b7")
	Purple400 = lipgloss.Color("#4b30ab")
	Purple500 = lipgloss.Color("#3d2896")
	Purple600 = lipgloss.Color("#2f2081")
	Purple700 = lipgloss.Color("#21186c")
	Purple800 = lipgloss.Color("#131057")
	Purple900 = lipgloss.Color("#050842")
)

// Green is the complementary color to purple.
// Shades range from lightest (Green50) to darkest (Green900).
var (
	Green50  = lipgloss.Color("#80dba9")
	Green100 = lipgloss.Color("#6ccf96")
	Green200 = lipgloss.Color("#58c383")
	Green300 = lipgloss.Color("#44b770")
	Green400 = lipgloss.Color("#30ab5d")
	Green500 = lipgloss.Color("#28964e")
	Green600 = lipgloss.Color("#20813f")
	Green700 = lipgloss.Color("#186c30")
	Green800 = lipgloss.Color("#105721")
	Green900 = lipgloss.Color("#084212")
)

// Orange is a triadic color to purple.
// Shades range from lightest (Orange50) to darkest (Orange900).
var (
	Orange50  = lipgloss.Color("#dba980")
	Orange100 = lipgloss.Color("#cf966c")
	Orange200 = lipgloss.Color("#c38358")
	Orange300 = lipgloss.Color("#b77044")
	Orange400 = lipgloss.Color("#ab5d30")
	Orange500 = lipgloss.Color("#964e28")
	Orange600 = lipgloss.Color("#813f20")
	Orange700 = lipgloss.Color("#6c3018")
	Orange800 = lipgloss.Color("#572110")
	Orange900 = lipgloss.Color("#421208")
)

// Red is used for error states.
// Shades range from lightest (Red50) to darkest (Red900).
var (
	Red50  = lipgloss.Color("#db8080")
	Red100 = lipgloss.Color("#cf6c6c")
	Red200 = lipgloss.Color("#c35858")
	Red300 = lipgloss.Color("#b74444")
	Red400 = lipgloss.Color("#ab3030")
	Red500 = lipgloss.Color("#962828")
	Red600 = lipgloss.Color("#812020")
	Red700 = lipgloss.Color("#6c1818")
	Red800 = lipgloss.Color("#571010")
	Red900 = lipgloss.Color("#420808")
)

// Blue is an analogous color to purple.
// Shades range from lightest (Blue50) to darkest (Blue900).
var (
	Blue50  = lipgloss.Color("#80a9db")
	Blue100 = lipgloss.Color("#6c96cf")
	Blue200 = lipgloss.Color("#5883c3")
	Blue300 = lipgloss.Color("#4470b7")
	Blue400 = lipgloss.Color("#305dab")
	Blue500 = lipgloss.Color("#284e96")
	Blue600 = lipgloss.Color("#203f81")
	Blue700 = lipgloss.Color("#18306c")
	Blue800 = lipgloss.Color("#102157")
	Blue900 = lipgloss.Color("#081242")
)

// ANSI standard colors (0-7).
var (
	Black   = lipgloss.Color("0")
	Red     = lipgloss.Color("1")
	Green   = lipgloss.Color("2")
	Yellow  = lipgloss.Color("3")
	Blue    = lipgloss.Color("4")
	Magenta = lipgloss.Color("5")
	Cyan    = lipgloss.Color("6")
	White   = lipgloss.Color("7")
)

// ANSI bright colors (8-15).
var (
	BrightBlack   = lipgloss.Color("8")
	BrightRed     = lipgloss.Color("9")
	BrightGreen   = lipgloss.Color("10")
	BrightYellow  = lipgloss.Color("11")
	BrightBlue    = lipgloss.Color("12")
	BrightMagenta = lipgloss.Color("13")
	BrightCyan    = lipgloss.Color("14")
	BrightWhite   = lipgloss.Color("15")
)
