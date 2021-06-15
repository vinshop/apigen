package models

import (
	"github.com/vinshop/apigen/pkg/util"
)

type Field struct {
	Name         string
	Type         string
	ColumnName   string
	IsPrimaryKey bool
	IsNullable   bool
}

type Model struct {
	Pkg    string
	Module string
	Name   string
	Table  string
	Fields []*Field

	NeedImport     bool
	NeedImportTime bool
	NeedImportGorm bool
}

func (m *Model) ToDbRepo() *Repo {
	return &Repo{
		Module:  m.Module,
		Pkg:     m.Module + "/repositories",
		Model:   m,
		Name:    m.Name + "Repository",
		ImpName: util.LowerTitle(m.Name) + "Repository",
	}
}

type Repo struct {
	Pkg     string
	Module  string
	Model   *Model
	Name    string
	ImpName string
}

func (r *Repo) ToService() *Service {
	return &Service{
		Module:  r.Module,
		Pkg:     r.Module + "/services",
		Repo:    r,
		Name:    r.Model.Name + "Service",
		ImpName: util.LowerTitle(r.Model.Name) + "Service",
	}
}

type Service struct {
	Pkg     string
	Module  string
	Repo    *Repo
	Name    string
	ImpName string
}

func (s *Service) ToController() *Controller {
	return &Controller{
		Pkg:     s.Pkg + "/controllers",
		Module:  s.Pkg,
		Service: s,
		Name:    s.Repo.Model.Name + "Controller",
		ImpName: util.LowerTitle(s.Repo.Model.Name) + "Controller",
	}
}

type Controller struct {
	Pkg     string
	Module  string
	Service *Service
	Name    string
	ImpName string
}
