package main

import (
	"godoBackend/database"
)

func main() {
	var items = database.GetAllItems()
	for _, item := range items {
		println(item.Group + ": " + item.Task)
	}
}
