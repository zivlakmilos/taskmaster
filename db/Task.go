package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id            string
	Description   string
	DateGet       DateTime `db:"dateGet"`
	DateStart     DateTime `db:"dateStart"`
	DateCompleted DateTime `db:"dateCompleted"`
	Priority      int
	Progress      int
	Note          string
	Status        TaskStatus
	Project       *Project
	ProjectId     string `db:"projectId"`
}

type TaskStore struct {
	con *sqlx.DB
}

func NewTask() *Task {
	return &Task{
		Status:   TaskStatusTodo,
		Priority: 3,
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
      dateCompleted,
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
      :dateCompleted,
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
      dateCompleted=:dateCompleted,
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

func (s *TaskStore) GetAllByProjectId(projectId string) ([]*Task, error) {
	var res []*Task

	err := s.con.Select(&res, "SELECT * FROM Task WHERE projectId=?", projectId)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *TaskStore) GetAllByProjectName(projectName string) ([]*Task, error) {
	var res []*Task

	err := s.con.Select(&res, "SELECT Task.* FROM Task RIGHT JOIN Project ON Task.projectId=Project.id WHERE Project.name=?", projectName)
	if err != nil {
		return res, err
	}

	return res, nil
}
