package reportepdf

import (
	"errors"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

type Opcion func(*configuracion)

type configuracion struct{}

func SinCompresion() Opcion {
	return func(*configuracion) {}
}

func GenerarReporteProyecto(reporte service.ReporteProyecto, nombreProyecto string, fechaGeneracion time.Time, opciones ...Opcion) ([]byte, error) {
	return nil, errors.New("no implementado")
}

func GenerarReporteSprint(reporte service.ReporteSprint, nombreProyecto string, fechaGeneracion time.Time, opciones ...Opcion) ([]byte, error) {
	return nil, errors.New("no implementado")
}
