package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	db "github.com/lundjrl/go-bubble-tea-playground/shared/database"
)

type mainModel struct {
	currentTab int
	state      sessionState
	table      table.Model
	textInput  textinput.Model
	err        error
}

// sessionState to track which model is focused.
type sessionState uint

const (
	tableView sessionState = iota
	inputView
	welcomeView
)

type Theme struct {
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

var theme = Theme{
	blue:     lipgloss.Color("#89b4fa"),
	pink:     lipgloss.Color("#f5c2e7"),
	yellow:   lipgloss.Color("#f9e2af"),
	lavender: lipgloss.Color("#b4befe"),
	bg:       lipgloss.Color("#11111b"),
	fg:       lipgloss.Color("#cdd6f4")}

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
)

func newModel() mainModel {
	m := mainModel{state: tableView}

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 15},
		{Title: "Count", Width: 24},
	}

	var items []db.GroceryItem
	result := db.DBConn.Find(&items)

	if result.Error != nil {
		panic(result.Error)
	}

	tableRows := []table.Row{}

	for _, item := range items {
		row := []string{fmt.Sprint(item.ID), item.Name, fmt.Sprint(item.Count)}
		tableRows = append(tableRows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(49),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.fg).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(theme.bg).
		Background(theme.yellow).
		Bold(true)
	t.SetStyles(s)

	m.table = t
	m.textInput = textinput.New()
	m.textInput.Placeholder = "add an item?"
	m.textInput.CharLimit = 156
	m.textInput.Width = 49
	m.err = nil
	m.state = tableView
	m.currentTab = 0

	return m
}

func getTabUI(m mainModel) string {
	gap := tabGap.Render(strings.Repeat(" ", max(0, 98)))

	switch m.currentTab {
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

func getWelcomeUI(m mainModel) string {
	if m.currentTab != 0 {
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
		MarginLeft(6).Render("\n\t- Check your inventory\n\t- Make a grocery list\n\t- Support James")

	// itemsStyle := lipgloss.NewStyle().Foreground(theme.fg).PaddingTop(-5).
	// 	MarginTop(0).
	// 	MarginLeft(8).

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Center, titleStyle, descriptionStyle), line)
}

func getTableUI(m mainModel) string {
	if m.currentTab != 1 {
		return ""
	}

	focusedTable := lipgloss.JoinHorizontal(lipgloss.Top, focusedTableStyle.Render(m.table.View()), modelStyle.Render(m.textInput.View())+"\n")
	unfocusedTable := lipgloss.JoinHorizontal(lipgloss.Top, baseTableStyle.Render(m.table.View()), focusedModelStyle.Render(m.textInput.View())+"\n")

	if m.state == tableView {
		return focusedTable
	}
	return unfocusedTable
}

func getInputUI(m mainModel) string {
	if m.currentTab != 1 {
		return ""
	}

	tableHelperText := tipContainerStyle.Render("tab: focus next • enter: create new item • q: exit")
	inputHelperText := tipContainerStyle.Render("tab: focus next • enter: view entry • q: exit")
	focusedInput := lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().PaddingTop(1).Render(), tableHelperText)
	unfocusedInput := lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().PaddingTop(1).Render(), inputHelperText)

	if m.state == tableView {
		return unfocusedInput
	}
	return focusedInput
}

func (m mainModel) View() string {
	var s string = getTabUI(m)

	// Tab 1 UI
	s += getWelcomeUI(m)

	// Tab 2 UI
	s += getTableUI(m)
	s += getInputUI(m)

	// Tab 3 UI

	// Tab 4 UI

	return s
}

// Add initial actions on mount.
func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.textInput.Focus(), textinput.Blink) // no batch?
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == inputView {
				item := m.textInput.Value()
				db.CreateGroceryItem(item)
				rows := m.table.Rows()
				id := len(m.table.Rows()) + 1
				row := []string{fmt.Sprint(id), item, fmt.Sprint(1)}
				rows = append(rows, row)
				m.table.SetRows(rows)
				m.table.GotoBottom()
				m.textInput.Reset()
				m.textInput.Cursor.SetMode(cursor.New().Mode())
			}
		case "tab":
			if m.state == tableView {
				m.state = inputView
				m.table.Blur()
				m.textInput.Focus()
			} else {
				m.state = tableView
				m.textInput.Blur()
				m.table.Focus()
			}

		case "shift+tab":
			switch m.currentTab {
			case 0:
				m.currentTab = 1
			case 1:
				m.currentTab = 2
			case 2:
				m.currentTab = 3
			case 3:
				m.currentTab = 0
			default:
				m.currentTab = 0
			}
		}

		switch m.state {
		// update whichever model is focused
		case inputView:
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
			cmds = append(cmds, textinput.Blink)
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func parseCommand(command string) (tea.Model, error) {
	switch command {
	case "init":
		model, err := tea.NewProgram(newModel()).Run()
		return model, err
	case "help":
		model, err := tea.NewProgram(newModel()).Run()
		return model, err
	default:
		model, err := tea.NewProgram(newModel()).Run()
		return model, err
	}
}

func main() {
	log.Info("Starting application...")

	db.InitDatabaseConnection()

	argsAfterCommandName := os.Args[1:]

	//	if len(argsAfterCommandName) == 0 {
	if false {
		log.Error("Please invoke with a command. \n\n\t`$ go run main.go <command>`\n")
		os.Exit(1)
	}

	for _, element := range argsAfterCommandName {
		_, err := parseCommand(element)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	log.Info("Program terminated.")
}
