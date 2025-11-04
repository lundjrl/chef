package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/lundjrl/chef/shared"
)

type Palette struct {
	blue     lipgloss.Color
	pink     lipgloss.Color
	yellow   lipgloss.Color
	lavender lipgloss.Color
	bg       lipgloss.Color
	fg       lipgloss.Color
}

// blue #89b4fa
// pink #f5c2e7
// yellow #f9e2af
// lavender #b4befe
// bg #11111b
// fg #cdd6f4

var theme = Palette{
	blue:     lipgloss.Color("#89b4fa"),
	pink:     lipgloss.Color("#f5c2e7"),
	yellow:   lipgloss.Color("#f9e2af"),
	lavender: lipgloss.Color("#b4befe"),
	bg:       lipgloss.Color("#11111b"),
	fg:       lipgloss.Color("#cdd6f4"),
}

var (
	modelStyle = lipgloss.NewStyle().
			Width(49).
			Height(2).
			BorderStyle(lipgloss.HiddenBorder()).
			MarginLeft(1).MarginTop(1)
	focusedModelStyle = lipgloss.NewStyle().
				Width(49).
				Height(2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(theme.pink).
				MarginLeft(1).MarginTop(1)
	tipContainerStyle = lipgloss.NewStyle().Foreground(theme.fg).Border(lipgloss.RoundedBorder()).BorderForeground(theme.yellow).MarginTop(1).MarginBottom(2).Width(100)
	baseTableStyle    = lipgloss.NewStyle().
				BorderStyle(lipgloss.HiddenBorder()).
				Width(49).Height(5).MarginTop(1)
	focusedTableStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(theme.pink).
				Width(49).Height(5).MarginTop(1)

	tabContainer = lipgloss.NewStyle().Render()

	horizontalRule = lipgloss.NewStyle().Render()

	highlight = lipgloss.NewStyle().Foreground(theme.pink)

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(theme.pink).
		Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true).Background(theme.blue)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	TableHeaderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(theme.fg).
				BorderBottom(true).
				Bold(false)

	TableSelectedStyle = lipgloss.NewStyle().Foreground(theme.bg).
				Background(theme.yellow).
				Bold(true)
)

func GetTabUI(m shared.MainModel) string {
	gap := tabGap.Render(strings.Repeat(" ", max(0, 98)))

	switch m.CurrentTab {
	case 1:
		return lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			tab.Render("(h) Home"),
			activeTab.Render("(i) Inventory"),
			tab.Render("(g) Grocery List"),
			tab.Render("(s) Settings"),
			gap)
	case 2:
		return lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			tab.Render("(h) Home"),
			tab.Render("(i) Inventory"),
			activeTab.Render("(g) Grocery List"),
			tab.Render("(s) Settings"),
			gap)
	case 3:
		return lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			tab.Render("(h) Home"),
			tab.Render("(i) Inventory"),
			tab.Render("(g) Grocery List"),
			activeTab.Render("(s) Settings"),
			gap)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		activeTab.Render("(h) Home"),
		tab.Render("(i) Inventory"),
		tab.Render("(g) Grocery List"),
		tab.Render("(s) Settings"),
		gap)
}

func enumerateList(items list.Items, i int) string {
	return "\t\t\t ✓ "
}

func GetWelcomeUI(m shared.MainModel) string {
	if m.CurrentTab != 0 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		PaddingTop(2).
		MarginLeft(6).
		Height(2).
		Bold(true).Foreground(theme.blue).
		Render("Chef!")

	descriptionStyle := lipgloss.NewStyle().
		Bold(true).PaddingTop(3).
		Foreground(theme.lavender).
		MarginLeft(2).
		Render("Your new inventory cli tool")

	line := lipgloss.NewStyle().
		BorderForeground(theme.pink).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingTop(-1).
		Width(50).
		MarginLeft(6).Render()

	itemStyle := lipgloss.NewStyle().
		Foreground(theme.pink).
		TabWidth(4) // Tab width can be different per terminal

	listItems := list.New(
		"Check your inventory",
		"Make a grocery list",
		"Support James",
	).ItemStyle(itemStyle).Enumerator(enumerateList)

	spacer := lipgloss.NewStyle().
		Height(3).Render(" ")

	homeHelperText := tipContainerStyle.MarginLeft(6).Width(60).Padding(1).Render("i: go to inventory • g: go to list • s: go to settings")

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Center, titleStyle, descriptionStyle), line, listItems.String(), spacer, homeHelperText, spacer)
}

func GetTableUI(m shared.MainModel) string {
	if m.CurrentTab != 1 {
		return ""
	}

	focusedTable := lipgloss.JoinHorizontal(lipgloss.Top, focusedTableStyle.Render(m.Table.View()), modelStyle.Render(m.TextInput.View())+"\n")
	unfocusedTable := lipgloss.JoinHorizontal(lipgloss.Top, baseTableStyle.Render(m.Table.View()), focusedModelStyle.Render(m.TextInput.View())+"\n")

	if m.State == shared.TableView {
		return focusedTable
	}
	return unfocusedTable
}

func GetInputUI(m shared.MainModel) string {
	if m.CurrentTab != 1 {
		return ""
	}

	tableHelperText := tipContainerStyle.Render("tab: focus next • enter: create new item • q: exit")
	inputHelperText := tipContainerStyle.Render("tab: focus next • enter: view entry • q: exit")
	focusedInput := lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().PaddingTop(1).Render(), tableHelperText)
	unfocusedInput := lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().PaddingTop(1).Render(), inputHelperText)

	if m.State == shared.TableView {
		return unfocusedInput
	}
	return focusedInput
}

func GetListUI(m shared.MainModel) string {
	if m.CurrentTab != 2 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		PaddingTop(2).
		MarginLeft(6).
		Height(2).
		Bold(true).Foreground(theme.blue).
		Render("Grocery List")

	descriptionStyle := lipgloss.NewStyle().
		Bold(true).PaddingTop(3).
		Foreground(theme.lavender).
		MarginLeft(2).
		Render("Coming soon...")

	line := lipgloss.NewStyle().
		BorderForeground(theme.pink).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingTop(-1).
		Width(50).
		MarginLeft(6).Render()

	spacer := lipgloss.NewStyle().
		Height(3).Render(" ")

	homeHelperText := tipContainerStyle.MarginLeft(6).Width(60).Padding(1).Render("i: go to inventory • g: go to list • s: go to settings")

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Center, titleStyle, descriptionStyle), line, spacer, homeHelperText, spacer)
}

func GetSettingsUI(m shared.MainModel) string {
	if m.CurrentTab != 3 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		PaddingTop(2).
		MarginLeft(6).
		Height(2).
		Bold(true).Foreground(theme.blue).
		Render("Chef!")

	descriptionStyle := lipgloss.NewStyle().
		Bold(true).PaddingTop(3).
		Foreground(theme.lavender).
		MarginLeft(2).
		Render("Your new inventory cli tool")

	line := lipgloss.NewStyle().
		BorderForeground(theme.pink).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingTop(-1).
		Width(50).
		MarginLeft(6).Render()

	bodyStyle := lipgloss.NewStyle().
		Bold(true).PaddingTop(2).
		Foreground(theme.lavender).
		MarginLeft(6).
		Render("Settings coming soon...\n\nEventually you'll be able to wipe your list.")

	spacer := lipgloss.NewStyle().
		Height(3).Render(" ")

	homeHelperText := tipContainerStyle.MarginLeft(6).Width(60).Padding(1).Render("i: go to inventory • g: go to list • s: go to settings")

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Center, titleStyle, descriptionStyle), line, bodyStyle, spacer, homeHelperText, spacer)
}
