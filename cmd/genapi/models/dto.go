package models

type DTOField struct {
	Name         string
	Type         string
	ColumnName   string
	IsNullable   bool
}

type DTO struct {
	Pkg    string
	Module string
	Name   string
	Table  string
	Fields []*DTOField

	NeedImport     bool
	NeedImportTime bool
	NeedImportGorm bool
}
