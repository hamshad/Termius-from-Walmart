package ui

import "github.com/charmbracelet/lipgloss"

// ── Color Palette ──────────────────────────────────────────────────────
// A cohesive cyberpunk-terminal palette: deep navy base, cyan/teal accents,
// warm highlights for actions, muted tones for secondary text.

var (
	// Primary brand colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // vibrant purple
	ColorAccent    = lipgloss.Color("#06B6D4") // cyan
	ColorAccentDim = lipgloss.Color("#0E7490") // darker cyan
	ColorHighlight = lipgloss.Color("#F59E0B") // amber/gold
	ColorSuccess   = lipgloss.Color("#10B981") // emerald green
	ColorDanger    = lipgloss.Color("#EF4444") // red
	ColorWarning   = lipgloss.Color("#F97316") // orange

	// Background tones
	ColorBg        = lipgloss.Color("#0F172A") // deep navy
	ColorBgPanel   = lipgloss.Color("#1E293B") // slate panel
	ColorBgActive  = lipgloss.Color("#334155") // active/focused bg
	ColorBgOverlay = lipgloss.Color("#1E1B4B") // modal/overlay bg

	// Text tones
	ColorText       = lipgloss.Color("#E2E8F0") // light slate
	ColorTextDim    = lipgloss.Color("#94A3B8") // muted
	ColorTextMuted  = lipgloss.Color("#64748B") // very muted
	ColorTextBright = lipgloss.Color("#F8FAFC") // near-white

	// Borders
	ColorBorder      = lipgloss.Color("#334155") // default border
	ColorBorderFocus = lipgloss.Color("#7C3AED") // focused border (primary)
)

// ── Layout Constants ───────────────────────────────────────────────────

const (
	PaddingH     = 2 // horizontal padding
	PaddingV     = 1 // vertical padding
	MinPaneWidth = 30
)

// ── App Chrome ─────────────────────────────────────────────────────────

var (
	// App header / branding
	AppHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Background(ColorBgPanel).
			Padding(0, PaddingH).
			MarginBottom(1)

	AppVersionStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Italic(true)

	// Status bar at bottom
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Background(ColorBgPanel).
			Padding(0, PaddingH)

	StatusBarKeyStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	StatusBarDescStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim)
)

// ── Title & Heading ────────────────────────────────────────────────────

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Italic(true).
			MarginLeft(1)

	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorBorder).
				MarginBottom(1).
				PaddingBottom(0)
)

// ── Panels & Boxes ─────────────────────────────────────────────────────

var (
	// Generic bordered panel
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(PaddingV, PaddingH)

	// Focused panel (e.g. active SFTP pane)
	PanelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorBorderFocus).
				Padding(PaddingV, PaddingH)

	// Dimmed panel (e.g. inactive SFTP pane)
	PanelDimStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder).
			Padding(PaddingV, PaddingH)

	// Modal overlay (e.g. confirmation dialog)
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(PaddingV, PaddingH+1).
			Background(ColorBgOverlay).
			Width(50)
)

// ── Messages & Feedback ────────────────────────────────────────────────

var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	MessageStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)
)

// ── Form Inputs ────────────────────────────────────────────────────────

var (
	InputLabelStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Width(8).
			Align(lipgloss.Right).
			MarginRight(1)

	InputActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorBorderFocus).
				Padding(0, 1).
				Width(46)

	InputInactiveStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder()).
				Padding(0, 1).
				Foreground(ColorTextDim).
				Width(46)

	SubmitButtonActiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextBright).
				Background(ColorPrimary).
				Bold(true).
				Padding(0, 3).
				MarginTop(1)

	SubmitButtonInactiveStyle = lipgloss.NewStyle().
					Foreground(ColorTextDim).
					Background(ColorBgActive).
					Padding(0, 3).
					MarginTop(1)
)

// ── List Items ─────────────────────────────────────────────────────────

var (
	// Server list item
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorTextBright).
				Background(ColorBgActive).
				Bold(true).
				Padding(0, 1)

	ListItemDescStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim)

	ListItemDescSelectedStyle = lipgloss.NewStyle().
					Foreground(ColorTextMuted).
					Background(ColorBgActive)

	// File picker items
	FileItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorText)

	FileSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorAccent).
				Bold(true)

	FileDirStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorHighlight)

	FileDirSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorHighlight).
				Bold(true)
)

// ── SFTP Specific ──────────────────────────────────────────────────────

var (
	SFTPHeaderStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Align(lipgloss.Center)

	SFTPHeaderDimStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Align(lipgloss.Center)

	SFTPPathStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Italic(true).
			Align(lipgloss.Center)

	SFTPDividerStyle = lipgloss.NewStyle().
				Foreground(ColorBorder)

	TransferBarStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)
)

// ── PEM Editor ─────────────────────────────────────────────────────────

var (
	PemBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccentDim).
			Padding(1, 2).
			Foreground(ColorText)

	PemLineNumStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Width(4).
			Align(lipgloss.Right).
			MarginRight(1)

	PemContentStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)
)

// ── Confirmation Dialog ────────────────────────────────────────────────

var (
	ConfirmTitleStyle = lipgloss.NewStyle().
				Foreground(ColorWarning).
				Bold(true).
				MarginBottom(1)

	ConfirmMessageStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				MarginBottom(1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextBright).
				Background(ColorPrimary).
				Bold(true).
				Padding(0, 2)

	ButtonDangerActiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextBright).
				Background(ColorDanger).
				Bold(true).
				Padding(0, 2)

	ButtonInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Background(ColorBgActive).
				Padding(0, 2)
)

// ── Menu ───────────────────────────────────────────────────────────────

var (
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			PaddingLeft(2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				PaddingLeft(1)

	MenuIconStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			MarginRight(1)
)

// ── Helper: Keybinding rendering ───────────────────────────────────────

// KeyStyle renders a keyboard shortcut key like [q], [enter], etc.
var KeyStyle = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Bold(true)

// KeyDescStyle renders the description next to a key.
var KeyDescStyle = lipgloss.NewStyle().
	Foreground(ColorTextDim)

// RenderKey creates a formatted key hint like "q quit"
func RenderKey(key, desc string) string {
	return KeyStyle.Render(key) + " " + KeyDescStyle.Render(desc)
}

// RenderKeyBar creates a horizontal bar of key hints
func RenderKeyBar(pairs ...string) string {
	var parts []string
	for i := 0; i < len(pairs)-1; i += 2 {
		parts = append(parts, RenderKey(pairs[i], pairs[i+1]))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, joinWith(parts, "  "))
}

func joinWith(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

// ── Helper: Responsive width ───────────────────────────────────────────

// ClampWidth ensures a value stays within min/max bounds.
func ClampWidth(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
