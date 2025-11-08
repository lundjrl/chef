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

type CheckboxItem struct {
	gorm.Model
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}

func GetGroceryItemByName(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	db := DBConn

	var item GroceryItem
	result := db.Find(&item, "Name = ?", name)
	return result.Name(), result.Error
}

func CreateGroceryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	db := DBConn
	item := new(GroceryItem)
	item.Name = name
	item.Count = 1

	result := db.Create(&item)

	return "", result.Error
}

func UpdateGroceryItem(id string, count int) (string, error) {
	db := DBConn

	var item GroceryItem
	db.First(&item, "Id")

	item.Count = count
	result := db.Save(&item)

	return "Item updated.", result.Error
}

func DeleteGroceryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)
	db := DBConn

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	var item GroceryItem
	db.First(&item, "Name = ?", name)

	if item.Name == "" {
		return "", errors.New("there's no grocery item with that name")
	}

	result := db.Delete(&item)
	log.Info("Item removed.")

	return "Item removed.", result.Error
}
