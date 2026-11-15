package service

import (
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// RepositorioRondas persiste el agregado planning_poker (solo PostgreSQL, D1).
type RepositorioRondas interface {
	// Iniciar inserta la ronda numero para la historia en estado Abierta.
	// domain.ErrRondaAbiertaExistente si la historia ya tiene una ronda Abierta
	// (defensa ante carreras; también traducida desde 23505 del índice parcial).
	Iniciar(proyectoID, itemID string, numero int) error

	// Obtener devuelve la ronda numero de la historia.
	// domain.ErrRondaNoEncontrada si no existe.
	Obtener(proyectoID, itemID string, numero int) (domain.RondaPoker, error)

	// ListarPorHistoria devuelve las rondas de la historia ordenadas por numero.
	ListarPorHistoria(proyectoID, itemID string) ([]domain.RondaPoker, error)

	// RegistrarVoto hace upsert del voto en la ronda, dentro de una transacción,
	// solo si está Abierta. domain.ErrRondaNoEncontrada si no existe;
	// domain.ErrRondaRevelada si ya está Revelada.
	RegistrarVoto(proyectoID, itemID string, numero int, voto domain.VotoPoker) error

	// YaVotaron devuelve los nombres de quienes votaron en la ronda, ordenados
	// alfabéticamente sin distinguir mayúsculas (D20), NUNCA los valores.
	YaVotaron(proyectoID, itemID string, numero int) ([]string, error)

	// Votos devuelve los votos de la ronda con sus valores, ordenados
	// alfabéticamente sin distinguir mayúsculas (D20). Solo se invoca con la
	// ronda Revelada; el servicio lo decide (R6, R17).
	Votos(proyectoID, itemID string, numero int) ([]domain.VotoPoker, error)

	// Revelar pasa la ronda a Revelada (irreversible, D6). domain.ErrRondaNoEncontrada
	// si no existe; domain.ErrRondaRevelada si ya estaba Revelada. Transacción (R11).
	Revelar(proyectoID, itemID string, numero int) error
}

// ServicioPlanningPoker orquesta los casos de uso de Planning Poker (D19).
type ServicioPlanningPoker struct {
	proyectos   RepositorioProyectos
	integrantes RepositorioIntegrantes
	backlog     RepositorioBacklog
	rondas      RepositorioRondas
}

// NuevoServicioPlanningPoker crea el servicio con los repositorios indicados.
func NuevoServicioPlanningPoker(
	proyectos RepositorioProyectos,
	integrantes RepositorioIntegrantes,
	backlog RepositorioBacklog,
	rondas RepositorioRondas,
) *ServicioPlanningPoker {
	return &ServicioPlanningPoker{
		proyectos:   proyectos,
		integrantes: integrantes,
		backlog:     backlog,
		rondas:      rondas,
	}
}

// verificarProyecto verifica que el proyecto exista; devuelve
// domain.ErrProyectoNoEncontrado sin envolver si no existe.
func (s *ServicioPlanningPoker) verificarProyecto(proyectoID string) error {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ErrProyectoNoEncontrado
	}
	return nil
}

// obtenerHistoria verifica que la historia exista; los errores de dominio
// conocidos pasan sin envolver vía envolverErrorRepositorio.
func (s *ServicioPlanningPoker) obtenerHistoria(proyectoID, itemID string) error {
	if _, err := s.backlog.ObtenerItem(proyectoID, itemID); err != nil {
		return envolverErrorRepositorio("obtener ítem del backlog", err)
	}
	return nil
}

// integranteCanonico valida el nombre del integrante al inicio del paso 3,
// antes de cualquier consulta; si está vacío o compuesto solo por espacios
// devuelve domain.ErrNombreIntegranteVacio (E12). Verifica la pertenencia al
// proyecto con comparación recortada sin mayúsculas y devuelve el nombre
// canónico almacenado (R13).
func (s *ServicioPlanningPoker) integranteCanonico(proyectoID, integranteNombre string) (string, error) {
	nombre := strings.TrimSpace(integranteNombre)
	if nombre == "" {
		return "", domain.ErrNombreIntegranteVacio
	}
	integrantes, err := s.integrantes.ListarPorProyecto(proyectoID)
	if err != nil {
		return "", envolverErrorRepositorio("listar integrantes", err)
	}
	registrado, encontrado := domain.BuscarIntegrante(integrantes, nombre)
	if !encontrado {
		return "", domain.ErrIntegranteNoDelProyecto
	}
	return registrado.Nombre, nil
}

// IniciarRonda crea la siguiente ronda de la historia en estado Abierta.
//
// Orden de evaluación: proyecto, historia, ronda. Si la historia tiene una
// ronda Abierta pendiente devuelve domain.ErrRondaAbiertaExistente.
func (s *ServicioPlanningPoker) IniciarRonda(proyectoID, itemID string) (domain.RondaPoker, error) {
	if err := s.verificarProyecto(proyectoID); err != nil {
		return domain.RondaPoker{}, err
	}
	if err := s.obtenerHistoria(proyectoID, itemID); err != nil {
		return domain.RondaPoker{}, err
	}
	rondas, err := s.rondas.ListarPorHistoria(proyectoID, itemID)
	if err != nil {
		return domain.RondaPoker{}, envolverErrorRepositorio("listar rondas de la historia", err)
	}
	numero, err := domain.SiguienteNumeroRonda(rondas)
	if err != nil {
		return domain.RondaPoker{}, err
	}
	if err := s.rondas.Iniciar(proyectoID, itemID, numero); err != nil {
		return domain.RondaPoker{}, envolverErrorRepositorio("iniciar ronda", err)
	}
	return domain.RondaPoker{Numero: numero, Estado: domain.EstadoRondaAbierta}, nil
}

// Votar registra o reemplaza el voto del integrante en la ronda Abierta.
//
// Orden de evaluación: proyecto, historia, integrante del proyecto, ronda
// existente y Abierta, valor. El voto persiste el nombre canónico almacenado
// del integrante (R13).
func (s *ServicioPlanningPoker) Votar(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) error {
	if err := s.verificarProyecto(proyectoID); err != nil {
		return err
	}
	if err := s.obtenerHistoria(proyectoID, itemID); err != nil {
		return err
	}
	nombreCanonico, err := s.integranteCanonico(proyectoID, integranteNombre)
	if err != nil {
		return err
	}
	if rondaNumero <= 0 {
		return domain.ErrNumeroRondaInvalido
	}
	ronda, err := s.rondas.Obtener(proyectoID, itemID, rondaNumero)
	if err != nil {
		return envolverErrorRepositorio("obtener ronda", err)
	}
	if ronda.Estado == domain.EstadoRondaRevelada {
		return domain.ErrRondaRevelada
	}
	voto, err := domain.NuevoVotoPoker(nombreCanonico, valor)
	if err != nil {
		return err
	}
	if err := s.rondas.RegistrarVoto(proyectoID, itemID, rondaNumero, voto); err != nil {
		return envolverErrorRepositorio("registrar voto", err)
	}
	return nil
}

// RevelarRonda pasa la ronda Abierta a Revelada; es irreversible (D6).
// Cualquier integrante del proyecto puede revelar.
func (s *ServicioPlanningPoker) RevelarRonda(proyectoID, itemID string, rondaNumero int, integranteNombre string) error {
	if err := s.verificarProyecto(proyectoID); err != nil {
		return err
	}
	if err := s.obtenerHistoria(proyectoID, itemID); err != nil {
		return err
	}
	if _, err := s.integranteCanonico(proyectoID, integranteNombre); err != nil {
		return err
	}
	if rondaNumero <= 0 {
		return domain.ErrNumeroRondaInvalido
	}
	ronda, err := s.rondas.Obtener(proyectoID, itemID, rondaNumero)
	if err != nil {
		return envolverErrorRepositorio("obtener ronda", err)
	}
	if ronda.Estado == domain.EstadoRondaRevelada {
		return domain.ErrRondaRevelada
	}
	if err := s.rondas.Revelar(proyectoID, itemID, rondaNumero); err != nil {
		return envolverErrorRepositorio("revelar ronda", err)
	}
	return nil
}

// ConsultarRonda devuelve el estado de la ronda y, según el estado:
// con la ronda Abierta, YaVotaron (nombres; Votos nil) y nunca los valores;
// con la ronda Revelada, Votos (YaVotaron nil) (R6, R17).
func (s *ServicioPlanningPoker) ConsultarRonda(proyectoID, itemID string, rondaNumero int) (domain.ConsultaRondaPoker, error) {
	if err := s.verificarProyecto(proyectoID); err != nil {
		return domain.ConsultaRondaPoker{}, err
	}
	if err := s.obtenerHistoria(proyectoID, itemID); err != nil {
		return domain.ConsultaRondaPoker{}, err
	}
	if rondaNumero <= 0 {
		return domain.ConsultaRondaPoker{}, domain.ErrNumeroRondaInvalido
	}
	ronda, err := s.rondas.Obtener(proyectoID, itemID, rondaNumero)
	if err != nil {
		return domain.ConsultaRondaPoker{}, envolverErrorRepositorio("obtener ronda", err)
	}
	if ronda.Estado == domain.EstadoRondaRevelada {
		votos, err := s.rondas.Votos(proyectoID, itemID, rondaNumero)
		if err != nil {
			return domain.ConsultaRondaPoker{}, envolverErrorRepositorio("listar votos de la ronda", err)
		}
		return domain.ConsultaRondaPoker{Ronda: ronda, Votos: votos}, nil
	}
	nombres, err := s.rondas.YaVotaron(proyectoID, itemID, rondaNumero)
	if err != nil {
		return domain.ConsultaRondaPoker{}, envolverErrorRepositorio("listar quiénes votaron", err)
	}
	return domain.ConsultaRondaPoker{Ronda: ronda, YaVotaron: nombres}, nil
}
