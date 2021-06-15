package models

import (
	"github.com/vinshop/apigen/pkg/util"
)

type ModelField struct {
	Name         string
	Type         string
	ColumnName   string
	IsPrimaryKey bool
	IsNullable   bool
}

func (f *ModelField) DTOField() *DTOField {

	var t string

	switch f.Type {
	case "time.Time":
		t = "int64"
	case "*time.Time":
		t = "*int64"
	default:
		t = f.Type
	}

	return &DTOField{
		Name:       f.Name,
		Type:       t,
		ColumnName: f.ColumnName,
		IsNullable: f.IsNullable,
	}
}

type Model struct {
	Pkg    string
	Module string
	Name   string
	Table  string
	Fields []*ModelField

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

func (m *Model) DTO() *DTO {
	fields := make([]*DTOField, len(m.Fields))

	for i := range m.Fields {
		fields[i] = m.Fields[i].DTOField()
	}

	return &DTO{
		Pkg:            m.Module + "/dtos",
		Module:         m.Module,
		Name:           m.Name,
		Table:          m.Table,
		Fields:         fields,
		NeedImport:     m.NeedImport,
		NeedImportTime: m.NeedImportTime,
		NeedImportGorm: m.NeedImportGorm,
	}
}
