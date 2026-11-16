package reportepdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/go-pdf/fpdf"
)

type Opcion func(*configuracion)

type configuracion struct {
	sinCompresion bool
}

// SinCompresion desactiva la compresión de streams para permitir verificar el
// texto del PDF en tests (D6, D10). No altera la validez estructural.
func SinCompresion() Opcion {
	return func(c *configuracion) {
		c.sinCompresion = true
	}
}

const (
	anchoPagina    = 210.0
	altoPagina     = 297.0
	margen         = 15.0
	altoLinea      = 5.0
	limiteInferior = altoPagina - margen
	sinValor       = "—"
)

// normalizar reemplaza los caracteres tipográficos fuera de cp1252 por
// equivalentes seguros y descarta cualquier rune no representable (D8, R12).
// Los caracteres del español (áéíóúñÁÉÍÓÚÑ¿¡) se conservan.
func normalizar(texto string) string {
	var b strings.Builder
	for _, r := range texto {
		switch r {
		case '—', '–', '‒', '―':
			b.WriteByte('-')
		case '“', '”', '„', '‟':
			b.WriteByte('"')
		case '‘', '’', '‚', '‛':
			b.WriteByte('\'')
		case '…':
			b.WriteString("...")
		case '\u00A0', '\u202F', '\u2007':
			b.WriteByte(' ')
		default:
			if r < 0x80 || (r >= 0xA0 && r <= 0xFF) {
				b.WriteRune(r)
			} else {
				b.WriteByte('?')
			}
		}
	}
	return b.String()
}

func formatoNumero(v float64) string { return fmt.Sprintf("%.2f", v) }

func formatoHoras(v float64) string { return fmt.Sprintf("%.2f h", v) }

func formatoHorasPtr(v *float64) string {
	if v == nil {
		return sinValor
	}
	return formatoHoras(*v)
}

func formatoDesviacion(v *float64) string {
	if v == nil {
		return sinValor
	}
	signo := ""
	if *v >= 0 {
		signo = "+"
	}
	return fmt.Sprintf("%s%.2f%%", signo, *v)
}

type generador struct {
	pdf *fpdf.Fpdf
}

func nuevoGenerador(titulo string, fechaGeneracion time.Time, cfg configuracion) *generador {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(margen, margen, margen)
	pdf.SetAutoPageBreak(true, margen)
	pdf.SetCompression(!cfg.sinCompresion)
	pdf.SetFont("Helvetica", "", 10)

	tituloNormalizado := normalizar(titulo)
	fechaNormalizada := normalizar(fechaGeneracion.Format("2006-01-02"))

	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.CellFormat(0, 8, tituloNormalizado, "", 1, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(0, 5, "Generado: "+fechaNormalizada, "", 1, "L", false, 0, "")
		pdf.Ln(2)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Helvetica", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Página %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	return &generador{pdf: pdf}
}

func (g *generador) salida() ([]byte, error) {
	var buf bytes.Buffer
	if err := g.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("no se pudo generar el PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func (g *generador) titulo(texto string) {
	g.pdf.SetFont("Helvetica", "B", 11)
	g.pdf.CellFormat(0, 7, normalizar(texto), "", 1, "L", false, 0, "")
	g.pdf.Ln(1)
}

func (g *generador) parrafo(texto string) {
	g.pdf.SetFont("Helvetica", "", 10)
	g.pdf.MultiCell(0, 5, normalizar(texto), "", "L", false)
}

func (g *generador) dato(clave, valor string) {
	const anchoClave = 65.0
	pdf := g.pdf
	x := pdf.GetX()
	y := pdf.GetY()

	pdf.SetFont("Helvetica", "B", 10)
	pdf.MultiCell(anchoClave, 5, normalizar(clave), "", "L", false)
	alto := pdf.GetY() - y

	pdf.SetXY(x+anchoClave, y)
	pdf.SetFont("Helvetica", "", 10)
	pdf.MultiCell(0, 5, normalizar(valor), "", "L", false)
	if altoValor := pdf.GetY() - y; altoValor > alto {
		alto = altoValor
	}
	pdf.SetXY(x, y+alto)
}

type columna struct {
	titulo  string
	ancho   float64
	derecha bool
}

func (g *generador) tabla(columnas []columna, filas [][]string) {
	if len(filas) == 0 {
		return
	}
	pdf := g.pdf
	if pdf.GetY()+g.altoFila(columnas, filas[0])+6 > limiteInferior {
		pdf.AddPage()
	}
	g.dibujarCabecera(columnas)
	for _, fila := range filas {
		if pdf.GetY()+g.altoFila(columnas, fila) > limiteInferior {
			pdf.AddPage()
			g.dibujarCabecera(columnas)
		}
		g.dibujarFila(columnas, fila)
	}
	pdf.Ln(2)
}

func (g *generador) dibujarCabecera(columnas []columna) {
	pdf := g.pdf
	pdf.SetFont("Helvetica", "B", 9)
	x := pdf.GetX()
	y := pdf.GetY()
	desplazamiento := 0.0
	for _, c := range columnas {
		pdf.SetXY(x+desplazamiento, y)
		pdf.CellFormat(c.ancho, 6, normalizar(c.titulo), "1", 0, "L", false, 0, "")
		desplazamiento += c.ancho
	}
	pdf.SetXY(x, y+6)
}

func (g *generador) dibujarFila(columnas []columna, fila []string) {
	pdf := g.pdf
	pdf.SetFont("Helvetica", "", 9)
	x := pdf.GetX()
	y := pdf.GetY()
	desplazamiento := 0.0
	altoMax := 0.0
	for i, c := range columnas {
		valor := ""
		if i < len(fila) {
			valor = fila[i]
		}
		alineacion := "L"
		if c.derecha {
			alineacion = "R"
		}
		pdf.SetXY(x+desplazamiento, y)
		pdf.MultiCell(c.ancho, altoLinea, normalizar(valor), "1", alineacion, false)
		if alto := pdf.GetY() - y; alto > altoMax {
			altoMax = alto
		}
		desplazamiento += c.ancho
	}
	pdf.SetXY(x, y+altoMax)
}

func (g *generador) altoFila(columnas []columna, fila []string) float64 {
	pdf := g.pdf
	pdf.SetFont("Helvetica", "", 9)
	lineasMax := 1
	for i, c := range columnas {
		valor := ""
		if i < len(fila) {
			valor = fila[i]
		}
		if lineas := len(pdf.SplitLines([]byte(normalizar(valor)), c.ancho)); lineas > lineasMax {
			lineasMax = lineas
		}
	}
	return float64(lineasMax) * altoLinea
}

var columnasHistorias = []columna{
	{titulo: "ID", ancho: 18},
	{titulo: "Título", ancho: 55},
	{titulo: "Estado", ancho: 25},
	{titulo: "Story Points", ancho: 22, derecha: true},
	{titulo: "Horas estimadas", ancho: 28, derecha: true},
	{titulo: "Horas reales (acumuladas)", ancho: 32, derecha: true},
}

var columnasDefectos = []columna{
	{titulo: "Descripción", ancho: 60},
	{titulo: "Severidad", ancho: 22},
	{titulo: "Estado", ancho: 25},
	{titulo: "Historia", ancho: 25},
	{titulo: "Sprint de detección", ancho: 23},
}

var columnasDefectosResueltos = []columna{
	{titulo: "Descripción", ancho: 65},
	{titulo: "Severidad", ancho: 25},
	{titulo: "Historia", ancho: 30},
	{titulo: "Sprint de resolución", ancho: 35},
}

func filasHistorias(historias []service.HistoriaReporte) [][]string {
	filas := make([][]string, 0, len(historias))
	for _, h := range historias {
		filas = append(filas, []string{
			h.ID,
			h.Titulo,
			h.Estado,
			fmt.Sprintf("%d", h.StoryPoints),
			formatoHorasPtr(h.HorasEstimadas),
			formatoHoras(h.HorasReales),
		})
	}
	return filas
}

func (g *generador) seccionHistorias(titulo, vacio string, historias []service.HistoriaReporte) {
	g.titulo(titulo)
	if len(historias) == 0 {
		g.parrafo(vacio)
		return
	}
	g.tabla(columnasHistorias, filasHistorias(historias))
}

func (g *generador) seccionDefectos(titulo string, resumen service.ResumenDefectos) {
	g.titulo(titulo)
	g.dato("Detectados", fmt.Sprintf("%d", resumen.Detectados))
	g.dato("Resueltos", fmt.Sprintf("%d", resumen.Resueltos))
	g.dato("Total", fmt.Sprintf("%d", resumen.Total))
	g.dato("Severidad alta", fmt.Sprintf("%d", resumen.Alta))
	g.dato("Severidad media", fmt.Sprintf("%d", resumen.Media))
	g.dato("Severidad baja", fmt.Sprintf("%d", resumen.Baja))
	g.dato("Abiertos", fmt.Sprintf("%d", resumen.Abierto))
	g.dato("En curso", fmt.Sprintf("%d", resumen.EnCurso))
	g.dato("Resueltos por estado", fmt.Sprintf("%d", resumen.Resuelto))
	g.pdf.Ln(2)

	if len(resumen.Defectos) == 0 {
		g.parrafo("No hay defectos")
		return
	}
	filas := make([][]string, 0, len(resumen.Defectos))
	for _, d := range resumen.Defectos {
		filas = append(filas, []string{
			d.Descripcion,
			d.Severidad,
			d.Estado,
			d.ItemID,
			fmt.Sprintf("%d", d.SprintDeteccion),
		})
	}
	g.tabla(columnasDefectos, filas)
}

// GenerarReporteProyecto genera el PDF del reporte del proyecto.
// nombreProyecto se usa en el encabezado; fechaGeneracion la provee la web
// (Dependencias.Ahora). No recalcula cifras: solo maqueta y formatea (D2).
func GenerarReporteProyecto(reporte service.ReporteProyecto, nombreProyecto string, fechaGeneracion time.Time, opciones ...Opcion) ([]byte, error) {
	cfg := configuracion{}
	for _, opcion := range opciones {
		opcion(&cfg)
	}

	g := nuevoGenerador("Reporte de "+nombreProyecto, fechaGeneracion, cfg)
	g.pdf.AddPage()

	g.dato("Story Points planificados", fmt.Sprintf("%d", reporte.SPPlanificadosTotal))
	g.dato("Story Points completados", fmt.Sprintf("%d", reporte.SPCompletadosTotal))
	g.dato("Velocidad promedio", formatoNumero(reporte.VelocidadPromedio))
	g.dato("Historias completadas", formatoNumero(reporte.PorcentajeHistoriasCompletadas)+"%")
	g.dato("Horas estimadas", formatoHoras(reporte.HorasEstimadasTotal))
	g.dato("Horas reales", formatoHoras(reporte.HorasRealesTotal))
	g.dato("Desviación de horas", formatoDesviacion(reporte.DesviacionHorasTotal))
	g.dato("Sprints cerrados", fmt.Sprintf("%d", reporte.SprintsCerrados))
	g.pdf.Ln(3)

	if !reporte.EstimacionesCompletas {
		g.parrafo("Hay historias sin estimar en horas; las horas estimadas suman solo las disponibles.")
	}
	g.parrafo("Horas reales (acumuladas): suma histórica de esfuerzo de la historia en el proyecto; no está acotada al Sprint. Una celda " + sinValor + " en Horas estimadas significa que la historia no está estimada en horas (sin estimar).")
	g.pdf.Ln(2)

	g.seccionHistorias("Historias planificadas", "No hay historias", reporte.HistoriasPlanificadas)
	g.seccionHistorias("Historias completadas", "No hay historias", reporte.HistoriasCompletadas)

	g.titulo("Sprints cerrados")
	if len(reporte.SprintsCerradosVelocidad) == 0 {
		g.parrafo("Todavía no hay Sprints cerrados")
	} else {
		filas := make([][]string, 0, len(reporte.SprintsCerradosVelocidad))
		for _, p := range reporte.SprintsCerradosVelocidad {
			filas = append(filas, []string{fmt.Sprintf("%d", p.SprintNumero), fmt.Sprintf("%d", p.Velocidad)})
		}
		g.tabla([]columna{{titulo: "Sprint", ancho: 30}, {titulo: "Velocidad", ancho: 30, derecha: true}}, filas)
	}

	g.seccionDefectos("Defectos", reporte.Defectos)

	return g.salida()
}

// GenerarReporteSprint genera el PDF del reporte de un Sprint.
// No recalcula cifras: solo maqueta y formatea (D2).
func GenerarReporteSprint(reporte service.ReporteSprint, nombreProyecto string, fechaGeneracion time.Time, opciones ...Opcion) ([]byte, error) {
	cfg := configuracion{}
	for _, opcion := range opciones {
		opcion(&cfg)
	}

	titulo := fmt.Sprintf("Reporte del Sprint %d de %s", reporte.Numero, nombreProyecto)
	g := nuevoGenerador(titulo, fechaGeneracion, cfg)
	g.pdf.AddPage()

	g.dato("Sprint Goal", reporte.Objetivo)
	g.dato("Fecha de inicio", reporte.FechaInicio.Format("2006-01-02"))
	g.dato("Fecha de fin", reporte.FechaFin.Format("2006-01-02"))
	g.dato("Estado", reporte.Estado)
	g.dato("Story Points planificados", fmt.Sprintf("%d", reporte.SPPlanificados))
	g.dato("Story Points completados", fmt.Sprintf("%d", reporte.SPCompletados))
	g.dato("Velocidad", fmt.Sprintf("%d", reporte.Velocidad))
	g.pdf.Ln(3)

	if !reporte.EstimacionesCompletas {
		g.parrafo("Hay historias sin estimar en horas; las horas estimadas suman solo las disponibles.")
	}
	g.parrafo("Horas reales (acumuladas): suma histórica de esfuerzo de la historia en el proyecto; no está acotada al Sprint.")
	g.pdf.Ln(2)

	g.seccionHistorias("Historias asignadas", "No hay historias asignadas", reporte.HistoriasAsignadas)

	g.titulo("Horas estimadas vs reales (horas)")
	if reporte.MetricasDisponibles {
		filas := [][]string{{
			formatoHoras(reporte.HorasEstimadas),
			formatoHoras(reporte.HorasReales),
			formatoDesviacion(reporte.DesviacionHoras),
		}}
		g.tabla([]columna{
			{titulo: "Horas estimadas (horas)", ancho: 60, derecha: true},
			{titulo: "Horas reales (horas)", ancho: 60, derecha: true},
			{titulo: "Desviación", ancho: 45, derecha: true},
		}, filas)
	} else {
		g.parrafo("métricas no disponibles")
	}

	g.seccionDefectos("Defectos detectados en el Sprint", reporte.Defectos)

	g.titulo("Defectos resueltos en el Sprint")
	if len(reporte.Defectos.DefectosResueltos) == 0 {
		g.parrafo("No hay defectos resueltos")
	} else {
		filas := make([][]string, 0, len(reporte.Defectos.DefectosResueltos))
		for _, d := range reporte.Defectos.DefectosResueltos {
			sprintResolucion := sinValor
			if d.SprintResolucion != nil {
				sprintResolucion = fmt.Sprintf("%d", *d.SprintResolucion)
			}
			filas = append(filas, []string{d.Descripcion, d.Severidad, d.ItemID, sprintResolucion})
		}
		g.tabla(columnasDefectosResueltos, filas)
	}

	return g.salida()
}
