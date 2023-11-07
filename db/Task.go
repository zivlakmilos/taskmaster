package db

import "time"

type Task struct {
	Id            string
	Description   string
	DateGet       time.Time
	DateStart     time.Time
	DateCompleted time.Time
	Priority      int
	Progress      int
	Node          string
	Status        TaskStatus
	Project       *Project
	ProjectId     string
}
