package models

import (
	"fmt"
	"github.com/vinshop/apigen/pkg/util"
)

type ModelField struct {
	Name         string
	NameLower    string
	Type         string
	ColumnName   string
	IsPrimaryKey bool
	IsNullable   bool
	IsPointer    bool
	DType        string
}

func (f *ModelField) DTOField() *DTOField {

	var t string

	var mType = "model." + f.Name
	var dType = "dto." + f.Name
	switch f.Type {
	case "time.Time":
		t = "int64"
		mType = "model." + f.Name + ".Unix()"
		dType = fmt.Sprintf("time.Unix(dto.%v,0)", f.Name)
	case "*time.Time":
		t = "*int64"
		mType = "model." + f.Name + ".Unix()"
		dType = fmt.Sprintf("time.Unix(*dto.%v,0)", f.Name)
	default:
		t = f.Type
	}

	f.DType = dType
	return &DTOField{
		Name:       f.Name,
		NameLower:  util.LowerTitle(f.Name),
		Type:       t,
		ColumnName: f.ColumnName,
		IsNullable: f.IsNullable,
		IsPointer:  f.IsPointer,
		MType:      mType,
	}
}

type Model struct {
	Pkg    string
	Module string
	Name   string
	Table  string
	Fields []*ModelField

	NeedImportGorm bool

	Import map[string]bool
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
		Pkg:    m.Module + "/dtos",
		Module: m.Module,
		Name:   m.Name,
		Table:  m.Table,
		Fields: fields,
		Import: m.Import,
	}
}
