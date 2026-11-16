package web_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// altaIntegranteT14 registra los argumentos crudos recibidos por el doble del
// servicio de alta de integrantes.
type altaIntegranteT14 struct {
	proyectoID string
	nombre     string
	rol        string
}

// servicioIntegrantesT14Falso es el doble en memoria de
// web.ServicioProyectosWeb para los tests de T-14. Reproduce el comportamiento
// del servicio de US-03: valida con el dominio y detecta duplicados.
type servicioIntegrantesT14Falso struct {
	registros       []service.ProyectoRegistrado
	integrantes     map[string][]domain.Integrante
	obtenerErr      error
	listarErr       error
	agregarErr      error
	obtenerLlamadas int
	listarLlamadas  int
	agregarLlamadas int
	ultimoObtenerID string
	ultimoListarID  string
	ultimaAlta      altaIntegranteT14
}

var _ web.ServicioProyectosWeb = (*servicioIntegrantesT14Falso)(nil)

func nuevoServicioIntegrantesT14(registros ...service.ProyectoRegistrado) *servicioIntegrantesT14Falso {
	falso := &servicioIntegrantesT14Falso{integrantes: map[string][]domain.Integrante{}}
	falso.registros = append(falso.registros, registros...)
	return falso
}

func (f *servicioIntegrantesT14Falso) existe(proyectoID string) bool {
	for _, registro := range f.registros {
		if registro.ID == proyectoID {
			return true
		}
	}
	return false
}

func (f *servicioIntegrantesT14Falso) CrearProyecto(nombre string, fechaInicio, fechaFin time.Time) (service.ProyectoRegistrado, error) {
	registro := registroProyectoUS22("PRJ-777", nombre, fechaInicio, fechaFin)
	f.registros = append(f.registros, registro)
	return registro, nil
}

func (f *servicioIntegrantesT14Falso) ListarProyectos() ([]service.ProyectoRegistrado, error) {
	return append([]service.ProyectoRegistrado(nil), f.registros...), nil
}

func (f *servicioIntegrantesT14Falso) ObtenerProyecto(id string) (service.ProyectoRegistrado, error) {
	f.obtenerLlamadas++
	f.ultimoObtenerID = id
	if f.obtenerErr != nil {
		return service.ProyectoRegistrado{}, f.obtenerErr
	}
	for _, registro := range f.registros {
		if registro.ID == id {
			return registro, nil
		}
	}
	return service.ProyectoRegistrado{}, domain.ErrProyectoNoEncontrado
}

func (f *servicioIntegrantesT14Falso) ListarIntegrantes(proyectoID string) ([]domain.Integrante, error) {
	f.listarLlamadas++
	f.ultimoListarID = proyectoID
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	if !f.existe(proyectoID) {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return append([]domain.Integrante(nil), f.integrantes[proyectoID]...), nil
}

func (f *servicioIntegrantesT14Falso) AgregarIntegrante(proyectoID, nombre, rol string) (domain.Integrante, error) {
	f.agregarLlamadas++
	f.ultimaAlta = altaIntegranteT14{proyectoID: proyectoID, nombre: nombre, rol: rol}
	if f.agregarErr != nil {
		return domain.Integrante{}, f.agregarErr
	}
	if !f.existe(proyectoID) {
		return domain.Integrante{}, domain.ErrProyectoNoEncontrado
	}
	nuevo, err := domain.NuevoIntegrante(nombre, rol)
	if err != nil {
		return domain.Integrante{}, err
	}
	if domain.EsDuplicado(f.integrantes[proyectoID], nuevo.Nombre) {
		return domain.Integrante{}, domain.ErrIntegranteDuplicado
	}
	f.integrantes[proyectoID] = append(f.integrantes[proyectoID], nuevo)
	return nuevo, nil
}

func nuevoHandlerIntegrantesT14(t *testing.T, proyectos *servicioIntegrantesT14Falso) http.Handler {
	t.Helper()
	if proyectos == nil {
		proyectos = nuevoServicioIntegrantesT14()
	}
	return web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})
}

func capturarLogT14(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buffer bytes.Buffer
	anterior := log.Writer()
	log.SetOutput(&buffer)
	t.Cleanup(func() {
		log.SetOutput(anterior)
	})
	return &buffer
}

func directorioWebT14(t *testing.T) string {
	t.Helper()
	_, archivo, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no se pudo obtener la ruta del test")
	}
	return filepath.Dir(archivo)
}

// CA1 / R2 / R7 / D9: la página muestra el nombre del proyecto obtenido con
// ObtenerProyecto y luego lista sus integrantes.
func TestT14_PaginaIntegrantesMuestraNombreProyecto(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Sistema de Métricas")

	if proyectos.obtenerLlamadas != 1 || proyectos.ultimoObtenerID != "PRJ-001" {
		t.Errorf("ObtenerProyecto: llamadas=%d id=%q, se esperaba 1 llamada con PRJ-001",
			proyectos.obtenerLlamadas, proyectos.ultimoObtenerID)
	}
	if proyectos.listarLlamadas != 1 || proyectos.ultimoListarID != "PRJ-001" {
		t.Errorf("ListarIntegrantes: llamadas=%d id=%q, se esperaba 1 llamada con PRJ-001",
			proyectos.listarLlamadas, proyectos.ultimoListarID)
	}
}

// CA1 / R2 / R7 / R9 / D8 / D13: la página muestra el listado en el orden del
// servicio y el formulario de alta con sus campos y botón.
func TestT14_PaginaIntegrantesMuestraListadoYFormulario(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	proyectos.integrantes["PRJ-001"] = []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
		{Nombre: "Victoria Caro", Rol: "Tester"},
	}
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Sistema de Métricas",
		"Dana Guzmán", "Desarrollador",
		"Victoria Caro", "Tester",
		"Nombre", "Rol",
		`name="nombre"`, `name="rol"`,
		"Registrar integrante",
	)
	exigirCuerpoNoContieneUS22(t, rec, "required")
}

// CA1 / R2 / R9 / D13: un proyecto sin integrantes muestra el estado vacío y el
// formulario sigue disponible.
func TestT14_PaginaIntegrantesSinIntegrantesMuestraVacio(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"El proyecto no tiene integrantes",
		`name="nombre"`, `name="rol"`,
		"Registrar integrante",
	)
}

// CA1 / R10 / D1: el detalle del proyecto ofrece "Ver integrantes" hacia la
// página de integrantes.
func TestT14_DetalleProyectoEnlaceIntegrantes(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Ver integrantes",
		`href="/proyectos/PRJ-001/integrantes"`,
	)
}

// CA1 / CA6 / R10: el detalle conserva los enlaces existentes de US-22, US-15,
// US-19 y US-20 al agregar el enlace de integrantes.
func TestT14_DetalleProyectoConservaEnlacesExistentes(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Ver integrantes", `href="/proyectos/PRJ-001/integrantes"`,
		"Ver backlog", `href="/proyectos/PRJ-001/backlog"`,
		"Ver sprints", `href="/proyectos/PRJ-001/sprints"`,
		"Ver defectos", `href="/proyectos/PRJ-001/defectos"`,
		"Ver dashboard", `href="/proyectos/PRJ-001/dashboard"`,
		"Ver reporte", `href="/proyectos/PRJ-001/reporte"`,
	)
}

// CA2 / R3 / R7: un alta válida responde 303 al listado, delega los valores
// crudos y el integrante aparece al volver a abrir la página.
func TestT14_AltaIntegranteValidaRedirige303(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"Dana Guzmán"},
		"rol":    {"Desarrollador"},
	})
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/integrantes")

	if proyectos.agregarLlamadas != 1 {
		t.Fatalf("AgregarIntegrante se llamó %d veces, se esperaba 1", proyectos.agregarLlamadas)
	}
	esperada := altaIntegranteT14{proyectoID: "PRJ-001", nombre: "Dana Guzmán", rol: "Desarrollador"}
	if proyectos.ultimaAlta != esperada {
		t.Errorf("AgregarIntegrante recibió %+v, se esperaba %+v", proyectos.ultimaAlta, esperada)
	}

	rec = solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Dana Guzmán", "Desarrollador")
	if len(proyectos.integrantes["PRJ-001"]) != 1 {
		t.Errorf("el proyecto debería tener 1 integrante, tiene %d", len(proyectos.integrantes["PRJ-001"]))
	}
}

// CA2 / R3: el primer integrante de un proyecto sin integrantes se registra con
// normalidad.
func TestT14_AltaPrimerIntegrante(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	if len(proyectos.integrantes["PRJ-001"]) != 0 {
		t.Fatalf("el proyecto debía empezar sin integrantes, tiene %d", len(proyectos.integrantes["PRJ-001"]))
	}

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"Dana Guzmán"},
		"rol":    {"Desarrollador"},
	})
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/integrantes")
	if len(proyectos.integrantes["PRJ-001"]) != 1 {
		t.Fatalf("el proyecto debería tener 1 integrante, tiene %d", len(proyectos.integrantes["PRJ-001"]))
	}
	if proyectos.integrantes["PRJ-001"][0].Nombre != "Dana Guzmán" {
		t.Errorf("Nombre = %q, se esperaba %q", proyectos.integrantes["PRJ-001"][0].Nombre, "Dana Guzmán")
	}
}

// CA1 / CA2 / R7: la web no reordena; respeta el orden que devuelve el
// servicio aunque no sea alfabético.
func TestT14_ListadoConservaOrdenDelServicio(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	proyectos.integrantes["PRJ-001"] = []domain.Integrante{
		{Nombre: "Victoria Caro", Rol: "Tester"},
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	cuerpo := rec.Body.String()
	if posicionUS22(t, cuerpo, "Victoria Caro") > posicionUS22(t, cuerpo, "Dana Guzmán") {
		t.Error("Victoria Caro debería aparecer antes que Dana Guzmán: la web no debe reordenar")
	}
}

// CA3 / R4 / E2: nombre vacío o de solo espacios responde 422 con el mensaje
// exacto del dominio y conserva el rol ingresado.
func TestT14_AltaIntegranteNombreVacio422(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"   "},
		"rol":    {"Desarrollador"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrNombreIntegranteVacio.Error(), `value="Desarrollador"`)
	if len(proyectos.integrantes["PRJ-001"]) != 0 {
		t.Errorf("no debería registrarse ningún integrante, hay %d", len(proyectos.integrantes["PRJ-001"]))
	}
}

// CA3 / R4 / E3: rol vacío o de solo espacios responde 422 con el mensaje
// exacto del dominio y conserva el nombre ingresado.
func TestT14_AltaIntegranteRolVacio422(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"Dana Guzmán"},
		"rol":    {"   "},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRolIntegranteVacio.Error(), `value="Dana Guzmán"`)
	if len(proyectos.integrantes["PRJ-001"]) != 0 {
		t.Errorf("no debería registrarse ningún integrante, hay %d", len(proyectos.integrantes["PRJ-001"]))
	}
}

// CA3 / R4 / E4: un nombre duplicado, incluso ignorando mayúsculas y espacios
// extremos, responde 422 con el mensaje exacto y conserva los valores.
func TestT14_AltaIntegranteDuplicado422(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	proyectos.integrantes["PRJ-001"] = []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"  dana guzmán  "},
		"rol":    {"Tester"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrIntegranteDuplicado.Error(), `value="Tester"`)
	if len(proyectos.integrantes["PRJ-001"]) != 1 {
		t.Errorf("el proyecto debería conservar 1 integrante, tiene %d", len(proyectos.integrantes["PRJ-001"]))
	}
}

// CA3 / R7 / E2 / E3 / E4: la web no recorta ni normaliza; envía al servicio y
// re-renderiza exactamente los valores ingresados.
func TestT14_AltaConservaValoresIngresados(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	proyectos.integrantes["PRJ-001"] = []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/integrantes", url.Values{
		"nombre": {"  Dana Guzmán  "},
		"rol":    {"  Desarrollador  "},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec,
		domain.ErrIntegranteDuplicado.Error(),
		`value="  Dana Guzmán  "`,
		`value="  Desarrollador  "`,
	)

	esperada := altaIntegranteT14{proyectoID: "PRJ-001", nombre: "  Dana Guzmán  ", rol: "  Desarrollador  "}
	if proyectos.ultimaAlta != esperada {
		t.Errorf("la web no debe recortar: AgregarIntegrante recibió %+v, se esperaba %+v",
			proyectos.ultimaAlta, esperada)
	}
}

// CA4 / R2 / R5 / E1: un proyecto inexistente en el GET responde 404 con el
// mensaje del dominio.
func TestT14_PaginaIntegrantesProyectoInexistente404(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14()
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if proyectos.listarLlamadas != 0 {
		t.Errorf("no debería listarse integrantes de un proyecto inexistente; llamadas=%d", proyectos.listarLlamadas)
	}
}

// CA4 / R5 / E1: un proyecto inexistente en el POST responde 404 con el mensaje
// del dominio y no registra nada.
func TestT14_AltaIntegranteProyectoInexistente404(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14()
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-999/integrantes", url.Values{
		"nombre": {"Dana Guzmán"},
		"rol":    {"Desarrollador"},
	})
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if len(proyectos.integrantes["PRJ-999"]) != 0 {
		t.Errorf("no debería registrarse ningún integrante, hay %d", len(proyectos.integrantes["PRJ-999"]))
	}
}

// CA4 / R1 / E7: un método distinto de GET/POST sobre la ruta responde 405 sin
// ejecutar el handler.
func TestT14_IntegrantesMetodoNoPermitido405(t *testing.T) {
	for _, metodo := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(metodo, func(t *testing.T) {
			proyectos := nuevoServicioIntegrantesT14(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			handler := nuevoHandlerIntegrantesT14(t, proyectos)

			rec := solicitarUS22(handler, metodo, "/proyectos/PRJ-001/integrantes", nil)
			exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
			if proyectos.obtenerLlamadas+proyectos.listarLlamadas+proyectos.agregarLlamadas != 0 {
				t.Errorf("el handler no debía ejecutarse: %+v", proyectos)
			}
		})
	}
}

// CA4 / R6 / E6: un error de infraestructura responde 500 con el mensaje
// genérico y el detalle queda solo en el log.
func TestT14_IntegrantesErrorInfraestructura500(t *testing.T) {
	detalle := "fallo de conexión a postgres: timeout al listar integrantes"
	casos := []struct {
		nombre     string
		metodo     string
		formulario url.Values
		semilla    func(*servicioIntegrantesT14Falso)
	}{
		{
			nombre: "obtener proyecto",
			metodo: http.MethodGet,
			semilla: func(f *servicioIntegrantesT14Falso) {
				f.obtenerErr = errors.New(detalle)
			},
		},
		{
			nombre: "listar integrantes",
			metodo: http.MethodGet,
			semilla: func(f *servicioIntegrantesT14Falso) {
				f.listarErr = errors.New(detalle)
			},
		},
		{
			nombre:     "agregar integrante",
			metodo:     http.MethodPost,
			formulario: url.Values{"nombre": {"Dana Guzmán"}, "rol": {"Desarrollador"}},
			semilla: func(f *servicioIntegrantesT14Falso) {
				f.agregarErr = errors.New(detalle)
			},
		},
		{
			nombre:     "re-render tras error de validación",
			metodo:     http.MethodPost,
			formulario: url.Values{"nombre": {"   "}, "rol": {"Desarrollador"}},
			semilla: func(f *servicioIntegrantesT14Falso) {
				f.listarErr = errors.New(detalle)
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos := nuevoServicioIntegrantesT14(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			caso.semilla(proyectos)
			handler := nuevoHandlerIntegrantesT14(t, proyectos)
			registro := capturarLogT14(t)

			rec := solicitarUS22(handler, caso.metodo, "/proyectos/PRJ-001/integrantes", caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusInternalServerError)
			exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
			exigirCuerpoNoContieneUS22(t, rec, detalle)
			if !strings.Contains(registro.String(), detalle) {
				t.Errorf("el detalle debía registrarse en el log; log:\n%s", registro.String())
			}
		})
	}
}

// CA4 / R3 / E5: un formulario ilegible responde 422 con el mensaje de la web y
// no llama al servicio de alta.
func TestT14_AltaFormularioIlegible422(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	peticion := httptest.NewRequest(http.MethodPost, "/proyectos/PRJ-001/integrantes", strings.NewReader("nombre=%zz"))
	peticion.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	grabador := httptest.NewRecorder()
	handler.ServeHTTP(grabador, peticion)

	rec := grabador.Result()
	defer rec.Body.Close()
	if rec.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("código = %d, se esperaba %d", rec.StatusCode, http.StatusUnprocessableEntity)
	}
	cuerpo := grabador.Body.String()
	if !strings.Contains(cuerpo, "No se pudo leer el formulario.") {
		t.Errorf("el cuerpo no contiene el mensaje de formulario ilegible; cuerpo:\n%s", cuerpo)
	}
	if proyectos.agregarLlamadas != 0 {
		t.Errorf("AgregarIntegrante se llamó %d veces y no debía ejecutarse", proyectos.agregarLlamadas)
	}
}

// CA6 / R8: un nombre con markup se muestra escapado y no se interpreta como
// HTML.
func TestT14_IntegrantesEscapaNombre(t *testing.T) {
	proyectos := nuevoServicioIntegrantesT14(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	proyectos.integrantes["PRJ-001"] = []domain.Integrante{
		{Nombre: `<script>alert(1)</script>`, Rol: "Tester"},
	}
	handler := nuevoHandlerIntegrantesT14(t, proyectos)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/integrantes", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "&lt;script&gt;alert(1)&lt;/script&gt;")
	exigirCuerpoNoContieneUS22(t, rec, "<script>alert(1)</script>")
}

// CA6 / R8: la capa web no usa template.HTML.
func TestT14_WebSinTemplateHTML(t *testing.T) {
	archivos, err := filepath.Glob(filepath.Join(directorioWebT14(t), "*.go"))
	if err != nil {
		t.Fatalf("no se pudieron listar los archivos de internal/web: %v", err)
	}
	for _, archivo := range archivos {
		if strings.HasSuffix(archivo, "_test.go") {
			continue
		}
		contenido, err := os.ReadFile(archivo)
		if err != nil {
			t.Fatalf("no se pudo leer %s: %v", archivo, err)
		}
		if strings.Contains(string(contenido), "template.HTML") {
			t.Errorf("%s no debe usar template.HTML", filepath.Base(archivo))
		}
	}
}

// CA6 / R7: los handlers de integrantes no duplican reglas de dominio ni
// normalizan valores.
func TestT14_IntegrantesNoDuplicaReglasDeDominio(t *testing.T) {
	contenido, err := os.ReadFile(filepath.Join(directorioWebT14(t), "integrantes.go"))
	if err != nil {
		t.Fatalf("no se pudo leer integrantes.go: %v", err)
	}
	for _, prohibido := range []string{
		"domain.NuevoIntegrante",
		"domain.EsDuplicado",
		"domain.BuscarIntegrante",
		"strings.TrimSpace",
		"sort.",
	} {
		if strings.Contains(string(contenido), prohibido) {
			t.Errorf("integrantes.go no debe contener %q: la regla vive en el dominio", prohibido)
		}
	}
}

// CA5 / R11 / D11: cuando el proyecto no tiene integrantes, Planning Poker
// muestra el aviso con enlace a la página de integrantes en cualquier estado de
// ronda y conserva los botones de inicio/nueva ronda.
func TestT14_PokerSinIntegrantesMuestraAvisoYEnlace(t *testing.T) {
	casos := []struct {
		nombre string
		poker  *servicioPokerFalso
		boton  string
	}{
		{
			nombre: "sin rondas",
			poker:  &servicioPokerFalso{},
			boton:  "Iniciar ronda",
		},
		{
			nombre: "ronda abierta",
			poker: &servicioPokerFalso{
				rondaActualExiste: true,
				rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
				consulta: domain.ConsultaRondaPoker{
					Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
				},
			},
		},
		{
			nombre: "ronda revelada",
			poker: &servicioPokerFalso{
				rondaActualExiste: true,
				rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
				consulta: domain.ConsultaRondaPoker{
					Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
					Votos: []domain.VotoPoker{},
				},
			},
			boton: "Nueva ronda",
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			proyectos.listarIntegrantes = nil
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, caso.poker)

			rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirCuerpoContieneUS22(t, rec,
				"Este proyecto todavía no tiene integrantes",
				"Registrar integrantes",
				`href="/proyectos/PRJ-001/integrantes"`,
			)
			if caso.boton != "" {
				exigirCuerpoContieneUS22(t, rec, caso.boton)
			}
			if proyectos.listarIntegrantesLlamadas != 1 {
				t.Errorf("ListarIntegrantes se llamó %d veces, se esperaba 1",
					proyectos.listarIntegrantesLlamadas)
			}
		})
	}
}

// CA5 / R11 / D11: sin integrantes no se renderizan los formularios que
// dependen de un desplegable de integrantes (votar, revelar y acordar).
func TestT14_PokerSinIntegrantesNoMuestraSelectsVacios(t *testing.T) {
	casos := []struct {
		nombre string
		poker  *servicioPokerFalso
	}{
		{
			nombre: "ronda abierta",
			poker: &servicioPokerFalso{
				rondaActualExiste: true,
				rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
				consulta: domain.ConsultaRondaPoker{
					Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
				},
			},
		},
		{
			nombre: "ronda revelada",
			poker: &servicioPokerFalso{
				rondaActualExiste: true,
				rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
				consulta: domain.ConsultaRondaPoker{
					Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
					Votos: []domain.VotoPoker{},
				},
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			proyectos.listarIntegrantes = nil
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, caso.poker)

			rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirCuerpoNoContieneUS22(t, rec,
				`name="integrante"`,
				"/poker/rondas/1/votos",
				"/poker/rondas/1/revelar",
				"/poker/rondas/1/acordar",
			)
		})
	}
}

// CA5 / R11: con integrantes se conservan los formularios y el comportamiento
// existente de Planning Poker (US-12/US-13).
func TestT14_PokerConIntegrantesConservaFormularios(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}
	poker := &servicioPokerFalso{
		rondaActualExiste: true,
		rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
		consulta: domain.ConsultaRondaPoker{
			Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
		},
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		`<select name="integrante">`,
		`value="Dana Guzmán"`,
		"Dana Guzmán",
		"Votar", "Revelar",
		"/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos",
		"/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar",
	)
	exigirCuerpoNoContieneUS22(t, rec, "Este proyecto todavía no tiene integrantes")
	if proyectos.listarIntegrantesLlamadas != 1 {
		t.Errorf("ListarIntegrantes se llamó %d veces, se esperaba 1",
			proyectos.listarIntegrantesLlamadas)
	}
}
