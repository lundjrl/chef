package database

import (
	"errors"
	"strings"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type GroceryItem struct {
	gorm.Model
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// func GetGroceryItems() (*sql.rows, error){
// 	db := DBConn
// 	var items []GroceryItem
// 	result := db.Find(&items)
// 	return result.Rows()
// }

func GetGroceryItemByName(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("Please type a grocery item.")
	}

	db := DBConn

	var item GroceryItem
	result := db.Find(&item, "Name = ?", name)
	return result.Name(), result.Error
}

func CreateGroceryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("Please type a grocery item.")
	}

	db := DBConn
	item := new(GroceryItem)
	item.Name = name
	item.Count = 1

	result := db.Create(&item)
	log.Info("Created ::", item)

	return "", result.Error
}

func DeleteGroceryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)
	db := DBConn

	if len(name) <= 0 {
		return "", errors.New("Please type a grocery item.")
	}

	var item GroceryItem
	db.First(&item, "Name = ?", name)

	if item.Name == "" {
		return "", errors.New("There's no grocery item with that name.")
	}

	result := db.Delete(&item)
	log.Info("Item removed.")

	return "Item removed.", result.Error
}
