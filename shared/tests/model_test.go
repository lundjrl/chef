package tests

import (
	"testing"

	db "github.com/lundjrl/chef/shared/database"
)

func TestGetInventoryItemByName(t *testing.T) {
	db.InitDatabaseConnection()

	val, err := db.GetInventoryItemByName("avocado")

	if err != nil {
		t.Log(err)
		// t.Error(err.Error())
	}

	if len(val) < 1 {
		t.Error("No Item was Retrieved.")
	}

}
