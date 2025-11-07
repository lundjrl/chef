package tests

import (
	"testing"

	model "github.com/lundjrl/go-bubble-tea-playground/shared/database"
)

func TestGetGroceryItemByName(t *testing.T) {
	model.InitDatabaseConnection()

	val, err := model.GetGroceryItemByName("avocado")

	if err != nil {
		t.Log(err)
		// t.Error(err.Error())
	}

	if len(val) < 1 {
		t.Error("No Item was Retrieved.")
	}

}
