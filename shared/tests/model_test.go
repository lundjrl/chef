package tests

import (
	"testing"

	model "github.com/lundjrl/chef/shared/database"
)

func TestGetInventoryItemByName(t *testing.T) {
	model.InitDatabaseConnection()

	val, err := model.GetInventoryItemByName("avocado")

	if err != nil {
		t.Log(err)
		// t.Error(err.Error())
	}

	if len(val) < 1 {
		t.Error("No Item was Retrieved.")
	}

}
