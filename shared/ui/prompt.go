package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

var input string

func Create() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Write something down...").
				Value(&input).
				Validate(func(str string) error {
					if str == "Stop" {
						// TODO: Some stop handler here?
						return errors.New("Stopping...")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
