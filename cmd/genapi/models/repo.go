package models

import "github.com/vinshop/apigen/pkg/util"

type Repo struct {
	Pkg     string
	Module  string
	Model   *Model
	Name    string
	ImpName string
}

func (r *Repo) Service() *Service {
	return &Service{
		Module:  r.Module,
		Pkg:     r.Module + "/services",
		Repo:    r,
		Name:    r.Model.Name + "Service",
		ImpName: util.LowerTitle(r.Model.Name) + "Service",
	}
}
