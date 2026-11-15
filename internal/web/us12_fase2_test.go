package web_test

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// Métodos FASE 2 (FASE 2 de US-12): dobles de web.ServicioPokerWeb (nuevo) y de
// web.ServicioProyectosWeb con la lectura de integrantes (D31). El struct
// servicioProyectosFalso de US-22 no se toca: este doble es propio de US-12 y
// ya implementa ListarIntegrantes para cuando la interfaz la incorpore.

type llamadaIniciarRondaPoker struct {
	proyectoID string
	itemID     string
}

type llamadaVotarPoker struct {
	proyectoID       string
	itemID           string
	rondaNumero      int
	integranteNombre string
	valor            int
}

type llamadaRevelarPoker struct {
	proyectoID       string
	itemID           string
	rondaNumero      int
	integranteNombre string
}

// servicioPokerFalso es el doble en memoria de web.ServicioPokerWeb.
type servicioPokerFalso struct {
	iniciarRondaResultado domain.RondaPoker
	iniciarRondaErr       error
	iniciarRondaLlamadas  int
	ultimaIniciar         llamadaIniciarRondaPoker

	votarErr      error
	votarLlamadas int
	ultimoVoto    llamadaVotarPoker

	revelarErr       error
	revelarLlamadas  int
	ultimaRevelacion llamadaRevelarPoker

	consulta             domain.ConsultaRondaPoker
	consultaErr          error
	consultaLlamadas     int
	ultimaConsultaNumero int

	rondaActual         domain.RondaPoker
	rondaActualExiste   bool
	rondaActualErr      error
	rondaActualLlamadas int
}

var _ web.ServicioPokerWeb = (*servicioPokerFalso)(nil)

func (f *servicioPokerFalso) IniciarRonda(proyectoID, itemID string) (domain.RondaPoker, error) {
	f.iniciarRondaLlamadas++
	f.ultimaIniciar = llamadaIniciarRondaPoker{proyectoID: proyectoID, itemID: itemID}
	if f.iniciarRondaErr != nil {
		return domain.RondaPoker{}, f.iniciarRondaErr
	}
	return f.iniciarRondaResultado, nil
}

func (f *servicioPokerFalso) Votar(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) error {
	f.votarLlamadas++
	f.ultimoVoto = llamadaVotarPoker{
		proyectoID:       proyectoID,
		itemID:           itemID,
		rondaNumero:      rondaNumero,
		integranteNombre: integranteNombre,
		valor:            valor,
	}
	return f.votarErr
}

func (f *servicioPokerFalso) RevelarRonda(proyectoID, itemID string, rondaNumero int, integranteNombre string) error {
	f.revelarLlamadas++
	f.ultimaRevelacion = llamadaRevelarPoker{
		proyectoID:       proyectoID,
		itemID:           itemID,
		rondaNumero:      rondaNumero,
		integranteNombre: integranteNombre,
	}
	return f.revelarErr
}

func (f *servicioPokerFalso) ConsultarRonda(proyectoID, itemID string, rondaNumero int) (domain.ConsultaRondaPoker, error) {
	f.consultaLlamadas++
	f.ultimaConsultaNumero = rondaNumero
	if f.consultaErr != nil {
		return domain.ConsultaRondaPoker{}, f.consultaErr
	}
	return f.consulta, nil
}

func (f *servicioPokerFalso) RondaActual(proyectoID, itemID string) (domain.RondaPoker, bool, error) {
	f.rondaActualLlamadas++
	if f.rondaActualErr != nil {
		return domain.RondaPoker{}, false, f.rondaActualErr
	}
	return f.rondaActual, f.rondaActualExiste, nil
}

// servicioProyectosPokerFalso es el doble en memoria de
// web.ServicioProyectosWeb ampliado con ListarIntegrantes (D31).
type servicioProyectosPokerFalso struct {
	registros                 []service.ProyectoRegistrado
	obtenerErr                error
	listarIntegrantes         []domain.Integrante
	listarIntegrantesErr      error
	listarIntegrantesLlamadas int
}

var _ web.ServicioProyectosWeb = (*servicioProyectosPokerFalso)(nil)

func (f *servicioProyectosPokerFalso) CrearProyecto(nombre string, fechaInicio, fechaFin time.Time) (service.ProyectoRegistrado, error) {
	registro := registroProyectoUS22("PRJ-777", nombre, fechaInicio, fechaFin)
	f.registros = append(f.registros, registro)
	return registro, nil
}

func (f *servicioProyectosPokerFalso) ListarProyectos() ([]service.ProyectoRegistrado, error) {
	return append([]service.ProyectoRegistrado(nil), f.registros...), nil
}

func (f *servicioProyectosPokerFalso) ObtenerProyecto(id string) (service.ProyectoRegistrado, error) {
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

func (f *servicioProyectosPokerFalso) ListarIntegrantes(proyectoID string) ([]domain.Integrante, error) {
	f.listarIntegrantesLlamadas++
	if f.listarIntegrantesErr != nil {
		return nil, f.listarIntegrantesErr
	}
	return append([]domain.Integrante(nil), f.listarIntegrantes...), nil
}

func integraPUS12(nombre, rol string) domain.Integrante {
	return domain.Integrante{Nombre: nombre, Rol: rol}
}

func nuevoEntornoPokerUS12(registro service.ProyectoRegistrado, items ...domain.ItemBacklog) (*servicioProyectosPokerFalso, *servicioBacklogFalso) {
	proyectos := &servicioProyectosPokerFalso{registros: []service.ProyectoRegistrado{registro}}
	backlog := nuevoServicioBacklogFalso()
	backlog.proyectosExistentes[registro.ID] = true
	for _, item := range items {
		backlog.agregarItem(registro.ID, item)
	}
	return proyectos, backlog
}

func nuevoHandlerPokerUS12(t *testing.T, proyectos *servicioProyectosPokerFalso, backlog *servicioBacklogFalso, poker *servicioPokerFalso) http.Handler {
	t.Helper()
	if proyectos == nil {
		proyectos = &servicioProyectosPokerFalso{}
	}
	if backlog == nil {
		backlog = nuevoServicioBacklogFalso()
		for _, registro := range proyectos.registros {
			backlog.proyectosExistentes[registro.ID] = true
		}
	}
	if poker == nil {
		poker = &servicioPokerFalso{}
	}

	sprints := nuevoServicioSprintsFalso()
	for _, registro := range proyectos.registros {
		sprints.marcarProyectoExistente(registro.ID)
	}

	return web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Poker:     poker,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})
}

func registroPokerUS12() service.ProyectoRegistrado {
	return registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30))
}

func integrantesPokerUS12() []domain.Integrante {
	return []domain.Integrante{
		integraPUS12("Ana", "Desarrolladora"),
		integraPUS12("Beto", "QA"),
		integraPUS12("Caro", "Product Owner"),
	}
}

func rutaPokerUS12(itemID string) string {
	return "/proyectos/PRJ-001/backlog/" + itemID + "/poker"
}

func formularioVotoUS12(integrante, valor string) url.Values {
	return url.Values{"integrante": {integrante}, "valor": {valor}}
}

// CA1 / R23: sin rondas la página ofrece iniciar una ronda y ningún formulario
// de voto ni tabla.
func TestUS12_Fase2_Poker_PaginaSinRondasOfreceIniciar(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{rondaActualExiste: false}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Planning Poker", "HU-01", "Iniciar ronda", "/proyectos/PRJ-001/backlog/HU-01/poker/rondas")
	exigirCuerpoNoContieneUS22(t, rec, `name="valor"`, "Revelar", "Nueva ronda", "<table")

	if poker.rondaActualLlamadas != 1 {
		t.Errorf("RondaActual se llamó %d veces, se esperaba 1", poker.rondaActualLlamadas)
	}
}

// CA1 / R26: iniciar una ronda desde la página responde 303 a la misma página y
// no renderiza contenido.
func TestUS12_Fase2_Poker_IniciarRondaRedirige303(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{iniciarRondaResultado: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", nil)
	exigirRedireccionUS22(t, rec, rutaPokerUS12("HU-01"))

	if poker.iniciarRondaLlamadas != 1 {
		t.Fatalf("IniciarRonda se llamó %d veces, se esperaba 1", poker.iniciarRondaLlamadas)
	}
	if poker.ultimaIniciar != (llamadaIniciarRondaPoker{proyectoID: "PRJ-001", itemID: "HU-01"}) {
		t.Errorf("IniciarRonda recibió %+v, se esperaba PRJ-001/HU-01", poker.ultimaIniciar)
	}
}

// CA1 / R26: con la ronda anterior Revelada la página ofrece abrir la siguiente
// y el POST redirige 303 a la misma página.
func TestUS12_Fase2_Poker_IniciarRondaTrasRevelarRedirige303(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{
		rondaActualExiste:     true,
		rondaActual:           domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
		consulta:              domain.ConsultaRondaPoker{Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada}, Votos: []domain.VotoPoker{{Integrante: "Ana", Valor: 5}}},
		iniciarRondaResultado: domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaAbierta},
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Revelada", "Nueva ronda")
	exigirCuerpoNoContieneUS22(t, rec, `name="valor"`)

	rec = solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", nil)
	exigirRedireccionUS22(t, rec, rutaPokerUS12("HU-01"))
	if poker.iniciarRondaLlamadas != 1 {
		t.Errorf("IniciarRonda se llamó %d veces, se esperaba 1", poker.iniciarRondaLlamadas)
	}
}

// CA2 / R24: con la ronda Abierta la página lista quién ya votó (solo nombres),
// nunca valores, y ofrece los formularios de voto y de revelado.
func TestUS12_Fase2_Poker_PaginaAbiertaListaVotantesSinValores(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := &servicioPokerFalso{
		rondaActualExiste: true,
		rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
		consulta:          domain.ConsultaRondaPoker{Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}, YaVotaron: []string{"Ana", "Beto"}},
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Planning Poker", "Abierta", "Ya votaron", "Ana", "Beto",
		"Los valores permanecen ocultos hasta revelar",
		`<select name="integrante">`, `<select name="valor">`,
		`value="1"`, `value="2"`, `value="3"`, `value="5"`, `value="8"`, `value="13"`, `value="21"`,
		"Votar", "Revelar",
		"/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos",
		"/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar",
	)
	exigirCuerpoNoContieneUS22(t, rec, "<table", `<th>Voto</th>`)

	if poker.consultaLlamadas != 1 {
		t.Errorf("ConsultarRonda se llamó %d veces, se esperaba 1", poker.consultaLlamadas)
	}
	if poker.ultimaConsultaNumero != 1 {
		t.Errorf("ConsultarRonda recibió el número %d, se esperaba 1", poker.ultimaConsultaNumero)
	}
	if proyectos.listarIntegrantesLlamadas != 1 {
		t.Errorf("ListarIntegrantes se llamó %d veces, se esperaba 1", proyectos.listarIntegrantesLlamadas)
	}
}

// CA2 / R33: con la ronda Abierta el HTML nunca transporta valores de votos,
// aunque el doble los devolviera. Se aserta sobre celdas y encabezados de la
// tabla (los valores 5 y 8 sí aparecen como <option> del select, que es
// presentación permitida).
func TestUS12_Fase2_Poker_HtmlAbiertaNoExponeValores(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := &servicioPokerFalso{
		rondaActualExiste: true,
		rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
		consulta: domain.ConsultaRondaPoker{
			Ronda:     domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
			YaVotaron: []string{"Ana", "Beto"},
			Votos:     []domain.VotoPoker{{Integrante: "Ana", Valor: 5}, {Integrante: "Beto", Valor: 8}},
		},
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Ana", "Beto")
	exigirCuerpoNoContieneUS22(t, rec, "<td>5</td>", "<td>8</td>", "<th>Voto</th>", "<table")
}

// CA1 / R27: votar desde la página responde 303 a la misma página y delega el
// voto en el servicio con los argumentos del formulario.
func TestUS12_Fase2_Poker_VotarRedirige303(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := &servicioPokerFalso{}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("Ana", "5"))
	exigirRedireccionUS22(t, rec, rutaPokerUS12("HU-01"))

	if poker.votarLlamadas != 1 {
		t.Fatalf("Votar se llamó %d veces, se esperaba 1", poker.votarLlamadas)
	}
	esperada := llamadaVotarPoker{proyectoID: "PRJ-001", itemID: "HU-01", rondaNumero: 1, integranteNombre: "Ana", valor: 5}
	if poker.ultimoVoto != esperada {
		t.Errorf("Votar recibió %+v, se esperaba %+v", poker.ultimoVoto, esperada)
	}
}

// CA4 / R29 / E17: votar con un integrante ajeno al proyecto muestra el
// mensaje exacto del dominio con HTTP 422.
func TestUS12_Fase2_Poker_VotarIntegranteNoDelProyecto422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{votarErr: domain.ErrIntegranteNoDelProyecto}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("Caro", "8"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrIntegranteNoDelProyecto.Error())
}

// CA2 / R25: con la ronda Revelada la página muestra todos los votos y un
// botón "Nueva ronda", y ningún formulario de voto ni de revelado.
func TestUS12_Fase2_Poker_RevelarMuestraVotos(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := &servicioPokerFalso{
		rondaActualExiste: true,
		rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
		consulta: domain.ConsultaRondaPoker{
			Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada},
			Votos: []domain.VotoPoker{{Integrante: "Ana", Valor: 5}, {Integrante: "Beto", Valor: 8}},
		},
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Revelada", "Ana", "<td>5</td>", "Beto", "<td>8</td>", "Nueva ronda")
	exigirCuerpoNoContieneUS22(t, rec, `name="valor"`, `name="integrante"`)
}

// CA2 / R29 / E20: revelar una ronda ya Revelada muestra el mensaje exacto del
// dominio con HTTP 422.
func TestUS12_Fase2_Poker_RevelarYaRevelada422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{revelarErr: domain.ErrRondaRevelada}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar", url.Values{"integrante": {"Ana"}})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaRevelada.Error())
}

// CA1 / R27 / R29 / E20: votar sobre una ronda Revelada muestra el mensaje
// exacto del dominio con HTTP 422.
func TestUS12_Fase2_Poker_VotarRondaRevelada422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{votarErr: domain.ErrRondaRevelada}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("Ana", "5"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaRevelada.Error())
}

// CA1 / R29 / E20: iniciar con una ronda Abierta pendiente muestra el mensaje
// exacto del dominio con HTTP 422.
func TestUS12_Fase2_Poker_IniciarConRondaAbierta422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{iniciarRondaErr: domain.ErrRondaAbiertaExistente}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaAbiertaExistente.Error())
}

// CA1 / R29 / E18: el integrante vacío o de solo espacios muestra el mensaje
// exacto del dominio con HTTP 422.
func TestUS12_Fase2_Poker_NombreIntegranteVacio422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{votarErr: domain.ErrNombreIntegranteVacio}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("   ", "5"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrNombreIntegranteVacio.Error())
}

// CA2 / R32 / E19: un valor fuera de la serie viaja al dominio y se muestra el
// mensaje exacto del dominio con HTTP 422 (la web no valida la serie).
func TestUS12_Fase2_Poker_ValorFueraDeSerie422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{votarErr: domain.ErrValorVotoPokerInvalido}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("Ana", "7"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrValorVotoPokerInvalido.Error())

	if poker.votarLlamadas != 1 {
		t.Errorf("Votar se llamó %d veces, se esperaba 1 (la web no prevalida la serie)", poker.votarLlamadas)
	}
}

// CA1 / R29 / E21: {numero} = 0 parsea pero el dominio responde
// ErrNumeroRondaInvalido; se muestra con HTTP 422.
func TestUS12_Fase2_Poker_NumeroRondaCero422(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "votar", ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/0/votos"},
		{nombre: "revelar", ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/0/revelar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{votarErr: domain.ErrNumeroRondaInvalido, revelarErr: domain.ErrNumeroRondaInvalido}
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			formulario := formularioVotoUS12("Ana", "5")
			if caso.nombre == "revelar" {
				formulario = url.Values{"integrante": {"Ana"}}
			}
			rec := solicitarUS22(handler, http.MethodPost, caso.ruta, formulario)
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, domain.ErrNumeroRondaInvalido.Error())

			if caso.nombre == "votar" && poker.votarLlamadas != 1 {
				t.Errorf("Votar se llamó %d veces, se esperaba 1 (el servicio valida el número)", poker.votarLlamadas)
			}
			if caso.nombre == "revelar" && poker.revelarLlamadas != 1 {
				t.Errorf("RevelarRonda se llamó %d veces, se esperaba 1 (el servicio valida el número)", poker.revelarLlamadas)
			}
		})
	}
}

// CA1 / R31 / V11: un valor no numérico en el formulario produce el mensaje de
// formato de la web con HTTP 422 sin llamar al servicio.
func TestUS12_Fase2_Poker_ValorFormularioNoNumerico422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formularioVotoUS12("Ana", "ocho"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "El voto debe ser un número entero")

	if poker.votarLlamadas != 0 {
		t.Errorf("Votar se llamó %d veces y no debía ejecutarse", poker.votarLlamadas)
	}
}

// CA1 / R29 / E16: votar sobre una ronda inexistente devuelve HTTP 404 con el
// mensaje del dominio.
func TestUS12_Fase2_Poker_VotarRondaInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{votarErr: domain.ErrRondaNoEncontrada}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/9/votos", formularioVotoUS12("Ana", "5"))
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaNoEncontrada.Error())
}

// CA2 / R29 / E15 / E16: la historia inexistente devuelve HTTP 404 con el
// mensaje del dominio en cualquier ruta de poker.
func TestUS12_Fase2_Poker_HistoriaInexistente404(t *testing.T) {
	casos := []struct {
		nombre     string
		metodo     string
		ruta       string
		formulario url.Values
		semilla    func(*servicioPokerFalso)
	}{
		{nombre: "página", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerFalso) { p.rondaActualErr = domain.ErrItemNoEncontrado }},
		{nombre: "iniciar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", semilla: func(p *servicioPokerFalso) { p.iniciarRondaErr = domain.ErrItemNoEncontrado }},
		{nombre: "votar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formulario: formularioVotoUS12("Ana", "5"), semilla: func(p *servicioPokerFalso) { p.votarErr = domain.ErrItemNoEncontrado }},
		{nombre: "revelar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar", formulario: url.Values{"integrante": {"Ana"}}, semilla: func(p *servicioPokerFalso) { p.revelarErr = domain.ErrItemNoEncontrado }},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			caso.semilla(poker)
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())
		})
	}
}

// CA4 / R29 / E14: el proyecto inexistente devuelve HTTP 404 con el mensaje del
// dominio en cualquier ruta de poker.
func TestUS12_Fase2_Poker_ProyectoInexistente404(t *testing.T) {
	casos := []struct {
		nombre     string
		metodo     string
		ruta       string
		formulario url.Values
		semilla    func(*servicioPokerFalso)
	}{
		{nombre: "página", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerFalso) { p.rondaActualErr = domain.ErrProyectoNoEncontrado }},
		{nombre: "iniciar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", semilla: func(p *servicioPokerFalso) { p.iniciarRondaErr = domain.ErrProyectoNoEncontrado }},
		{nombre: "votar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formulario: formularioVotoUS12("Ana", "5"), semilla: func(p *servicioPokerFalso) { p.votarErr = domain.ErrProyectoNoEncontrado }},
		{nombre: "revelar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar", formulario: url.Values{"integrante": {"Ana"}}, semilla: func(p *servicioPokerFalso) { p.revelarErr = domain.ErrProyectoNoEncontrado }},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			caso.semilla(poker)
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
		})
	}
}

// CA1 / R20-D: cuando ConsultarRonda reporta una ronda inexistente en la
// página actual, la respuesta es HTTP 404 con el mensaje del dominio.
func TestUS12_Fase2_Poker_ConsultarRondaInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerFalso{
		rondaActualExiste: true,
		rondaActual:       domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta},
		consultaErr:       domain.ErrRondaNoEncontrada,
	}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaNoEncontrada.Error())
}

// CA1 / R28 / D28: el backlog muestra, por historia, un enlace "Planning
// Poker" hacia la página de poker de esa historia.
func TestUS12_Fase2_Poker_BacklogMuestraEnlacePorHistoria(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(
		registroPokerUS12(),
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cerrar sesión", "Media", "Pendiente", 0),
	)
	poker := &servicioPokerFalso{}
	handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	cuerpo := rec.Body.String()
	for _, itemID := range []string{"HU-01", "HU-02"} {
		href := "/proyectos/PRJ-001/backlog/" + itemID + "/poker"
		if strings.Count(cuerpo, href) != 1 {
			t.Errorf("el enlace %q aparece %d veces, se esperaba 1", href, strings.Count(cuerpo, href))
		}
	}
	if strings.Count(cuerpo, "Planning Poker") < 2 {
		t.Errorf("se esperaba el enlace \"Planning Poker\" por historia; cuerpo:\n%s", cuerpo)
	}
}

// CA2 / R31 / V11 / E22: un parámetro de ruta {numero} no numérico produce el
// mensaje de formato con HTTP 422 sin llamar al servicio.
func TestUS12_Fase2_Poker_NumeroRutaNoNumerico422(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "votar", ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/abc/votos"},
		{nombre: "revelar", ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/abc/revelar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			formulario := formularioVotoUS12("Ana", "5")
			if caso.nombre == "revelar" {
				formulario = url.Values{"integrante": {"Ana"}}
			}
			rec := solicitarUS22(handler, http.MethodPost, caso.ruta, formulario)
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, "El número de la ronda debe ser un número entero")

			if poker.votarLlamadas != 0 {
				t.Errorf("Votar se llamó %d veces y no debía ejecutarse", poker.votarLlamadas)
			}
			if poker.revelarLlamadas != 0 {
				t.Errorf("RevelarRonda se llamó %d veces y no debía ejecutarse", poker.revelarLlamadas)
			}
		})
	}
}

// CA2 / R31 / V12 / E23: un itemID que no cumple el formato HU-NN produce el
// mensaje de formato con HTTP 422 sin llamar al servicio, en cualquier ruta de
// poker.
func TestUS12_Fase2_Poker_ItemIDFormatoInvalido422(t *testing.T) {
	casos := []struct {
		nombre     string
		metodo     string
		ruta       string
		formulario url.Values
	}{
		{nombre: "página", metodo: http.MethodGet, ruta: rutaPokerUS12("XX-1")},
		{nombre: "iniciar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/XX-1/poker/rondas"},
		{nombre: "votar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/XX-1/poker/rondas/1/votos", formulario: formularioVotoUS12("Ana", "5")},
		{nombre: "revelar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/XX-1/poker/rondas/1/revelar", formulario: url.Values{"integrante": {"Ana"}}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, "El ID de la historia debe tener el formato HU-NN")

			if poker.rondaActualLlamadas != 0 || poker.iniciarRondaLlamadas != 0 || poker.votarLlamadas != 0 || poker.revelarLlamadas != 0 || poker.consultaLlamadas != 0 {
				t.Errorf("el servicio de poker se llamó pese al formato inválido: %+v", poker)
			}
		})
	}
}

// CA2 / R30 / E24: un error de infraestructura se registra en el log y se
// responde con el mensaje genérico y HTTP 500; el detalle nunca llega al
// cuerpo.
func TestUS12_Fase2_Poker_ErrorInfraestructura500(t *testing.T) {
	infra := errors.New("fallo de conexión a postgres: timeout al listar rondas")
	casos := []struct {
		nombre     string
		metodo     string
		ruta       string
		formulario url.Values
		semilla    func(*servicioPokerFalso, *servicioProyectosPokerFalso)
	}{
		{nombre: "página (RondaActual)", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) { p.rondaActualErr = infra }},
		{nombre: "página (ConsultarRonda)", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) {
			p.rondaActualExiste = true
			p.rondaActual = domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}
			p.consultaErr = infra
		}},
		{nombre: "página (ListarIntegrantes)", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) {
			p.rondaActualExiste = true
			p.rondaActual = domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}
			p.consulta = domain.ConsultaRondaPoker{Ronda: domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}, YaVotaron: []string{"Ana"}}
			s.listarIntegrantesErr = infra
		}},
		{nombre: "iniciar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas", semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) { p.iniciarRondaErr = infra }},
		{nombre: "votar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos", formulario: formularioVotoUS12("Ana", "5"), semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) { p.votarErr = infra }},
		{nombre: "revelar", metodo: http.MethodPost, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar", formulario: url.Values{"integrante": {"Ana"}}, semilla: func(p *servicioPokerFalso, s *servicioProyectosPokerFalso) { p.revelarErr = infra }},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			caso.semilla(poker, proyectos)
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusInternalServerError)
			exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
			exigirCuerpoNoContieneUS22(t, rec, "fallo de conexión a postgres", "timeout")
		})
	}
}

// CA2 / R25 / E25: un método no permitido en una ruta registrada devuelve HTTP
// 405 y no ejecuta el handler.
func TestUS12_Fase2_Poker_MetodoNoPermitido405(t *testing.T) {
	casos := []struct {
		nombre string
		metodo string
		ruta   string
	}{
		{nombre: "POST en la página", metodo: http.MethodPost, ruta: rutaPokerUS12("HU-01")},
		{nombre: "GET en iniciar", metodo: http.MethodGet, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas"},
		{nombre: "GET en votar", metodo: http.MethodGet, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/votos"},
		{nombre: "GET en revelar", metodo: http.MethodGet, ruta: "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/revelar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerFalso{}
			handler := nuevoHandlerPokerUS12(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)

			if poker.rondaActualLlamadas != 0 || poker.iniciarRondaLlamadas != 0 || poker.votarLlamadas != 0 || poker.revelarLlamadas != 0 {
				t.Errorf("el servicio de poker se ejecutó con un método no permitido: %+v", poker)
			}
		})
	}
}
