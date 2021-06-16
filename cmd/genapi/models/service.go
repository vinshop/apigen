package models

import "github.com/vinshop/apigen/pkg/util"

type Service struct {
	Pkg     string
	Module  string
	Repo    *Repo
	Name    string
	ImpName string
}

func (s *Service) Controller() *Controller {
	return &Controller{
		Pkg:     s.Pkg + "/controllers",
		Module:  s.Pkg,
		Service: s,
		Name:    s.Repo.Model.Name + "Controller",
		ImpName: util.LowerTitle(s.Repo.Model.Name) + "Controller",
	}
}
