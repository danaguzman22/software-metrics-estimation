package reportepdf

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

func requierePDFValido(t *testing.T, pdf []byte) {
	t.Helper()
	if len(pdf) == 0 {
		t.Fatal("el PDF está vacío")
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatal("el PDF no comienza con la cabecera esperada")
	}
	if !bytes.HasSuffix(bytes.TrimRight(pdf, "\r\n"), []byte("%%"+"EOF")) {
		t.Fatal("el PDF no termina con el marcador EOF esperado")
	}
	if !bytes.Contains(pdf, []byte("xref")) {
		t.Fatal("el PDF no contiene xref")
	}
	if !bytes.Contains(pdf, []byte("startxref")) {
		t.Fatal("el PDF no contiene startxref")
	}
	if !bytes.Contains(pdf, []byte("/Type /Page")) && !bytes.Contains(pdf, []byte("/Type/Page")) {
		t.Fatal("el PDF no contiene al menos una página")
	}
}

func requiereContenido(t *testing.T, pdf []byte, textos ...string) {
	t.Helper()
	for _, texto := range textos {
		if !bytes.Contains(pdf, []byte(texto)) {
			t.Fatalf("el PDF no contiene %q", texto)
		}
	}
}

func ptrFloat(v float64) *float64 { return &v }

func dtoProyectoConDatos() service.ReporteProyecto {
	return service.ReporteProyecto{
		SPPlanificadosTotal:            88,
		SPCompletadosTotal:             34,
		VelocidadPromedio:              77.25,
		PorcentajeHistoriasCompletadas: 93.75,
		HorasEstimadasTotal:            111.11,
		HorasRealesTotal:               222.22,
		EstimacionesCompletas:          false,
		DesviacionHorasTotal:           ptrFloat(25.0),
		SprintsCerrados:                4,
		SprintsCerradosVelocidad:       []service.PuntoVelocidadSprint{{SprintNumero: 1, Velocidad: 8}},
		HistoriasPlanificadas: []service.HistoriaReporte{{
			ID: "HU-02", Titulo: "Registrar horas", Estado: "Pendiente", StoryPoints: 3,
			HorasEstimadas: ptrFloat(10.0), HorasReales: 9.5,
		}},
		HistoriasCompletadas: []service.HistoriaReporte{{
			ID: "HU-01", Titulo: "Iniciar sesión", Estado: "Completada", StoryPoints: 5,
			HorasEstimadas: ptrFloat(40.0), HorasReales: 35.5,
		}},
		Defectos: service.ResumenDefectos{
			Alta: 9, Media: 8, Baja: 7, Abierto: 5, EnCurso: 4, Resuelto: 3, Total: 24, Detectados: 61, Resueltos: 57,
			Defectos: []domain.Defecto{{ID: 11, Descripcion: "Falla grave", Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto, ItemID: "HU-01", SprintDeteccion: 1}},
		},
	}
}

func dtoSprintConDatos() service.ReporteSprint {
	return service.ReporteSprint{
		Numero: 1, Objetivo: "Objetivo alfa", FechaInicio: time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC),
		Estado: domain.EstadoSprintCerrado,
		HistoriasAsignadas: []service.HistoriaReporte{{ID: "HU-01", Titulo: "Iniciar sesión", Estado: "Completada", StoryPoints: 5, HorasEstimadas: ptrFloat(40.0), HorasReales: 35.5}},
		SPPlanificados: 88, SPCompletados: 34, Velocidad: 13,
		HorasEstimadas: 111.11, HorasReales: 222.22, EstimacionesCompletas: false, MetricasDisponibles: true, DesviacionHoras: ptrFloat(-12.5),
		Defectos: service.ResumenDefectos{Total: 2, Detectados: 2, Resueltos: 1, Defectos: []domain.Defecto{{ID: 21, Descripcion: "Defecto sprint", Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto, ItemID: "HU-01", SprintDeteccion: 1}}},
	}
}

// R11,R1,CA1: PDF de proyecto válido (estructura comprobable).
func TestUS21_Generador_ReporteProyecto_PDFValido(t *testing.T) {
	pdf, err := GenerarReporteProyecto(dtoProyectoConDatos(), "Sistema de Métricas", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
}

// R3,R4,CA1: contenido tomado del DTO tal cual (sin recalcular cifras), visible sin compresión.
func TestUS21_Generador_ReporteProyecto_ContenidoSinCompresion(t *testing.T) {
	pdf, err := GenerarReporteProyecto(dtoProyectoConDatos(), "Sistema de Métricas", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	// Valores del DTO que no son derivables de las listas (R4): deben aparecer tal cual.
	requiereContenido(t, pdf,
		"Sistema de Métricas",
		"HU-01", "HU-02", "Iniciar sesión", "Registrar horas",
		"88", "34", "77.25", "93.75", "111.11", "222.22", "25.00",
		"Falla grave",
	)
}

// R5,R11,CA2: proyecto sin datos genera PDF válido con mensajes vacíos.
func TestUS21_Generador_ReporteProyecto_SinDatos(t *testing.T) {
	pdf, err := GenerarReporteProyecto(service.ReporteProyecto{}, "Proyecto Vacío", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	requiereContenido(t, pdf, "No hay historias", "No hay defectos")
}

// R2,CA1: PDF de Sprint válido.
func TestUS21_Generador_ReporteSprint_PDFValido(t *testing.T) {
	pdf, err := GenerarReporteSprint(dtoSprintConDatos(), "Sistema de Métricas", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
}

// R3,R4,CA1: contenido de Sprint tomado del DTO tal cual.
func TestUS21_Generador_ReporteSprint_ContenidoSinCompresion(t *testing.T) {
	pdf, err := GenerarReporteSprint(dtoSprintConDatos(), "Sistema de Métricas", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	requiereContenido(t, pdf, "Sistema de Métricas", "Objetivo alfa", "HU-01", "88", "34", "13", "111.11", "222.22", "12.50", "Defecto sprint")
}

// R5,CA2: Sprint sin datos genera PDF válido con mensajes vacíos.
func TestUS21_Generador_ReporteSprint_SinDatos(t *testing.T) {
	pdf, err := GenerarReporteSprint(service.ReporteSprint{Numero: 1, Estado: domain.EstadoSprintActivo}, "Proyecto Vacío", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	requiereContenido(t, pdf, "No hay historias asignadas", "No hay defectos")
}

// R5,CA2: Sprint cerrado sin snapshot muestra "métricas no disponibles".
func TestUS21_Generador_ReporteSprint_SinMetricas(t *testing.T) {
	pdf, err := GenerarReporteSprint(service.ReporteSprint{Numero: 1, Estado: domain.EstadoSprintCerrado, MetricasDisponibles: false,
		HistoriasAsignadas: []service.HistoriaReporte{{ID: "HU-01", Titulo: "Historia", Estado: "Completada", StoryPoints: 5}},
	}, "Proyecto", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	requiereContenido(t, pdf, "métricas no disponibles")
}

// R12,CA1: acentos, eñes y normalización de caracteres fuera de cp1252.
func TestUS21_Generador_NormalizacionCaracteres(t *testing.T) {
	dto := service.ReporteProyecto{
		HistoriasPlanificadas: []service.HistoriaReporte{{
			ID: "HU-Ñ", Titulo: "Configuración — módulo “principal” … áéíóú ÁÉÍÓÚ ñ Ñ ¿¡", Estado: "Pendiente",
		}},
	}
	pdf, err := GenerarReporteProyecto(dto, "Proyecto ñandú", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), SinCompresion())
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	requiereContenido(t, pdf, "Configuración", "ñandú")
	// Guion largo/comillas tipográficas/elipsis normalizados a ASCII.
	texto := string(pdf)
	if bytes.Contains(pdf, []byte("—")) || bytes.Contains(pdf, []byte("“")) || bytes.Contains(pdf, []byte("”")) || bytes.Contains(pdf, []byte("…")) {
		t.Fatalf("el PDF contiene caracteres tipográficos sin normalizar")
	}
	if !strings.Contains(texto, "-") || !strings.Contains(texto, "\"") || !strings.Contains(texto, "...") {
		t.Fatalf("el PDF no contiene la normalización esperada: %q", texto)
	}
}

// R11 (caso límite): títulos largos y muchas historias/defectos generan PDF válido con más de una página.
func TestUS21_Generador_MuchasHistoriasYDefectos_Paginacion(t *testing.T) {
	dto := service.ReporteProyecto{}
	for i := 0; i < 80; i++ {
		dto.HistoriasPlanificadas = append(dto.HistoriasPlanificadas, service.HistoriaReporte{ID: "HU-L", Titulo: strings.Repeat("Título muy largo ", 12), Estado: "Pendiente", StoryPoints: 3, HorasEstimadas: ptrFloat(1.0)})
	}
	for i := 0; i < 60; i++ {
		dto.Defectos.Defectos = append(dto.Defectos.Defectos, domain.Defecto{ID: int64(i + 1), Descripcion: strings.Repeat("Descripción extensa ", 10), Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto, ItemID: "HU-L", SprintDeteccion: 1})
	}
	dto.Defectos.Total = len(dto.Defectos.Defectos)
	dto.Defectos.Detectados = dto.Defectos.Total
	pdf, err := GenerarReporteProyecto(dto, "Proyecto Grande", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("error al generar PDF: %v", err)
	}
	requierePDFValido(t, pdf)
	paginas := bytes.Count(pdf, []byte("/Type /Page")) + bytes.Count(pdf, []byte("/Type/Page"))
	if paginas < 2 {
		t.Fatalf("se esperaban múltiples páginas, se encontraron %d", paginas)
	}
}
