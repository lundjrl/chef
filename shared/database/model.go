package database

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GroceryItem struct {
	gorm.Model
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetGroceryItems(c *fiber.Ctx) error {
	db := DBConn
	var items []GroceryItem
	db.Find(&items)
	return c.JSON(items)
}

func GetGroceryItemById(c *fiber.Ctx) error {
	id := c.Params("id")
	db := DBConn

	var item GroceryItem
	db.Find(&item, id)
	return c.JSON(item)
}

func CreateGroceryItem(c *fiber.Ctx) error {
	db := DBConn
	item := new(GroceryItem)

	if err := c.BodyParser(item); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&item)
	return c.JSON(item)
}

func DeleteGroceryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	db := DBConn

	var item GroceryItem
	db.First(&item, id)

	if item.Name == "" {
		return c.Status(500).SendString("No Item Found with ID")
	}
	db.Delete(&item)
	return c.SendString("Item Successfully deleted")
}
