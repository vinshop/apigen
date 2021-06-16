package models

type DTOField struct {
	Name       string
	NameLower  string
	Type       string
	ColumnName string
	IsNullable bool
	Cast       string
	FuncCall   string
	IsPointer  bool
	MType      string
}

type DTO struct {
	Pkg    string
	Module string
	Name   string
	Table  string
	Fields []*DTOField
	Import map[string]bool
}
