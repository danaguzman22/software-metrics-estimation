package service

import "github.com/danaguzman22/software-metrics-estimation/internal/domain"

type ProyectoRegistrado struct {
	ID       string
	Proyecto domain.Proyecto
}

type RepositorioProyectos interface {
}

type RepositorioIntegrantes interface {
}

type RepositorioBacklog interface {
}
