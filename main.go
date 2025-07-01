package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	db "github.com/lundjrl/go-bubble-tea-playground/shared/database"
	table "github.com/lundjrl/go-bubble-tea-playground/shared/table"
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
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

		log.Info(UserInput)
		db.CreateGroceryItem(UserInput)

		return model, err
	case "list":
		model, err := table.Main()
		return model, err
	case "remove":
		model, err := tea.NewProgram(initialModel()).Run()

		db.DeleteGroceryItem(UserInput)

		return model, err
	case "write":
		model, err := tea.NewProgram(initialModel()).Run()

		db.CreateGroceryItem(UserInput)

		return model, err
	default:
		model, err := tea.NewProgram(model{}).Run()
		return model, err
	}
}

func main() {
	log.Info("Starting application...")

	db.InitDatabaseConnection()

	argsAfterCommandName := os.Args[1:]

	fmt.Println(len(argsAfterCommandName))

	if len(argsAfterCommandName) == 0 {
		log.Error("Please invoke with a command. \n\n\t`$ go run main.go <command>`\n")
		os.Exit(1)
	}

	for _, element := range argsAfterCommandName {
		// for i := 1; i < 1; i++ {
		_, err := parseCommand(element)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	log.Info("User Input::" + UserInput)

	log.Info("Program terminated.")
}
