package web_test

// Tests de la FASE 2 de US-13 (interfaz web de Planning Poker: diferencias,
// historial y acuerdo). Los dobles de US-12 (servicioPokerFalso) y US-22
// quedan intactos; aquí vive un doble propio que implementa la superficie
// ampliada de web.ServicioPokerWeb propuesta en la sección 20 de
// specs/US-13.md: DiferenciasRonda, HistorialRondas, AcordarEstimacion e
// HistoriaActual (D30 y D36 opción A).

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

type llamadaAcordarPoker struct {
	proyectoID       string
	itemID           string
	rondaNumero      int
	integranteNombre string
	valor            int
}

// servicioPokerUS13Falso es el doble en memoria de web.ServicioPokerWeb
// ampliado con los métodos de US-13. acordarFunc permite que un test conecte
// el acuerdo con el doble del backlog (como hace el servicio real al delegar
// en EstimarItem); si es nil se usa acordarResultado/acordarErr.
type servicioPokerUS13Falso struct {
	rondaActual         domain.RondaPoker
	rondaActualExiste   bool
	rondaActualErr      error
	rondaActualLlamadas int

	consulta             domain.ConsultaRondaPoker
	consultaErr          error
	consultaLlamadas     int
	ultimaConsultaNumero int

	diferencias             domain.DiferenciaPoker
	diferenciasErr          error
	diferenciasLlamadas     int
	ultimaDiferenciasNumero int

	historial         []domain.RondaHistorialPoker
	historialErr      error
	historialLlamadas int

	acordarResultado domain.ItemBacklog
	acordarErr       error
	acordarLlamadas  int
	ultimoAcuerdo    llamadaAcordarPoker
	acordarFunc      func(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) (domain.ItemBacklog, error)

	historiaActualResultado domain.ItemBacklog
	historiaActualErr       error
	historiaActualLlamadas  int
}

var _ web.ServicioPokerWeb = (*servicioPokerUS13Falso)(nil)

func (f *servicioPokerUS13Falso) IniciarRonda(proyectoID, itemID string) (domain.RondaPoker, error) {
	return domain.RondaPoker{}, nil
}

func (f *servicioPokerUS13Falso) Votar(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) error {
	return nil
}

func (f *servicioPokerUS13Falso) RevelarRonda(proyectoID, itemID string, rondaNumero int, integranteNombre string) error {
	return nil
}

func (f *servicioPokerUS13Falso) ConsultarRonda(proyectoID, itemID string, rondaNumero int) (domain.ConsultaRondaPoker, error) {
	f.consultaLlamadas++
	f.ultimaConsultaNumero = rondaNumero
	if f.consultaErr != nil {
		return domain.ConsultaRondaPoker{}, f.consultaErr
	}
	return f.consulta, nil
}

func (f *servicioPokerUS13Falso) RondaActual(proyectoID, itemID string) (domain.RondaPoker, bool, error) {
	f.rondaActualLlamadas++
	if f.rondaActualErr != nil {
		return domain.RondaPoker{}, false, f.rondaActualErr
	}
	return f.rondaActual, f.rondaActualExiste, nil
}

func (f *servicioPokerUS13Falso) DiferenciasRonda(proyectoID, itemID string, rondaNumero int) (domain.DiferenciaPoker, error) {
	f.diferenciasLlamadas++
	f.ultimaDiferenciasNumero = rondaNumero
	if f.diferenciasErr != nil {
		return domain.DiferenciaPoker{}, f.diferenciasErr
	}
	return f.diferencias, nil
}

func (f *servicioPokerUS13Falso) HistorialRondas(proyectoID, itemID string) ([]domain.RondaHistorialPoker, error) {
	f.historialLlamadas++
	if f.historialErr != nil {
		return nil, f.historialErr
	}
	return append([]domain.RondaHistorialPoker(nil), f.historial...), nil
}

func (f *servicioPokerUS13Falso) AcordarEstimacion(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) (domain.ItemBacklog, error) {
	f.acordarLlamadas++
	f.ultimoAcuerdo = llamadaAcordarPoker{
		proyectoID:       proyectoID,
		itemID:           itemID,
		rondaNumero:      rondaNumero,
		integranteNombre: integranteNombre,
		valor:            valor,
	}
	if f.acordarFunc != nil {
		return f.acordarFunc(proyectoID, itemID, rondaNumero, integranteNombre, valor)
	}
	if f.acordarErr != nil {
		return domain.ItemBacklog{}, f.acordarErr
	}
	return f.acordarResultado, nil
}

func (f *servicioPokerUS13Falso) HistoriaActual(proyectoID, itemID string) (domain.ItemBacklog, error) {
	f.historiaActualLlamadas++
	if f.historiaActualErr != nil {
		return domain.ItemBacklog{}, f.historiaActualErr
	}
	return f.historiaActualResultado, nil
}

func nuevoHandlerPokerUS13(t *testing.T, proyectos *servicioProyectosPokerFalso, backlog *servicioBacklogFalso, poker *servicioPokerUS13Falso) http.Handler {
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
		poker = &servicioPokerUS13Falso{}
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

func rutaAcordarPokerUS13(itemID string, ronda int) string {
	return fmt.Sprintf("/proyectos/PRJ-001/backlog/%s/poker/rondas/%d/acordar", itemID, ronda)
}

func formularioAcuerdoUS13(integrante, valor string) url.Values {
	return url.Values{"integrante": {integrante}, "valor": {valor}}
}

func rondaReveladaUS13(numero int) domain.RondaPoker {
	return domain.RondaPoker{Numero: numero, Estado: domain.EstadoRondaRevelada}
}

func rondaAbiertaUS13(numero int) domain.RondaPoker {
	return domain.RondaPoker{Numero: numero, Estado: domain.EstadoRondaAbierta}
}

func consultaReveladaUS13(numero int, votos ...domain.VotoPoker) domain.ConsultaRondaPoker {
	return domain.ConsultaRondaPoker{Ronda: rondaReveladaUS13(numero), Votos: votos}
}

func pokerReveladaUS13() *servicioPokerUS13Falso {
	diferencia := domain.DiferenciaPoker{
		HayDiferencia: true,
		Minimo:        5,
		Maximo:        8,
		QuienesMinimo: []string{"Ana"},
		QuienesMaximo: []string{"Beto"},
	}
	return &servicioPokerUS13Falso{
		rondaActualExiste: true,
		rondaActual:       rondaReveladaUS13(1),
		consulta: consultaReveladaUS13(1,
			domain.VotoPoker{Integrante: "Ana", Valor: 5},
			domain.VotoPoker{Integrante: "Beto", Valor: 8},
		),
		diferencias: diferencia,
		historial: []domain.RondaHistorialPoker{
			{Ronda: rondaReveladaUS13(1), Votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 5},
				{Integrante: "Beto", Valor: 8},
			}},
		},
		historiaActualResultado: domain.ItemBacklog{ID: "HU-01", StoryPoints: 3},
	}
}

// CA1 / R1 / D31 / D32 / D33: con la ronda Revelada la página muestra las
// diferencias (mínimo, máximo y quiénes votaron cada extremo), el historial
// de rondas y el formulario para acordar la estimación (select de integrante
// y de valor).
func TestUS13_Fase2_Poker_DiferenciasVisiblesTrasRevelar(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
		"Diferencias", "Hay diferencia", "mínimo", "5", "Ana", "máximo", "8", "Beto",
	)
	exigirCuerpoContieneUS22(t, rec,
		"Historial", "Ana: 5", "Beto: 8",
	)
	exigirCuerpoContieneUS22(t, rec,
		`action="/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/acordar"`,
		"Acordar estimación", `name="integrante"`, `name="valor"`,
		`value="1"`, `value="2"`, `value="3"`, `value="5"`, `value="8"`, `value="13"`, `value="21"`,
	)

	if poker.diferenciasLlamadas != 1 {
		t.Errorf("DiferenciasRonda se llamó %d veces, se esperaba 1", poker.diferenciasLlamadas)
	}
	if poker.ultimaDiferenciasNumero != 1 {
		t.Errorf("DiferenciasRonda recibió el número %d, se esperaba 1", poker.ultimaDiferenciasNumero)
	}
	if poker.historialLlamadas != 1 {
		t.Errorf("HistorialRondas se llamó %d veces, se esperaba 1", poker.historialLlamadas)
	}
}

// CA1 / R1 / R15 / D29: con la ronda Abierta la página no muestra las
// diferencias ni el formulario de acordar ni los valores de la ronda actual;
// el historial lista las rondas reveladas previas con sus votos y las
// pendientes como "votos ocultos".
func TestUS13_Fase2_Poker_AbiertaSinDiferenciasNiValores(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := &servicioPokerUS13Falso{
		rondaActualExiste: true,
		rondaActual:       rondaAbiertaUS13(2),
		consulta: domain.ConsultaRondaPoker{
			Ronda:     rondaAbiertaUS13(2),
			YaVotaron: []string{"Ana"},
			Votos:     []domain.VotoPoker{{Integrante: "Ana", Valor: 5}},
		},
		historial: []domain.RondaHistorialPoker{
			{Ronda: rondaReveladaUS13(1), Votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 5},
				{Integrante: "Vicky", Valor: 8},
			}},
			{Ronda: rondaAbiertaUS13(2), Votos: nil},
		},
	}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec, "Abierta", "votos ocultos", "Ana: 5", "Vicky: 8")
	exigirCuerpoNoContieneUS22(t, rec,
		"/acordar", "Acordar estimación", "Diferencias", "Hay diferencia", "<td>5</td>",
	)
	if poker.diferenciasLlamadas != 0 {
		t.Errorf("DiferenciasRonda se llamó %d veces y no debía ejecutarse con la ronda Abierta", poker.diferenciasLlamadas)
	}
}

// CA2 / R4 / D32: la sección Historial lista las rondas en orden ascendente;
// las Revelada muestran "Integrante: valor" y las Abierta muestran
// "votos ocultos".
func TestUS13_Fase2_Poker_HistorialVisible(t *testing.T) {
	t.Run("con ronda actual Revelada", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
		proyectos.listarIntegrantes = integrantesPokerUS12()
		poker := pokerReveladaUS13()
		poker.historial = []domain.RondaHistorialPoker{
			{Ronda: rondaReveladaUS13(1), Votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 5},
				{Integrante: "Vicky", Valor: 8},
			}},
			{Ronda: rondaReveladaUS13(2), Votos: []domain.VotoPoker{
				{Integrante: "Caro", Valor: 13},
			}},
		}
		poker.rondaActual = rondaReveladaUS13(2)
		poker.consulta = consultaReveladaUS13(2, domain.VotoPoker{Integrante: "Caro", Valor: 13})
		handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

		rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec,
			"Historial", "Ronda", "Estado", "Votos",
			"Ana: 5", "Vicky: 8", "Caro: 13",
		)
		if poker.historialLlamadas != 1 {
			t.Errorf("HistorialRondas se llamó %d veces, se esperaba 1", poker.historialLlamadas)
		}
	})
	t.Run("ronda actual Abierta muestra votos ocultos en el historial", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
		proyectos.listarIntegrantes = integrantesPokerUS12()
		poker := &servicioPokerUS13Falso{
			rondaActualExiste: true,
			rondaActual:       rondaAbiertaUS13(2),
			consulta: domain.ConsultaRondaPoker{
				Ronda:     rondaAbiertaUS13(2),
				YaVotaron: []string{"Ana"},
			},
			historial: []domain.RondaHistorialPoker{
				{Ronda: rondaReveladaUS13(1), Votos: []domain.VotoPoker{
					{Integrante: "Ana", Valor: 5},
					{Integrante: "Vicky", Valor: 8},
				}},
				{Ronda: rondaAbiertaUS13(2), Votos: nil},
			},
		}
		handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

		rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Historial", "Ana: 5", "Vicky: 8", "votos ocultos")
	})
}

// CA3 / R5 / D28: acordar desde la página responde 303 a la misma página, sin
// contenido, y delega en el servicio con el integrante y el valor del
// formulario.
func TestUS13_Fase2_Poker_AcordarRedirige303(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.acordarResultado = domain.ItemBacklog{ID: "HU-01", StoryPoints: 8}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "8"))
	exigirRedireccionUS22(t, rec, rutaPokerUS12("HU-01"))

	if poker.acordarLlamadas != 1 {
		t.Fatalf("AcordarEstimacion se llamó %d veces, se esperaba 1", poker.acordarLlamadas)
	}
	esperada := llamadaAcordarPoker{proyectoID: "PRJ-001", itemID: "HU-01", rondaNumero: 1, integranteNombre: "Ana", valor: 8}
	if poker.ultimoAcuerdo != esperada {
		t.Errorf("AcordarEstimacion recibió %+v, se esperaba %+v", poker.ultimoAcuerdo, esperada)
	}
}

// CA3 / R8: volver a acordar en una ronda posterior reemplaza el valor
// anterior en la historia (Story Points).
func TestUS13_Fase2_Poker_ReAcordarReemplazaStoryPoints(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12(), itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.acordarFunc = func(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) (domain.ItemBacklog, error) {
		return backlog.EstimarItem(proyectoID, itemID, valor)
	}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "5"))
	exigirEstadoUS22(t, rec, http.StatusSeeOther)

	poker.rondaActual = rondaReveladaUS13(2)
	poker.consulta = consultaReveladaUS13(2,
		domain.VotoPoker{Integrante: "Ana", Valor: 3},
		domain.VotoPoker{Integrante: "Beto", Valor: 8},
	)
	poker.historial = append(poker.historial, domain.RondaHistorialPoker{
		Ronda: rondaReveladaUS13(2),
		Votos: []domain.VotoPoker{
			{Integrante: "Ana", Valor: 3},
			{Integrante: "Beto", Valor: 8},
		},
	})

	rec = solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 2), formularioAcuerdoUS13("Beto", "8"))
	exigirEstadoUS22(t, rec, http.StatusSeeOther)

	if poker.acordarLlamadas != 2 {
		t.Fatalf("AcordarEstimacion se llamó %d veces, se esperaban 2", poker.acordarLlamadas)
	}
	if backlog.estimarLlamadas != 2 {
		t.Fatalf("EstimarItem se llamó %d veces, se esperaban 2", backlog.estimarLlamadas)
	}
	if actual := backlog.items["PRJ-001"][0].StoryPoints; actual != 8 {
		t.Errorf("StoryPoints = %d, se esperaba 8 tras el segundo acuerdo", actual)
	}
}

// CA3 / D36: tras acordar, el backlog muestra los Story Points acordados de la
// historia.
func TestUS13_Fase2_Poker_BacklogMuestraStoryPointsTrasAcordar(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12(), itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.acordarFunc = func(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) (domain.ItemBacklog, error) {
		return backlog.EstimarItem(proyectoID, itemID, valor)
	}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "8"))
	exigirRedireccionUS22(t, rec, rutaPokerUS12("HU-01"))

	if poker.acordarLlamadas != 1 {
		t.Fatalf("AcordarEstimacion se llamó %d veces, se esperaba 1", poker.acordarLlamadas)
	}
	if poker.ultimoAcuerdo.valor != 8 {
		t.Errorf("AcordarEstimacion recibió valor %d, se esperaba 8", poker.ultimoAcuerdo.valor)
	}

	rec = solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "HU-01", "<td>8</td>")
}

// CA4 / R5 / D34: acordar sobre una ronda no Revelada muestra el mensaje
// exacto del dominio con HTTP 422 y la página re-renderizada.
func TestUS13_Fase2_Poker_AcordarEnRondaAbierta422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.acordarErr = domain.ErrRondaNoRevelada
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "8"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrRondaNoRevelada.Error())

	if poker.acordarLlamadas != 1 {
		t.Errorf("AcordarEstimacion se llamó %d veces, se esperaba 1", poker.acordarLlamadas)
	}
}

// CA3 / R6 / D34: un valor fuera de la serie muestra el mensaje exacto del
// dominio con HTTP 422; la web no prevalida la serie.
func TestUS13_Fase2_Poker_AcordarValorFueraSerie422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.acordarErr = domain.ErrStoryPointsInvalidos
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "4"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrStoryPointsInvalidos.Error())

	if poker.acordarLlamadas != 1 {
		t.Errorf("AcordarEstimacion se llamó %d veces, se esperaba 1 (la web no prevalida la serie)", poker.acordarLlamadas)
	}
}

// CA3 / R12 / D34: los demás errores de dominio del acuerdo devuelven HTTP
// 422 con el mensaje exacto del dominio.
func TestUS13_Fase2_Poker_AcordarErroresDeDominio422(t *testing.T) {
	casos := []struct {
		nombre       string
		errorDominio error
	}{
		{nombre: "integrante no del proyecto", errorDominio: domain.ErrIntegranteNoDelProyecto},
		{nombre: "nombre de integrante vacío", errorDominio: domain.ErrNombreIntegranteVacio},
		{nombre: "número de ronda inválido", errorDominio: domain.ErrNumeroRondaInvalido},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			proyectos.listarIntegrantes = integrantesPokerUS12()
			poker := pokerReveladaUS13()
			poker.acordarErr = caso.errorDominio
			handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "8"))
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.errorDominio.Error())

			if poker.acordarLlamadas != 1 {
				t.Errorf("AcordarEstimacion se llamó %d veces, se esperaba 1", poker.acordarLlamadas)
			}
		})
	}
}

// CA4 / R12 / D34: proyecto, historia o ronda inexistentes tras el POST
// devuelven HTTP 404 con el mensaje del dominio.
func TestUS13_Fase2_Poker_AcordarInexistente404(t *testing.T) {
	casos := []struct {
		nombre       string
		errorDominio error
	}{
		{nombre: "proyecto", errorDominio: domain.ErrProyectoNoEncontrado},
		{nombre: "historia", errorDominio: domain.ErrItemNoEncontrado},
		{nombre: "ronda", errorDominio: domain.ErrRondaNoEncontrada},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			proyectos.listarIntegrantes = integrantesPokerUS12()
			poker := pokerReveladaUS13()
			poker.acordarErr = caso.errorDominio
			handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, http.MethodPost, rutaAcordarPokerUS13("HU-01", 1), formularioAcuerdoUS13("Ana", "8"))
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			exigirCuerpoContieneUS22(t, rec, caso.errorDominio.Error())

			if poker.acordarLlamadas != 1 {
				t.Errorf("AcordarEstimacion se llamó %d veces, se esperaba 1", poker.acordarLlamadas)
			}
		})
	}
}

// CA3 / D34: un error de infraestructura se muestra con HTTP 500 genérico; el
// detalle nunca llega al cuerpo.
func TestUS13_Fase2_Poker_ErrorInfraestructura500(t *testing.T) {
	infra := errors.New("detalle interno: conexión caída a postgres")
	casos := []struct {
		nombre     string
		metodo     string
		ruta       string
		formulario url.Values
		semilla    func(*servicioPokerUS13Falso)
	}{
		{nombre: "página (DiferenciasRonda)", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerUS13Falso) {
			p.rondaActualExiste = true
			p.rondaActual = rondaReveladaUS13(1)
			p.consulta = consultaReveladaUS13(1)
			p.diferenciasErr = infra
		}},
		{nombre: "página (HistorialRondas)", metodo: http.MethodGet, ruta: rutaPokerUS12("HU-01"), semilla: func(p *servicioPokerUS13Falso) {
			p.rondaActualExiste = true
			p.rondaActual = rondaReveladaUS13(1)
			p.consulta = consultaReveladaUS13(1)
			p.historialErr = infra
		}},
		{nombre: "acordar", metodo: http.MethodPost, ruta: rutaAcordarPokerUS13("HU-01", 1), formulario: formularioAcuerdoUS13("Ana", "8"), semilla: func(p *servicioPokerUS13Falso) {
			p.acordarErr = infra
		}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			proyectos.listarIntegrantes = integrantesPokerUS12()
			poker := pokerReveladaUS13()
			caso.semilla(poker)
			handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, caso.formulario)
			exigirEstadoUS22(t, rec, http.StatusInternalServerError)
			exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
			exigirCuerpoNoContieneUS22(t, rec, "detalle interno", "postgres")
		})
	}
}

// CA3 / D34: un {numero} no numérico produce el mensaje de formato de la web
// con HTTP 422 sin llamar al servicio.
func TestUS13_Fase2_Poker_NumeroRutaNoNumerico422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerUS13Falso{}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/poker/rondas/abc/acordar", formularioAcuerdoUS13("Ana", "8"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "El número de la ronda debe ser un número entero")

	if poker.acordarLlamadas != 0 {
		t.Errorf("AcordarEstimacion se llamó %d veces y no debía ejecutarse", poker.acordarLlamadas)
	}
}

// CA3 / D34: un itemID que no cumple el formato HU-NN produce el mensaje de
// formato con HTTP 422 sin llamar al servicio.
func TestUS13_Fase2_Poker_ItemIDFormatoInvalido422(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerUS13Falso{}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/XX-1/poker/rondas/1/acordar", formularioAcuerdoUS13("Ana", "8"))
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "El ID de la historia debe tener el formato HU-NN")

	if poker.acordarLlamadas != 0 {
		t.Errorf("AcordarEstimacion se llamó %d veces y no debía ejecutarse", poker.acordarLlamadas)
	}
}

// CA3 / R30: un método no permitido en la ruta de acuerdo devuelve HTTP 405 y
// no ejecuta el handler.
func TestUS13_Fase2_Poker_MetodoNoPermitido405(t *testing.T) {
	casos := []struct {
		metodo string
	}{
		{metodo: http.MethodGet},
		{metodo: http.MethodDelete},
	}

	for _, caso := range casos {
		t.Run(caso.metodo, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
			poker := &servicioPokerUS13Falso{}
			handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

			rec := solicitarUS22(handler, caso.metodo, rutaAcordarPokerUS13("HU-01", 1), nil)
			exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)

			if poker.acordarLlamadas != 0 {
				t.Errorf("AcordarEstimacion se llamó %d veces y no debía ejecutarse", poker.acordarLlamadas)
			}
		})
	}
}

// CA1 / D31: una ronda revelada sin votos muestra "La ronda no tiene votos"
// sin mínimo, máximo ni nombres, y el formulario de acuerdo sigue disponible
// (acordar con ronda Revelada sin votos está permitido, R17).
func TestUS13_Fase2_Poker_DiferenciasRondaSinVotos(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	proyectos.listarIntegrantes = integrantesPokerUS12()
	poker := pokerReveladaUS13()
	poker.consulta = consultaReveladaUS13(1)
	poker.diferencias = domain.DiferenciaPoker{}
	poker.diferenciasErr = domain.ErrRondaSinVotos
	poker.historial = []domain.RondaHistorialPoker{
		{Ronda: rondaReveladaUS13(1), Votos: nil},
	}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "La ronda no tiene votos")
	exigirCuerpoNoContieneUS22(t, rec, "Hay diferencia", "mínimo", "máximo")
	exigirCuerpoContieneUS22(t, rec,
		`action="/proyectos/PRJ-001/backlog/HU-01/poker/rondas/1/acordar"`,
		"Acordar estimación",
	)
}

// CA2 / D32: una historia sin rondas no renderiza la sección Historial.
func TestUS13_Fase2_Poker_SinRondasSinHistorial(t *testing.T) {
	proyectos, backlog := nuevoEntornoPokerUS12(registroPokerUS12())
	poker := &servicioPokerUS13Falso{rondaActualExiste: false}
	handler := nuevoHandlerPokerUS13(t, proyectos, backlog, poker)

	rec := solicitarUS22(handler, http.MethodGet, rutaPokerUS12("HU-01"), nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Planning Poker", "Iniciar ronda")
	exigirCuerpoNoContieneUS22(t, rec, "Historial", "votos ocultos", "/acordar")

	if poker.historialLlamadas != 0 {
		t.Errorf("HistorialRondas se llamó %d veces y no debía ejecutarse sin rondas", poker.historialLlamadas)
	}
}
