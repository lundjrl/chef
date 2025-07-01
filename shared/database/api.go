package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabaseConnection() {
	var err error

	DBConn, err = gorm.Open(sqlite.Open("app.db"))

	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Database connection started")

	DBConn.AutoMigrate(&GroceryItem{})

	fmt.Println("Database Migrated")
}
