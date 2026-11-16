package web

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
	return nil
}
