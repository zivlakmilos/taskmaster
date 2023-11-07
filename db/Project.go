package db

import (
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

func NewProjectStore(con *sqlx.DB) *ProjectStore {
	return &ProjectStore{
		con: con,
	}
}

func (s *ProjectStore) Insert(project *Project) error {
	_, err := s.con.Exec("INSERT INTO Project(id, name, status) values(?, ?, ?)", project)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) Update(project *Project) error {
	_, err := s.con.Exec("UPDATE Project SET name=?, status=? WHERE id=?", project)
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

func (s *ProjectStore) Delete(project *Project) error {
	_, err := s.con.Exec("DELETE FROM Project WHERE id=?", project)
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
