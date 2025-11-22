package database

import (
	"gorm.io/gorm"
)

type GroceryItem struct {
	gorm.Model
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}

func GetGroceryItems() []GroceryItem {
	var items []GroceryItem
	result := DBConn.Find(&items)

	if result.Error != nil {
		panic(result.Error)
	}

	return items
}

func UpdateGroceryItem(id string, checked bool) (string, error) {
	db := DBConn

	var item GroceryItem
	db.Find(&item, "ID = ?", id)

	item.Checked = checked

	result := db.Save(&item)

	return "Item updated.", result.Error
}
