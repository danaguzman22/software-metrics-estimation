package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// valoresPokerVoto es la serie de valores admitidos para renderizar las
// <option> del select de voto (D33). La web no valida contra ella: si un
// formulario manipulado envía otro valor, lo rechaza el dominio.
var valoresPokerVoto = []int{1, 2, 3, 5, 8, 13, 21}

// pokerVista es la representación de la página de Planning Poker lista para
// renderizar. Con la ronda Abierta, Votos queda nil (R33). Diferencias solo se
// completa con la ronda Revelada (D31); Historial siempre que haya rondas
// (D32); StoryPointsActuales lo devuelve el servicio (D36).
type pokerVista struct {
	HayRonda    bool
	RondaNumero int
	Estado      string
	YaVotaron   []string
	Votos       []domain.VotoPoker
	Integrantes []domain.Integrante
	Valores     []int
	Diferencia          *domain.DiferenciaPoker
	DiferenciasSinVotos bool
	Historial           []domain.RondaHistorialPoker
	StoryPointsActuales int
}

// paginaPoker muestra la ronda actual de la historia (la de mayor número) o,
// si no existe ninguna, la posibilidad de iniciar una (R23, R24, R25).
func (a *aplicacion) paginaPoker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	a.renderizarPoker(w, id, itemID, http.StatusOK, "")
}

// iniciarRondaPoker delega el inicio en el servicio y redirige a la página
// con 303 (R26).
func (a *aplicacion) iniciarRondaPoker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	if _, err := a.deps.Poker.IniciarRonda(id, itemID); err != nil {
		a.responderErrorPokerPost(w, id, itemID, "iniciar ronda de poker", err)
		return
	}

	redirigirPoker(w, id, itemID)
}

// votarPoker delega el voto en el servicio y redirige a la página con 303
// (R27). El valor no numérico se rechaza con 422 sin llamar al servicio.
func (a *aplicacion) votarPoker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número de la ronda debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		a.responderErrorRuta(w, id, "No se pudo leer el formulario.")
		return
	}

	valor, err := parsearEntero(strings.TrimSpace(r.PostFormValue("valor")), "El voto debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	integrante := r.PostFormValue("integrante")
	if err := a.deps.Poker.Votar(id, itemID, numero, integrante, valor); err != nil {
		a.responderErrorPokerPost(w, id, itemID, "votar en poker", err)
		return
	}

	redirigirPoker(w, id, itemID)
}

// revelarPoker delega el revelado en el servicio y redirige a la página con
// 303 (R27).
func (a *aplicacion) revelarPoker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número de la ronda debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		a.responderErrorRuta(w, id, "No se pudo leer el formulario.")
		return
	}

	integrante := r.PostFormValue("integrante")
	if err := a.deps.Poker.RevelarRonda(id, itemID, numero, integrante); err != nil {
		a.responderErrorPokerPost(w, id, itemID, "revelar ronda de poker", err)
		return
	}

	redirigirPoker(w, id, itemID)
}

// acordarPoker delega el acuerdo en el servicio y redirige a la página con
// 303 (D28). El valor no numérico se rechaza con 422 sin llamar al servicio.
func (a *aplicacion) acordarPoker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número de la ronda debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		a.responderErrorRuta(w, id, "No se pudo leer el formulario.")
		return
	}

	valor, err := parsearEntero(strings.TrimSpace(r.PostFormValue("valor")), "El valor acordado debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	integrante := r.PostFormValue("integrante")
	if _, err := a.deps.Poker.AcordarEstimacion(id, itemID, numero, integrante, valor); err != nil {
		a.responderErrorPokerPost(w, id, itemID, "acordar estimación en poker", err)
		return
	}

	redirigirPoker(w, id, itemID)
}

// redirigirPoker responde 303 a la propia página del poker (D25).
func redirigirPoker(w http.ResponseWriter, proyectoID, itemID string) {
	w.Header().Set("Location", "/proyectos/"+proyectoID+"/backlog/"+itemID+"/poker")
	w.WriteHeader(http.StatusSeeOther)
}

// renderizarPoker muestra la página de poker con la ronda actual (si la hay) y
// un mensaje opcional. Con la ronda Abierta lista solo los nombres de quienes
// votaron; con la Revelada incluye los votos con sus valores (R33), las
// diferencias y el formulario de acuerdo. Siempre informa los Story Points
// vigentes y, con rondas, el historial (D31, D32, D33, D36).
func (a *aplicacion) renderizarPoker(w http.ResponseWriter, proyectoID, itemID string, estado int, mensaje string) {
	ronda, existe, err := a.deps.Poker.RondaActual(proyectoID, itemID)
	if err != nil {
		a.responderErrorPoker(w, proyectoID, err)
		return
	}

	historia, err := a.deps.Poker.HistoriaActual(proyectoID, itemID)
	if err != nil {
		a.responderErrorPoker(w, proyectoID, err)
		return
	}

	integrantes, err := a.deps.Proyectos.ListarIntegrantes(proyectoID)
	if err != nil {
		a.responderErrorPoker(w, proyectoID, err)
		return
	}

	vista := pokerVista{
		StoryPointsActuales: historia.StoryPoints,
		Integrantes:         integrantes,
		Valores:             valoresPokerVoto,
	}
	if existe {
		consulta, err := a.deps.Poker.ConsultarRonda(proyectoID, itemID, ronda.Numero)
		if err != nil {
			a.responderErrorPoker(w, proyectoID, err)
			return
		}
		vista.HayRonda = true
		vista.RondaNumero = ronda.Numero
		vista.Estado = ronda.Estado

		historial, err := a.deps.Poker.HistorialRondas(proyectoID, itemID)
		if err != nil {
			a.responderErrorPoker(w, proyectoID, err)
			return
		}
		vista.Historial = historial

		if ronda.Estado == domain.EstadoRondaRevelada {
			vista.Votos = consulta.Votos
			diferencia, err := a.deps.Poker.DiferenciasRonda(proyectoID, itemID, ronda.Numero)
			switch {
			case err == nil:
				vista.Diferencia = &diferencia
			case errors.Is(err, domain.ErrRondaSinVotos):
				vista.DiferenciasSinVotos = true
			case errors.Is(err, domain.ErrRondaNoRevelada):
			default:
				a.responderErrorPoker(w, proyectoID, err)
				return
			}
		} else {
			vista.YaVotaron = consulta.YaVotaron
		}
	}

	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorPoker(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "poker", datosPagina{
		Titulo:   "Planning Poker",
		Nav:      navConProyecto(proyectoID),
		Proyecto: &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		ItemID:   itemID,
		Poker:    &vista,
		Error:    mensaje,
		Migas: migasDeProyecto(proyecto.Proyecto.Nombre, proyectoID,
			Miga{Texto: "Backlog", URL: "/proyectos/" + proyectoID + "/backlog"},
			Miga{Texto: "Planning Poker", URL: "/proyectos/" + proyectoID + "/backlog/" + itemID + "/poker"},
		),
		Volver: "/proyectos/" + proyectoID + "/backlog#item-" + itemID,
	})
}

// responderErrorPoker muestra 404 con el mensaje del dominio para un proyecto,
// una historia o una ronda inexistente y 500 genérico para el resto de los
// errores.
func (a *aplicacion) responderErrorPoker(w http.ResponseWriter, proyectoID string, err error) {
	switch {
	case errors.Is(err, domain.ErrProyectoNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Proyecto no encontrado",
			Nav:    navSinProyecto(false),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrItemNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Ítem no encontrado",
			Nav:    navConProyecto(proyectoID),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrRondaNoEncontrada):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Ronda no encontrada",
			Nav:    navConProyecto(proyectoID),
			Error:  err.Error(),
		})
	default:
		a.responderErrorInfraestructura(w, "operar el poker", err)
	}
}

// responderErrorPokerPost traduce el error del servicio tras un POST: 404 para
// proyecto, historia o ronda inexistente; 422 re-renderizando la página con el
// mensaje exacto del dominio para los errores de dominio; 500 genérico para el
// resto.
func (a *aplicacion) responderErrorPokerPost(w http.ResponseWriter, proyectoID, itemID, contexto string, err error) {
	switch {
	case errors.Is(err, domain.ErrProyectoNoEncontrado),
		errors.Is(err, domain.ErrItemNoEncontrado),
		errors.Is(err, domain.ErrRondaNoEncontrada):
		a.responderErrorPoker(w, proyectoID, err)
	case errors.Is(err, domain.ErrRondaRevelada),
		errors.Is(err, domain.ErrRondaAbiertaExistente),
		errors.Is(err, domain.ErrRondaNoRevelada),
		errors.Is(err, domain.ErrStoryPointsInvalidos),
		errors.Is(err, domain.ErrIntegranteNoDelProyecto),
		errors.Is(err, domain.ErrValorVotoPokerInvalido),
		errors.Is(err, domain.ErrNombreIntegranteVacio),
		errors.Is(err, domain.ErrNumeroRondaInvalido):
		a.renderizarPoker(w, proyectoID, itemID, http.StatusUnprocessableEntity, err.Error())
	default:
		a.responderErrorInfraestructura(w, contexto, err)
	}
}
