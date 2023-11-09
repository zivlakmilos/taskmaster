package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	Id     string
	Name   string
	Status Status
	Tasks  []*Task
}

type ProjectStore struct {
	con *sqlx.DB
}

func NewProject() *Project {
	return &Project{
		Status: StatusActive,
	}
}

func NewProjectStore(con *sqlx.DB) *ProjectStore {
	return &ProjectStore{
		con: con,
	}
}

func (s *ProjectStore) Insert(project *Project) error {
	project.Id = uuid.New().String()
	_, err := s.con.NamedExec("INSERT INTO Project(id, name, status) values(:id, :name, :status)", project)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) Update(project *Project) error {
	_, err := s.con.NamedExec("UPDATE Project SET name=:name, status=:status WHERE id=:id", project)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) Save(project *Project) error {
	if project.Id == "" {
		return s.Insert(project)
	}

	return s.Update(project)
}

func (s *ProjectStore) Delete(id string) error {
	_, err := s.con.Exec("DELETE FROM Project WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) DeleteByName(name string) error {
	_, err := s.con.Exec("DELETE FROM Project WHERE name=?", name)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) GetAll() ([]*Project, error) {
	var res []*Project

	err := s.con.Select(&res, "SELECT * FROM Project")
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *ProjectStore) Get(id string) (*Project, error) {
	var res *Project

	err := s.con.Get(res, "SELECT * FROM Project")
	if err != nil {
		return nil, err
	}

	return res, nil
}
