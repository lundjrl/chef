package database

import (
	"errors"
	"strings"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type InventoryItem struct {
	gorm.Model
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetInventoryItems() []InventoryItem {
	var items []InventoryItem
	result := DBConn.Find(&items)

	if result.Error != nil {
		panic(result.Error)
	}

	return items
}

func GetInventoryItemByName(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	db := DBConn

	var item InventoryItem
	result := db.Find(&item, "Name = ?", name)
	return result.Name(), result.Error
}

func CreateInventoryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	db := DBConn
	item := new(InventoryItem)
	item.Name = name
	item.Count = 1

	result := db.Create(&item)

	return "", result.Error
}

func UpdateInventoryItem(id string, count int) (string, error) {
	db := DBConn

	var item InventoryItem
	db.Find(&item, "ID = ?", id)

	item.Count = count
	result := db.Save(&item)

	return "Item updated.", result.Error
}

func DeleteInventoryItem(itemName string) (string, error) {
	name := strings.ToLower(itemName)
	db := DBConn

	if len(name) <= 0 {
		return "", errors.New("please type a grocery item")
	}

	var item InventoryItem
	db.First(&item, "Name = ?", name)

	if item.Name == "" {
		return "", errors.New("there's no grocery item with that name")
	}

	result := db.Delete(&item)
	log.Info("Item removed.")

	return "Item removed.", result.Error
}
