package web

import (
	"fmt"
	"html/template"
	"math"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

const (
	graficoMargenIzq = 48
	graficoAnchoPlot = 576
	graficoVelTope   = 32
	graficoVelBase   = 200
	graficoVelAlto   = 168
	graficoHorasTope = 40
	graficoHorasBase = 216
	graficoHorasAlto = 176
)

// formatoEscala formatea una etiqueta del eje Y con hasta dos decimales,
// sin ceros sobrantes.
func formatoEscala(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimSuffix(s, ".")
	if s == "" || s == "-" {
		return "0"
	}
	return s
}

// ticksEjeY divide el eje Y en cuatro partes (cinco etiquetas, de 0 a max),
// con la posición vertical ya calculada.
func ticksEjeY(max float64, base, alto int) []dashboardTick {
	ticks := make([]dashboardTick, 0, 5)
	for i := 0; i <= 4; i++ {
		valor := max * float64(i) / 4
		y := base - int(math.Round(valor*float64(alto)/max))
		ticks = append(ticks, dashboardTick{Y: y, Etiqueta: formatoEscala(valor)})
	}
	return ticks
}

// alturaBarra convierte un valor en píxeles de alto según la escala (D15).
func alturaBarra(valor, max float64, alto int) int {
	if max <= 0 {
		return 0
	}
	return int(math.Round(valor * float64(alto) / max))
}

// puntosVelocidad calcula la geometría de las barras del gráfico (a).
func puntosVelocidad(puntos []service.PuntoVelocidadSprint, escala int) []dashboardPuntoVelocidad {
	if len(puntos) == 0 {
		return nil
	}
	anchoBanda := float64(graficoAnchoPlot) / float64(len(puntos))
	anchoBarra := math.Min(48, anchoBanda*0.6)
	maximo := float64(escala)
	vista := make([]dashboardPuntoVelocidad, 0, len(puntos))
	for i, punto := range puntos {
		alto := alturaBarra(float64(punto.Velocidad), maximo, graficoVelAlto)
		x := graficoMargenIzq + int(float64(i)*anchoBanda+(anchoBanda-anchoBarra)/2)
		y := graficoVelBase - alto
		vista = append(vista, dashboardPuntoVelocidad{
			Sprint:    punto.SprintNumero,
			Velocidad: punto.Velocidad,
			X:         x,
			Y:         y,
			Ancho:     int(anchoBarra),
			Alto:      alto,
			TextoX:    x + int(anchoBarra)/2,
			TextoY:    y - 4,
			Etiqueta:  fmt.Sprintf("%d", punto.Velocidad),
		})
	}
	return vista
}

// puntosHoras calcula la geometría de las barras agrupadas del gráfico (b).
// Los Sprints sin métricas disponibles no dibujan barras (D12).
func puntosHoras(puntos []service.PuntoHorasSprint, escala float64) []dashboardPuntoHoras {
	if len(puntos) == 0 {
		return nil
	}
	anchoBanda := float64(graficoAnchoPlot) / float64(len(puntos))
	anchoBarra := math.Min(28, anchoBanda*0.3)
	const separacion = 6.0
	total := 2*anchoBarra + separacion
	vista := make([]dashboardPuntoHoras, 0, len(puntos))
	for i, punto := range puntos {
		inicio := graficoMargenIzq + float64(i)*anchoBanda + (anchoBanda-total)/2
		xEst := int(inicio)
		xReal := int(inicio + anchoBarra + separacion)
		grupoX := int(graficoMargenIzq + float64(i)*anchoBanda + anchoBanda/2)
		v := dashboardPuntoHoras{
			Sprint:              punto.SprintNumero,
			HorasEstimadas:      punto.HorasEstimadas,
			HorasReales:         punto.HorasReales,
			MetricasDisponibles: punto.MetricasDisponibles,
			Ancho:               int(anchoBarra),
			EstX:                xEst,
			RealX:               xReal,
			TextoEstX:           xEst + int(anchoBarra)/2,
			TextoRealX:          xReal + int(anchoBarra)/2,
			GrupoX:              grupoX,
			EstTexto:            fmt.Sprintf("%.2f", punto.HorasEstimadas),
			RealTexto:           fmt.Sprintf("%.2f", punto.HorasReales),
			DesvTexto:           desviacionTextoPunto(punto),
		}
		if punto.MetricasDisponibles {
			v.EstAlto = alturaBarra(punto.HorasEstimadas, escala, graficoHorasAlto)
			v.RealAlto = alturaBarra(punto.HorasReales, escala, graficoHorasAlto)
			v.EstY = graficoHorasBase - v.EstAlto
			v.RealY = graficoHorasBase - v.RealAlto
			v.TextoEstY = v.EstY - 4
			v.TextoRealY = v.RealY - 4
		}
		vista = append(vista, v)
	}
	return vista
}

// desviacionTextoPunto formatea la desviación de un Sprint para la tabla;
// nil se muestra como "—" (D12, R14).
func desviacionTextoPunto(punto service.PuntoHorasSprint) template.HTML {
	if !punto.MetricasDisponibles || !punto.EstimacionesCompletas {
		return "—"
	}
	desviacion := domain.CalcularDesviacion(punto.HorasEstimadas, punto.HorasReales)
	if desviacion == nil {
		return "—"
	}
	return template.HTML(fmt.Sprintf("%+.2f%%", *desviacion))
}
