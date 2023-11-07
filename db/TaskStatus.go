package db

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress            = "inProgress"
	TaskStatusCompleted             = "completed"
)
