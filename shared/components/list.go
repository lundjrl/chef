package components

import (
	"github.com/charmbracelet/bubbles/list"
	db "github.com/lundjrl/chef/shared/database"
)

const listHeight = 14

// TODO: Move to UI

type listItem struct {
	checked     bool
	title, desc string
}

func Init() {
	groceryItems := db.GetGroceryItems()

	items := []list.Item{
		listItem{desc: "", title: "", checked: false},
	}

	for _, grocery := range groceryItems {
		record := listItem{checked: grocery.Checked, desc: "", title: grocery.Name}

		items = append(items, record)
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

}
