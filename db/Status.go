package db

type Status string

const (
	StatusActive    Status = "active"
	StatusPaused           = "paused"
	StatusCompleted        = "completed"
)
