package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id            string
	Description   string
	DateGet       time.Time
	DateStart     time.Time
	DateCompleted time.Time
	Priority      int
	Progress      int
	Note          string
	Status        TaskStatus
	Project       *Project
	ProjectId     string
}

type TaskStore struct {
	con *sqlx.DB
}

func NewTask() *Task {
	return &Task{
		Status: TaskStatusTodo,
	}
}

func NewTaskStore(con *sqlx.DB) *TaskStore {
	return &TaskStore{
		con: con,
	}
}

func (s *TaskStore) Insert(task *Task) error {
	task.Id = uuid.New().String()
	_, err := s.con.NamedExec(`INSERT INTO Task(
      id,
      description,
      dateGet,
      dateStart,
      dateComplete,
      priority,
      progress,
      note,
      status,
      projectId
    ) VALUES(
      :id,
      :description,
      :dateGet,
      :dateStart,
      :dateComplete,
      :priority,
      :progress,
      :note,
      :status,
      :projectId
    )`, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskStore) Update(task *Task) error {
	_, err := s.con.NamedExec(`UPDATE Task SET
      description=:description,
      dateGet=:dateGet,
      dateStart=:dateStart,
      dateComplete=:dateComplete,
      priority=:priority,
      progress=:progress,
      note=:note,
      status=:status,
      projectId=:projectId
    WHERE id=:id`, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskStore) Save(task *Task) error {
	if task.Id == "" {
		return s.Insert(task)
	}

	return s.Update(task)
}

func (s *TaskStore) Delete(id string) error {
	_, err := s.con.Exec("DELETE FROM Task WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskStore) GetAll() ([]*Task, error) {
	var res []*Task

	err := s.con.Select(&res, "SELECT * FROM Task")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *TaskStore) Get(id string) (*Task, error) {
	var res *Task

	err := s.con.Get(res, "SELECT * FROM Task WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
