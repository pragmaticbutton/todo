package todo

type ToDoService interface {
	GetCategory(id string) string
}
