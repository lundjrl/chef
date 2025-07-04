package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	db "github.com/lundjrl/go-bubble-tea-playground/shared/database"
)

type (
	errMsg error
)

type mainModel struct {
	state     sessionState
	table     table.Model
	textInput textinput.Model
	err       error
	index     int
}

// sessionState to track which model is focused.
type sessionState uint

const (
	tableView sessionState = iota
	inputView
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("68")).
			MarginLeft(17)
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).PaddingLeft(2).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69")).MarginLeft(17)
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	baseTableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).Width(50).Height(5)
)

type statusMsg int

var UserInput string

func newModel() mainModel {
	m := mainModel{state: tableView}

	m.table = table.New()
	m.textInput = textinput.New()
	return m
}

// Add initial actions on mount.
func (m mainModel) Init() tea.Cmd {
	m.state = tableView
	return tea.Batch(m.textInput.Focus()) // no batch?
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == tableView {
				m.state = inputView
			} else {
				m.state = tableView
			}
		}

		switch m.state {
		// update whichever model is focused
		case inputView:
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
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

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == tableView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, baseTableStyle.Render(fmt.Sprintf("%4s", m.table.View()))+"\n")
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.textInput.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == inputView {
		return "textInput"
	}
	return "table"
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

	if len(argsAfterCommandName) == 0 {
		log.Error("Please invoke with a command. \n\n\t`$ go run main.go <command>`\n")
		os.Exit(1)
	}

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Count", Width: 10},
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
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	for _, element := range argsAfterCommandName {
		_, err := parseCommand(element)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	log.Info("Program terminated.")
}
