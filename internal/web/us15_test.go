package web_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

type servicioDefectosFalsoUS15 struct {
	defectos          []domain.Defecto
	listarErr         error
	registrarErr      error
	cambiarErr        error
	listarLlamadas    int
	registrarLlamadas int
	cambiarLlamadas   int
	ultimoProyecto    string
	ultimoDefecto     domain.Defecto
	ultimoID          int64
	ultimoEstado      string
	ultimoSprint      *int
}

func (f *servicioDefectosFalsoUS15) ListarDefectos(proyectoID string) ([]domain.Defecto, error) {
	f.listarLlamadas++
	f.ultimoProyecto = proyectoID
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return append([]domain.Defecto(nil), f.defectos...), nil
}

func (f *servicioDefectosFalsoUS15) RegistrarDefecto(proyectoID string, defecto domain.Defecto) (domain.Defecto, error) {
	f.registrarLlamadas++
	f.ultimoProyecto = proyectoID
	f.ultimoDefecto = defecto
	if f.registrarErr != nil {
		return domain.Defecto{}, f.registrarErr
	}
	defecto.ID = 1
	f.defectos = append(f.defectos, defecto)
	return defecto, nil
}

func (f *servicioDefectosFalsoUS15) CambiarEstadoDefecto(proyectoID string, defectoID int64, estado string, sprintResolucion *int) (domain.Defecto, error) {
	f.cambiarLlamadas++
	f.ultimoProyecto = proyectoID
	f.ultimoID = defectoID
	f.ultimoEstado = estado
	f.ultimoSprint = sprintResolucion
	if f.cambiarErr != nil {
		return domain.Defecto{}, f.cambiarErr
	}
	defecto := f.defectos[0]
	defecto.Estado = estado
	defecto.SprintResolucion = sprintResolucion
	return defecto, nil
}

func (f *servicioDefectosFalsoUS15) MetricasDefectosPorSprint(string, int) (domain.MetricasDefectos, error) {
	return domain.MetricasDefectos{}, nil
}

func (f *servicioDefectosFalsoUS15) MetricasDefectosPorProyecto(string) (domain.MetricasDefectos, error) {
	return domain.MetricasDefectos{}, nil
}

func TestUS15_WebListaDefectosDelProyecto(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(registroProyectoUS22("PRJ-015", "Proyecto de prueba", fechaUS22(2026, 1, 1), fechaUS22(2026, 12, 31)))
	backlog.proyectosExistentes["PRJ-015"] = true
	sprints := nuevoServicioSprintsFalso()
	sprints.marcarProyectoExistente("PRJ-015")
	defectos := &servicioDefectosFalsoUS15{
		defectos: []domain.Defecto{{ID: 1, Descripcion: "No guarda", Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto}},
	}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Defectos:  defectos,
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-015/defectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No guarda", "Alta", "Abierto")
	if defectos.listarLlamadas != 1 || defectos.ultimoProyecto != "PRJ-015" {
		t.Fatalf("listado: llamadas=%d proyecto=%q", defectos.listarLlamadas, defectos.ultimoProyecto)
	}
}

func TestUS15_WebRegistraDefectoYActualizaResolucion(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(registroProyectoUS22("PRJ-015", "Proyecto de prueba", fechaUS22(2026, 1, 1), fechaUS22(2026, 12, 31)))
	backlog.proyectosExistentes["PRJ-015"] = true
	sprints := nuevoServicioSprintsFalso()
	sprints.marcarProyectoExistente("PRJ-015")
	defectos := &servicioDefectosFalsoUS15{}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Defectos:  defectos,
	})

	formulario := url.Values{
		"descripcion":      {"No guarda"},
		"severidad":        {"Alta"},
		"estado":           {"Abierto"},
		"item_id":          {"HU-01"},
		"sprint_deteccion": {"1"},
	}
	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-015/defectos", formulario)
	exigirEstadoUS22(t, rec, http.StatusSeeOther)
	if defectos.registrarLlamadas != 1 || defectos.ultimoDefecto.ItemID != "HU-01" ||
		defectos.ultimoDefecto.SprintDeteccion != 1 {
		t.Fatalf("registro recibido: %+v, llamadas=%d", defectos.ultimoDefecto, defectos.registrarLlamadas)
	}

	formulario = url.Values{"estado": {"Resuelto"}, "sprint_resolucion": {"2"}}
	rec = solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-015/defectos/1/estado", formulario)
	exigirEstadoUS22(t, rec, http.StatusSeeOther)
	if defectos.cambiarLlamadas != 1 || defectos.ultimoEstado != "Resuelto" ||
		defectos.ultimoSprint == nil || *defectos.ultimoSprint != 2 {
		t.Fatalf("actualización recibida: estado=%q sprint=%v llamadas=%d", defectos.ultimoEstado, defectos.ultimoSprint, defectos.cambiarLlamadas)
	}
}

func TestUS15_WebRechazaIdentificadorDefectoInvalido(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(registroProyectoUS22("PRJ-015", "Proyecto de prueba", fechaUS22(2026, 1, 1), fechaUS22(2026, 12, 31)))
	sprints := nuevoServicioSprintsFalso()
	defectos := &servicioDefectosFalsoUS15{}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Defectos:  defectos,
	})

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-015/defectos/no-numero/estado", url.Values{"estado": {"Resuelto"}})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("código = %d, se esperaba %d; cuerpo: %s", rec.Code, http.StatusUnprocessableEntity, rec.Body.String())
	}
	if strings.TrimSpace(rec.Body.String()) == "" {
		t.Fatal("se esperaba un mensaje de validación visible")
	}
	if defectos.cambiarLlamadas != 0 {
		t.Fatal("el servicio no debía llamarse con un ID inválido")
	}
}
