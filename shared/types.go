package shared

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

type MainModel struct {
	CurrentTab int
	State      SessionState
	Table      table.Model
	TextInput  textinput.Model
	Err        error
}

// sessionState to track which model is focused.
type SessionState uint

const (
	TableView SessionState = iota
	InputView
	WelcomeView
)
