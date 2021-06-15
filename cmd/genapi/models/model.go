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

func (m *Model) Repo() *Repo {
	return &Repo{
		Module:  m.Module,
		Pkg:     m.Module + "/repositories",
		Model:   m,
		Name:    m.Name + "Repository",
		ImpName: util.LowerTitle(m.Name) + "Repository",
	}
}

