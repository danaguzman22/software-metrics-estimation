package web

import (
	"errors"
	"net/http"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

// dashboardProyecto muestra el dashboard del proyecto. El estado se calcula
// con domain.EstadoPorFechas y el reloj inyectable (D14); las métricas se
// muestran tal como las devuelve el servicio (D5, R19).
func (a *aplicacion) dashboardProyecto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	estado, err := domain.EstadoPorFechas(registro.Proyecto.FechaInicio, registro.Proyecto.FechaFin, a.deps.Ahora())
	if err != nil {
		a.responderErrorInfraestructura(w, "calcular estado del proyecto "+id, err)
		return
	}

	datos, err := a.deps.Dashboard.ObtenerDashboard(id)
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
			return
		}
		a.responderErrorInfraestructura(w, "obtener dashboard del proyecto "+id, err)
		return
	}

	a.renderizar(w, http.StatusOK, "dashboard", datosPagina{
		Titulo:    "Dashboard",
		Nav:       navConProyecto(id),
		Dashboard: dashboardVistaDesde(datos, registro, estado),
		Migas: migasDeProyecto(registro.Proyecto.Nombre, id,
			Miga{Texto: "Dashboard", URL: "/proyectos/" + id + "/dashboard"},
		),
		Volver: "/proyectos/" + id,
	})
}

// dashboardVistaDesde copia los valores del servicio sin recalcularlos; solo
// formatea la desviación con signo explícito (D16, R18, R23).
func dashboardVistaDesde(datos service.DashboardProyecto, registro service.ProyectoRegistrado, estado domain.Estado) *dashboardVista {
	vista := &dashboardVista{
		ProyectoID:                    registro.ID,
		ProyectoNombre:                registro.Proyecto.Nombre,
		Estado:                        estado,
		SPPlanificadosTotal:           datos.SPPlanificadosTotal,
		SPCompletadosTotal:            datos.SPCompletadosTotal,
		VelocidadPromedio:             datos.VelocidadPromedio,
		PorcentajeHistoriasCompletadas: datos.PorcentajeHistoriasCompletadas,
		HorasEstimadasTotal:           datos.HorasEstimadasTotal,
		HorasRealesTotal:              datos.HorasRealesTotal,
		EstimacionesCompletas:         datos.EstimacionesCompletas,
		SprintsCerrados:               datos.SprintsCerrados,
	}
	if datos.DesviacionHorasTotal != nil {
		vista.DesviacionTotalDisponible = true
		vista.DesviacionHorasTotal = *datos.DesviacionHorasTotal
	}

	valoresVelocidad := make([]float64, 0, len(datos.VelocidadPorSprint))
	for _, punto := range datos.VelocidadPorSprint {
		valoresVelocidad = append(valoresVelocidad, float64(punto.Velocidad))
	}
	valoresHoras := make([]float64, 0, 2*len(datos.HorasPorSprint))
	for _, punto := range datos.HorasPorSprint {
		valoresHoras = append(valoresHoras, punto.HorasEstimadas, punto.HorasReales)
	}
	vista.EscalaVelocidadMax = int(escalaConMinimoUno(valoresVelocidad...))
	vista.EscalaHorasMax = escalaConMinimoUno(valoresHoras...)
	vista.TicksVelocidad = ticksEjeY(float64(vista.EscalaVelocidadMax), graficoVelBase, graficoVelAlto)
	vista.TicksHoras = ticksEjeY(vista.EscalaHorasMax, graficoHorasBase, graficoHorasAlto)
	vista.Velocidad = puntosVelocidad(datos.VelocidadPorSprint, vista.EscalaVelocidadMax)
	vista.Horas = puntosHoras(datos.HorasPorSprint, vista.EscalaHorasMax)
	for _, punto := range vista.Horas {
		if !punto.MetricasDisponibles {
			vista.HayMetricasNoDisponibles = true
			break
		}
	}
	return vista
}

// escalaConMinimoUno devuelve el tope del eje Y de un gráfico: el máximo de las
// barras o 1 cuando todas valen cero (D15).
func escalaConMinimoUno(valores ...float64) float64 {
	maximo := 0.0
	for _, valor := range valores {
		if valor > maximo {
			maximo = valor
		}
	}
	if maximo < 1 {
		return 1
	}
	return maximo
}
