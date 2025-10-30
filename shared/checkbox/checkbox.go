package checkbox

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

func ItemRow(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[X] " + label)
	}

	return fmt.Sprintf("[ ] %s", label)
}
