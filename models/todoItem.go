package models

type TodoItem struct {
	Id    string `json:"id"`
	Group string `json:"group"`
	Task  string `json:"task"`
}

func NewTodoItem(group string, task string) *TodoItem {
	var item TodoItem
	item.Id = "testid"
	item.Group = group
	item.Task = task

	return &item
}
