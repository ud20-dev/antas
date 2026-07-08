package console

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

var BaseStyle = lipgloss.NewStyle().Background(lipgloss.Black)

var ErrorStyle = BaseStyle.Foreground(lipgloss.Red)

var InfoStyle = BaseStyle.Foreground(lipgloss.Color("#1d1d1d"))

var SuccessStyle = BaseStyle.Foreground(lipgloss.Green)

func PrintWithStyle(style lipgloss.Style, format string, v ...any) {
	formatted_string := fmt.Sprintf(format, v...)
	render_result := style.Render(formatted_string)
	lipgloss.Println(render_result)
}