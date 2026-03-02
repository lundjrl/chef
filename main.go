package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/lundjrl/chef/shared"
	db "github.com/lundjrl/chef/shared/database"
	"github.com/lundjrl/chef/shared/ui"
)

type mainModel shared.MainModel

func initTable() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 15},
		{Title: "Count", Width: 24},
	}

	items := db.GetInventoryItems()

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
	s.Header = ui.TableHeaderStyle
	s.Selected = ui.TableSelectedStyle

	t.SetStyles(s)

	return t
}

// func initList() {
// 	items :=
// }

func newModel() mainModel {
	m := mainModel{State: shared.TableView}

	m.Table = initTable()
	// m.GroceryList = initList()
	m.TextInput = textinput.New()
	m.TextInput.Placeholder = "add an item?"
	m.TextInput.CharLimit = 156
	m.TextInput.Width = 49
	m.Err = nil
	m.State = shared.TableView
	m.CurrentTab = 0

	return m
}

func (m mainModel) View() string {
	var s string = ui.GetTabUI(shared.MainModel(m))

	// Tab 1 UI
	s += ui.GetWelcomeUI(shared.MainModel(m))

	// Tab 2 UI
	s += ui.GetTableUI(shared.MainModel(m))
	s += ui.GetInputUI(shared.MainModel(m))

	// Tab 3 UI
	s += ui.GetListUI(shared.MainModel(m))

	// Tab 4 UI
	s += ui.GetSettingsUI(shared.MainModel(m))

	return s
}

// Add initial actions on mount.
func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.TextInput.Focus(), textinput.Blink) // no batch?
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		count, err := strconv.Atoi(msg.String())

		if err == nil && m.State == shared.TableView {
			row := m.Table.SelectedRow()

			db.UpdateInventoryItem(row[0], count)
			rows := m.Table.Rows()
			row[2] = msg.String()

			index := m.Table.Cursor()
			rows[index] = row

			m.Table.SetRows(rows)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.State == shared.InputView {
				item := m.TextInput.Value()
				db.CreateInventoryItem(item)
				rows := m.Table.Rows()
				id := len(m.Table.Rows()) + 1
				row := []string{fmt.Sprint(id), item, fmt.Sprint(1)}
				rows = append(rows, row)
				m.Table.SetRows(rows)
				m.Table.GotoBottom()
				m.TextInput.Reset()
				m.TextInput.Cursor.SetMode(cursor.New().Mode())
			}

		case "tab":
			switch m.State {
			case shared.TableView:
				m.State = shared.InputView
				m.Table.Blur()
				m.TextInput.Focus()
			case shared.InputView:
				m.State = shared.TableView
				m.TextInput.Blur()
				m.Table.Focus()
			}

		case "shift+tab":
			switch m.CurrentTab {
			case 0:
				m.CurrentTab = 1
			case 1:
				m.CurrentTab = 2
			case 2:
				m.CurrentTab = 3
			case 3:
				m.CurrentTab = 0
			default:
				m.CurrentTab = 0
			}

		case "h":
			if !m.TextInput.Focused() {
				m.CurrentTab = 0
			}

		case "i":
			if !m.TextInput.Focused() {
				m.CurrentTab = 1
			}

		case "g":
			if !m.TextInput.Focused() {
				m.CurrentTab = 2
			}
		case "s":
			if !m.TextInput.Focused() {
				m.CurrentTab = 3
			}
		}

		switch m.State {
		// update whichever model is focused
		case shared.
		case shared.InputView:
			m.TextInput, cmd = m.TextInput.Update(msg)
			cmds = append(cmds, cmd)
			cmds = append(cmds, textinput.Blink)
		case shared.TableView:
			m.Table, cmd = m.Table.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.Table, cmd = m.Table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// Note: This is set up to add future commands.
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
