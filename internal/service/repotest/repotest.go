package repotest

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

// Entorno define la API mínima de la suite de contrato compartida (T-08).
type Entorno struct {
	Proyectos       service.RepositorioProyectos
	Integrantes     service.RepositorioIntegrantes
	Backlog         service.RepositorioBacklog
	SembrarProyecto func(t *testing.T, p service.ProyectoRegistrado)
	InyectarFallo   func(t *testing.T) (causa error, ejecutar func() error)
}

// Ejecutar ejecuta la suite de contrato compartida para un adaptador.
func Ejecutar(t *testing.T, nuevo func(t *testing.T) Entorno) {
	t.Helper()
	casos := []struct {
		nombre string
		fn     func(*testing.T, Entorno)
	}{
		{"TestT08_Proyectos_RepositorioVacioYRoundTrip", ejecutarProyectos},
		{"TestT08_Integrantes_OrdenYAislamiento", ejecutarIntegrantes},
		{"TestT08_Backlog_NumeracionOrdenYAislamiento", ejecutarBacklogNumeracion},
		{"TestT08_Backlog_RoundTripYNoEncontrado", ejecutarBacklogRoundTrip},
		{"TestT08_ErroresInfraestructura_EnvuelvenConW", ejecutarErroresInfra},
		{"TestUS11_EstimarItem_BacklogSinArchivo", ejecutarUS11BacklogSinArchivo},
		{"TestUS11_EstimarItem_MismoIDEnOtroProyecto", ejecutarUS11MismoIDEnOtroProyecto},
		{"TestUS09_MarcarItemCompletada_BacklogEnAmbosAdaptadores", ejecutarUS09MarcarItemCompletada},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			c.fn(t, nuevo(t))
		})
	}
}

func ejecutarProyectos(t *testing.T, e Entorno) {
	t.Helper()
	existe, err := e.Proyectos.Existe("PRJ-001")
	debeNoError(t, err, "Existe vacío")
	if existe {
		t.Fatalf("Existe vacío: esperado false")
	}
	lista, err := e.Proyectos.Listar()
	debeNoError(t, err, "Listar vacío")
	verificarVacio(t, "Listar proyectos", lista)
	_, err = e.Proyectos.Obtener("PRJ-001")
	verificarErrDominio(t, err, domain.ErrProyectoNoEncontrado, "Obtener inexistente")
	p1 := nuevoProyecto(t, "Proyecto Uno", time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC), time.Date(2026, 6, 20, 23, 59, 0, 0, time.UTC))
	err = e.Proyectos.Actualizar("PRJ-001", p1)
	verificarErrDominio(t, err, domain.ErrProyectoNoEncontrado, "Actualizar inexistente")

	e.SembrarProyecto(t, service.ProyectoRegistrado{ID: "PRJ-001", Proyecto: p1})
	existe, err = e.Proyectos.Existe("PRJ-001")
	debeNoError(t, err, "Existe tras siembra")
	if !existe {
		t.Fatalf("Existe tras siembra: esperado true")
	}
	obtenido, err := e.Proyectos.Obtener("PRJ-001")
	debeNoError(t, err, "Obtener tras siembra")
	if obtenido.ID != "PRJ-001" {
		t.Fatalf("ID: esperado PRJ-001, obtenido %q", obtenido.ID)
	}
	verificarProyectoFechas(t, p1, obtenido.Proyecto)

	p2 := nuevoProyecto(t, "Proyecto Uno Actualizado", time.Date(2026, 4, 1, 8, 0, 0, 0, time.FixedZone("X", -3*3600)), time.Date(2026, 7, 1, 18, 0, 0, 0, time.FixedZone("X", -3*3600)))
	debeNoError(t, e.Proyectos.Actualizar("PRJ-001", p2), "Actualizar")
	actualizado, err := e.Proyectos.Obtener("PRJ-001")
	debeNoError(t, err, "Obtener actualizado")
	if actualizado.ID != "PRJ-001" {
		t.Fatalf("ID tras actualizar: esperado PRJ-001, obtenido %q", actualizado.ID)
	}
	verificarProyectoFechas(t, p2, actualizado.Proyecto)
	lista, err = e.Proyectos.Listar()
	debeNoError(t, err, "Listar tras actualizar")
	if len(lista) != 1 {
		t.Fatalf("Listar: esperado 1, obtenido %d", len(lista))
	}
	verificarProyectoFechas(t, p2, lista[0].Proyecto)
}

func ejecutarIntegrantes(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-002", "Proyecto B"))

	lista, err := e.Integrantes.ListarPorProyecto("PRJ-001")
	debeNoError(t, err, "ListarPorProyecto vacío")
	verificarVacio(t, "ListarPorProyecto integrantes", lista)

	i1 := debeNuevoIntegrante(t, "Ana", "Dev")
	i2 := debeNuevoIntegrante(t, "Beto", "QA")
	i3 := debeNuevoIntegrante(t, "Caro", "PM")
	debeNoError(t, e.Integrantes.Agregar("PRJ-001", i1), "Agregar i1")
	debeNoError(t, e.Integrantes.Agregar("PRJ-001", i2), "Agregar i2")
	debeNoError(t, e.Integrantes.Agregar("PRJ-002", i3), "Agregar i3")

	got1, err := e.Integrantes.ListarPorProyecto("PRJ-001")
	debeNoError(t, err, "Listar PRJ-001")
	esperado1 := []domain.Integrante{i1, i2}
	if !slicesEquivalentes(got1, esperado1) {
		t.Fatalf("orden/aislamiento PRJ-001: esperado %#v, obtenido %#v", esperado1, got1)
	}
	got2, err := e.Integrantes.ListarPorProyecto("PRJ-002")
	debeNoError(t, err, "Listar PRJ-002")
	esperado2 := []domain.Integrante{i3}
	if !slicesEquivalentes(got2, esperado2) {
		t.Fatalf("aislamiento PRJ-002: esperado %#v, obtenido %#v", esperado2, got2)
	}
}

func ejecutarBacklogNumeracion(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-002", "Proyecto B"))

	for i := 1; i <= 100; i++ {
		item := debeNuevoItem(t, fmt.Sprintf("Item %d", i), "desc", "Media", nil)
		agregado, err := e.Backlog.AgregarItem("PRJ-001", item)
		debeNoError(t, err, fmt.Sprintf("AgregarItem %d", i))
		esperadoID := fmt.Sprintf("HU-%02d", i)
		if agregado.ID != esperadoID {
			t.Fatalf("ID item %d: esperado %s, obtenido %s", i, esperadoID, agregado.ID)
		}
	}
	items, err := e.Backlog.ListarPorProyecto("PRJ-001")
	debeNoError(t, err, "ListarPorProyecto PRJ-001")
	if len(items) != 100 {
		t.Fatalf("cantidad: esperado 100, obtenido %d", len(items))
	}
	for i, it := range items {
		esperadoID := fmt.Sprintf("HU-%02d", i+1)
		if it.ID != esperadoID {
			t.Fatalf("orden[%d]: esperado %s, obtenido %s", i, esperadoID, it.ID)
		}
	}
	if items[98].ID != "HU-99" || items[99].ID != "HU-100" {
		t.Fatalf("HU-99/HU-100: obtenido %s, %s", items[98].ID, items[99].ID)
	}
	itemB := debeNuevoItem(t, "Item B", "desc", "Alta", []string{})
	agregadoB, err := e.Backlog.AgregarItem("PRJ-002", itemB)
	debeNoError(t, err, "AgregarItem PRJ-002")
	if agregadoB.ID != "HU-01" {
		t.Fatalf("numeración independiente: esperado HU-01 en PRJ-002, obtenido %s", agregadoB.ID)
	}
	itemsB, err := e.Backlog.ListarPorProyecto("PRJ-002")
	debeNoError(t, err, "ListarPorProyecto PRJ-002")
	if len(itemsB) != 1 || itemsB[0].ID != "HU-01" {
		t.Fatalf("aislamiento PRJ-002: %#v", itemsB)
	}
}

func ejecutarBacklogRoundTrip(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-002", "Proyecto B"))

	_, err := e.Backlog.ObtenerItem("PRJ-001", "HU-01")
	verificarErrDominio(t, err, domain.ErrItemNoEncontrado, "ObtenerItem vacío")

	vacio, err := e.Backlog.ListarPorProyecto("PRJ-001")
	debeNoError(t, err, "ListarPorProyecto vacío")
	verificarVacio(t, "ListarPorProyecto backlog", vacio)

	item := debeNuevoItem(t, "Historia", "desc", "Alta", []string{"c1", "c2"})
	agregado, err := e.Backlog.AgregarItem("PRJ-001", item)
	debeNoError(t, err, "AgregarItem")
	if agregado.ID == "" {
		t.Fatalf("AgregarItem debe asignar ID")
	}
	obtenido, err := e.Backlog.ObtenerItem("PRJ-001", agregado.ID)
	debeNoError(t, err, "ObtenerItem")
	if !itemsEquivalentes(agregado, obtenido) {
		t.Fatalf("round-trip agregar/obtener: esperado %#v, obtenido %#v", agregado, obtenido)
	}

	agregado.Titulo = "Historia actualizada"
	agregado.Descripcion = "nueva desc"
	agregado.Prioridad = "Baja"
	agregado.StoryPoints = 5
	agregado.Criterios = []string{"n1"}
	debeNoError(t, e.Backlog.ActualizarItem("PRJ-001", agregado), "ActualizarItem")
	actualizado, err := e.Backlog.ObtenerItem("PRJ-001", agregado.ID)
	debeNoError(t, err, "ObtenerItem actualizado")
	if !itemsEquivalentes(agregado, actualizado) {
		t.Fatalf("round-trip actualizar: esperado %#v, obtenido %#v", agregado, actualizado)
	}

	itemInexistente := debeNuevoItem(t, "X", "", "Media", nil)
	itemInexistente.ID = "HU-99"
	err = e.Backlog.ActualizarItem("PRJ-001", itemInexistente)
	verificarErrDominio(t, err, domain.ErrItemNoEncontrado, "ActualizarItem inexistente")
}

func ejecutarUS09MarcarItemCompletada(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	item := debeNuevoItem(t, "Historia", "desc", "Alta", []string{"criterio"})
	agregado, err := e.Backlog.AgregarItem("PRJ-001", item)
	debeNoError(t, err, "AgregarItem")
	estimado, err := agregado.Estimar(3)
	debeNoError(t, err, "Estimar")
	debeNoError(t, e.Backlog.ActualizarItem("PRJ-001", estimado), "ActualizarItem estimado")

	completado, err := estimado.MarcarComoCompletada()
	debeNoError(t, err, "MarcarComoCompletada")
	debeNoError(t, e.Backlog.ActualizarItem("PRJ-001", completado), "ActualizarItem completado")
	obtenido, err := e.Backlog.ObtenerItem("PRJ-001", agregado.ID)
	debeNoError(t, err, "ObtenerItem completado")
	if obtenido.Estado != domain.EstadoCompletada {
		t.Fatalf("Estado = %q, se esperaba %q", obtenido.Estado, domain.EstadoCompletada)
	}
	if !itemsEquivalentes(completado, obtenido) {
		t.Fatalf("ítem completado: esperado %#v, obtenido %#v", completado, obtenido)
	}
}

func ejecutarErroresInfra(t *testing.T, e Entorno) {
	t.Helper()
	if e.InyectarFallo == nil {
		t.Skip("InyectarFallo no disponible")
	}
	causa, ejecutar := e.InyectarFallo(t)
	err := ejecutar()
	if err == nil {
		t.Fatalf("esperado error de infraestructura")
	}
	if !errors.Is(err, causa) {
		t.Fatalf("errors.Is(err, causa) = false; err=%v", err)
	}
	if err == causa {
		t.Fatalf("error no debe ser la causa pelada")
	}
	var uw interface{ Unwrap() error }
	if !errors.As(err, &uw) || uw.Unwrap() == nil {
		t.Fatalf("error debe estar envuelto con %%w: %T %v", err, err)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("error de infraestructura no debe mapearse a dominio: %v", err)
	}
}

func ejecutarUS11BacklogSinArchivo(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	_, err := e.Backlog.ObtenerItem("PRJ-001", "HU-01")
	verificarErrDominio(t, err, domain.ErrItemNoEncontrado, "backlog sin datos")
}

func ejecutarUS11MismoIDEnOtroProyecto(t *testing.T, e Entorno) {
	t.Helper()
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-001", "Proyecto A"))
	e.SembrarProyecto(t, nuevoRegistroProyecto(t, "PRJ-002", "Proyecto B"))
	item1 := debeNuevoItem(t, "Uno", "", "Alta", nil)
	item2 := debeNuevoItem(t, "Dos", "", "Media", nil)
	_, err := e.Backlog.AgregarItem("PRJ-001", item1)
	debeNoError(t, err, "agregar PRJ-001")
	_, err = e.Backlog.AgregarItem("PRJ-002", item2)
	debeNoError(t, err, "agregar PRJ-002")

	_, err = e.Backlog.ObtenerItem("PRJ-001", "HU-99")
	verificarErrDominio(t, err, domain.ErrItemNoEncontrado, "HU-99 no existe en PRJ-001")
	_, err = e.Backlog.ObtenerItem("PRJ-002", "HU-01")
	debeNoError(t, err, "HU-01 debe existir en PRJ-002")
	_, err = e.Backlog.ObtenerItem("PRJ-001", "HU-02")
	verificarErrDominio(t, err, domain.ErrItemNoEncontrado, "HU-02 no existe en PRJ-001")

	items1, err := e.Backlog.ListarPorProyecto("PRJ-001")
	debeNoError(t, err, "listar PRJ-001")
	items2, err := e.Backlog.ListarPorProyecto("PRJ-002")
	debeNoError(t, err, "listar PRJ-002")
	if len(items1) != 1 || items1[0].ID != "HU-01" || items1[0].Titulo != "Uno" {
		t.Fatalf("aislamiento PRJ-001: %#v", items1)
	}
	if len(items2) != 1 || items2[0].ID != "HU-01" || items2[0].Titulo != "Dos" {
		t.Fatalf("aislamiento PRJ-002: %#v", items2)
	}
}

func nuevoRegistroProyecto(t *testing.T, id, nombre string) service.ProyectoRegistrado {
	t.Helper()
	return service.ProyectoRegistrado{
		ID:       id,
		Proyecto: nuevoProyecto(t, nombre, time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC), time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)),
	}
}

func nuevoProyecto(t *testing.T, nombre string, inicio, fin time.Time) domain.Proyecto {
	t.Helper()
	p, err := domain.NuevoProyecto(nombre, inicio, fin)
	if err != nil {
		t.Fatalf("NuevoProyecto(%q): %v", nombre, err)
	}
	return p
}

func debeNoError(t *testing.T, err error, ctx string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", ctx, err)
	}
}

func debeNuevoIntegrante(t *testing.T, nombre, rol string) domain.Integrante {
	t.Helper()
	i, err := domain.NuevoIntegrante(nombre, rol)
	if err != nil {
		t.Fatalf("NuevoIntegrante(%q): %v", nombre, err)
	}
	return i
}

func debeNuevoItem(t *testing.T, titulo, descripcion, prioridad string, criterios []string) domain.ItemBacklog {
	t.Helper()
	item, err := domain.NuevoItemBacklog(titulo, descripcion, prioridad, criterios)
	if err != nil {
		t.Fatalf("NuevoItemBacklog(%q): %v", titulo, err)
	}
	return item
}

func verificarErrDominio(t *testing.T, err, esperado error, ctx string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: esperado error %v", ctx, esperado)
	}
	if !errors.Is(err, esperado) {
		t.Fatalf("%s: esperado %v, obtenido %v", ctx, esperado, err)
	}
	if err != esperado {
		t.Fatalf("%s: error de dominio no debe estar envuelto: %v", ctx, err)
	}
}

func verificarVacio(t *testing.T, ctx string, v any) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Slice || rv.Len() != 0 {
		t.Fatalf("%s: esperado slice vacío, obtenido %#v", ctx, v)
	}
}

func normalizarSlice(v any) any {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Slice {
		return v
	}
	if rv.IsNil() {
		return reflect.MakeSlice(rv.Type(), 0, 0).Interface()
	}
	return v
}

func slicesEquivalentes(a, b any) bool {
	return reflect.DeepEqual(normalizarSlice(a), normalizarSlice(b))
}

func mismaFecha(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func verificarProyectoFechas(t *testing.T, esperado, obtenido domain.Proyecto) {
	t.Helper()
	if esperado.Nombre != obtenido.Nombre || esperado.Estado != obtenido.Estado {
		t.Fatalf("proyecto: esperado %#v, obtenido %#v", esperado, obtenido)
	}
	if !mismaFecha(esperado.FechaInicio, obtenido.FechaInicio) || !mismaFecha(esperado.FechaFin, obtenido.FechaFin) {
		t.Fatalf("fechas: esperado %v-%v, obtenido %v-%v", esperado.FechaInicio, esperado.FechaFin, obtenido.FechaInicio, obtenido.FechaFin)
	}
}

func itemsEquivalentes(a, b domain.ItemBacklog) bool {
	a.Criterios = normalizarSlice(a.Criterios).([]string)
	b.Criterios = normalizarSlice(b.Criterios).([]string)
	return reflect.DeepEqual(a, b)
}
