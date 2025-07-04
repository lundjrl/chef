package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
}

type statusMsg int

var UserInput string

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func saveItem(m model) {
	// TODO: Would save info into db.
	m.textInput.SetValue("")
	return
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			saveItem(m)
			return m, tea.SetWindowTitle("Added item!")
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	UserInput = m.textInput.Value()
	return m, cmd
}

func (m model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		log.Error(m.err)
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	return fmt.Sprintf(
		"What do you want to add to the list?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "vegetables?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func parseCommand(command string) (tea.Model, error) {
	log.Info("Command:: " + command)
	switch command {
	case "add":
		model, err := tea.NewProgram(initialModel()).Run()
		tea.ClearScreen()
		return model, err
	case "list":
		model, err := tea.NewProgram(initialModel()).Run()
		tea.ClearScreen()
		return model, err
	case "write":
		model, err := tea.NewProgram(initialModel()).Run()
		tea.ClearScreen()
		return model, err
	default:
		model, err := tea.NewProgram(model{}).Run()
		tea.ClearScreen()
		return model, err
	}
}

func main() {
	log.Info("Starting application...")

	argsAfterCommandName := os.Args[1:]

	fmt.Println(len(argsAfterCommandName))

	if len(argsAfterCommandName) == 0 {
		log.Error("Please invoke with a command. \n\n\t`$ go run main.go <command>`\n")
		os.Exit(1)
	}

	_, err := parseCommand(argsAfterCommandName[0])

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("User Input::" + UserInput)

	log.Info("Program terminated.")
}
