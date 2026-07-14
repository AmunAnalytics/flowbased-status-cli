package config

import "charm.land/lipgloss/v2"

var (
	CellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(14)
	RowStyle = CellStyle.Foreground(lipgloss.White).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true)
	BorderStyle = lipgloss.NewStyle().Foreground(lipgloss.White)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f59e0b")).
			Bold(true).
			Align(lipgloss.Left)
	Green  = "#65a30d"
	Yellow = "#fbbf24"
	Red    = "#b91c1c"
)
