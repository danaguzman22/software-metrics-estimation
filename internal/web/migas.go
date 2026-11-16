package web

import "strings"

// Miga es un nivel de la migaja de pan de la navegación secundaria (T-15).
// Actual marca el nivel vigente y habilita aria-current="page" (D9).
type Miga struct {
	Texto  string
	URL    string
	Actual bool
}

// migasDeProyecto antepone los niveles Proyectos y el proyecto indicado a
// los niveles de la pantalla y marca el último nivel como vigente.
func migasDeProyecto(nombreProyecto, proyectoID string, niveles ...Miga) []Miga {
	nombre := strings.TrimSpace(nombreProyecto)
	if nombre == "" {
		nombre = "Proyecto"
	}

	migas := make([]Miga, 0, len(niveles)+2)
	migas = append(migas,
		Miga{Texto: "Proyectos", URL: "/proyectos"},
		Miga{Texto: nombre, URL: "/proyectos/" + proyectoID},
	)
	migas = append(migas, niveles...)
	migas[len(migas)-1].Actual = true
	return migas
}
