package models

type Mapper struct {
	Name   string
	Pkg    string
	Module string
	Model  *Model
	DTO    *DTO
	Import map[string]bool
}

func NewMapper(model *Model, dto *DTO) *Mapper {

	mp := make(map[string]bool)

	for k := range model.Import {
		mp[k] = true
	}
	for k := range dto.Import {
		mp[k] = true
	}

	mp[model.Pkg] = true
	mp[dto.Pkg] = true

	return &Mapper{
		Name:   model.Name + "Mapper",
		Pkg:    model.Module + "/dtos",
		Module: model.Module,
		Model:  model,
		DTO:    dto,
		Import: mp,
	}
}
