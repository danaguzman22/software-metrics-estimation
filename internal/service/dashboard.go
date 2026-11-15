package service

import "errors"

// PuntoVelocidadSprint es la velocidad (Story Points completados) de un
// Sprint cerrado del proyecto (CA2).
type PuntoVelocidadSprint struct {
	SprintNumero int
	Velocidad    int
}

// PuntoHorasSprint son las horas estimadas y reales de un Sprint cerrado.
// MetricasDisponibles es false cuando el Sprint no tiene snapshot de esfuerzo
// (cerrado antes de US-17): en ese caso HorasEstimadas y HorasReales son 0 y
// el dashboard se muestra igualmente (D12).
type PuntoHorasSprint struct {
	SprintNumero          int
	HorasEstimadas        float64
	HorasReales           float64
	EstimacionesCompletas bool
	MetricasDisponibles   bool
}

// DashboardProyecto es el "modelo de lectura" del dashboard: lo arma
// ServicioDashboard reutilizando los servicios existentes. Crece de forma
// append-only entre historias (D19).
type DashboardProyecto struct {
	SPPlanificadosTotal            int
	SPCompletadosTotal             int
	VelocidadPromedio              float64
	PorcentajeHistoriasCompletadas float64
	HorasEstimadasTotal            float64
	HorasRealesTotal               float64
	EstimacionesCompletas          bool
	DesviacionHorasTotal           *float64
	SprintsCerrados                int
	VelocidadPorSprint             []PuntoVelocidadSprint
	HorasPorSprint                 []PuntoHorasSprint
}

// ServicioDashboard arma los datos del dashboard del proyecto.
type ServicioDashboard struct {
	proyectos RepositorioProyectos
	sprints   *ServicioSprints
	esfuerzo  *ServicioEsfuerzo
}

// NuevoServicioDashboard crea el servicio con los servicios existentes.
func NuevoServicioDashboard(
	proyectos RepositorioProyectos,
	sprints *ServicioSprints,
	esfuerzo *ServicioEsfuerzo,
) *ServicioDashboard {
	return &ServicioDashboard{proyectos: proyectos, sprints: sprints, esfuerzo: esfuerzo}
}

// ObtenerDashboard devuelve las métricas clave y las listas por Sprint del
// dashboard. Devuelve domain.ErrProyectoNoEncontrado sin envolver si el
// proyecto no existe; el resto de los errores de repositorio se envuelven con
// %w.
func (s *ServicioDashboard) ObtenerDashboard(proyectoID string) (DashboardProyecto, error) {
	return DashboardProyecto{}, errors.New("service: no implementado")
}
