package models

type Mapper struct {
	Name   string
	Pkg    string
	Module string
	Model  *Model
	DTO    *DTO
}

func NewMapper(model *Model, dto *DTO) *Mapper {
	return &Mapper{
		Name:   model.Name + "Mapper",
		Pkg:    model.Module + "/dtos",
		Module: model.Module,
		Model:  model,
		DTO:    dto,
	}
}
