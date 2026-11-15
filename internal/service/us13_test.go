package service

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// llamadaEstimador registra una invocación al doble fakeEstimadorStoryPoints.
type llamadaEstimador struct {
	proyectoID  string
	itemID      string
	storyPoints int
}

// fakeEstimadorStoryPoints es el doble de EstimadorStoryPoints para los
// tests de US-13 (D14): registra los llamados a EstimarItem y permite
// inyectar un resultado o un error.
type fakeEstimadorStoryPoints struct {
	llamadas  []llamadaEstimador
	resultado domain.ItemBacklog
	err       error
}

func (f *fakeEstimadorStoryPoints) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
	f.llamadas = append(f.llamadas, llamadaEstimador{proyectoID: proyectoID, itemID: itemID, storyPoints: storyPoints})
	if f.err != nil {
		return domain.ItemBacklog{}, f.err
	}
	if f.resultado.ID != "" {
		return f.resultado, nil
	}
	return domain.ItemBacklog{ID: itemID, StoryPoints: storyPoints}, nil
}

func (e *entornoPoker) servicioConEstimador(estimador EstimadorStoryPoints) *ServicioPlanningPoker {
	return NuevoServicioPlanningPokerConEstimador(
		e.fakeRepositorioProyectos,
		e.fakeRepositorioIntegrantes,
		e.fakeRepositorioBacklog,
		e.fakeRepositorioRondas,
		estimador,
	)
}

func (e *entornoPoker) servicioConEstimadorReal() *ServicioPlanningPoker {
	return NuevoServicioPlanningPokerConEstimador(
		e.fakeRepositorioProyectos,
		e.fakeRepositorioIntegrantes,
		e.fakeRepositorioBacklog,
		e.fakeRepositorioRondas,
		NuevoServicioBacklog(e.fakeRepositorioProyectos, e.fakeRepositorioBacklog),
	)
}

func historiaStoryPointsPoker(t *testing.T, entorno *entornoPoker, proyectoID, itemID string) int {
	t.Helper()
	item, err := entorno.fakeRepositorioBacklog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		t.Fatalf("ObtenerItem(%q, %q) devolvió un error inesperado: %v", proyectoID, itemID, err)
	}
	return item.StoryPoints
}

// R4 / CA2 (US-13): el historial informa las rondas reveladas con sus votos y
// las abiertas con Votos nil, en orden ascendente por número.
func TestUS13_Servicio_HistorialRondas_ReveladasConVotosYAbiertasSinValores(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 3, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaAbierta})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
	)
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 3,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 13},
	)
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 2,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 21},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	historial, err := servicio.HistorialRondas("p1", "HU-01")
	if err != nil {
		t.Fatalf("HistorialRondas() devolvió un error inesperado: %v", err)
	}
	if len(historial) != 3 {
		t.Fatalf("len(historial) = %d, se esperaban 3 rondas", len(historial))
	}
	if historial[0].Ronda != (domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada}) {
		t.Errorf("historial[0].Ronda = %+v, se esperaba {Numero:1 Estado:Revelada}", historial[0].Ronda)
	}
	if historial[1].Ronda != (domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaAbierta}) {
		t.Errorf("historial[1].Ronda = %+v, se esperaba {Numero:2 Estado:Abierta}", historial[1].Ronda)
	}
	if historial[2].Ronda != (domain.RondaPoker{Numero: 3, Estado: domain.EstadoRondaRevelada}) {
		t.Errorf("historial[2].Ronda = %+v, se esperaba {Numero:3 Estado:Revelada}", historial[2].Ronda)
	}
	if !reflect.DeepEqual(historial[0].Votos, []domain.VotoPoker{
		{Integrante: "Dana Guzmán", Valor: 5},
		{Integrante: "Victoria Caro", Valor: 8},
	}) {
		t.Errorf("historial[0].Votos = %+v, se esperaban los votos ordenados de la ronda 1", historial[0].Votos)
	}
	if historial[1].Votos != nil {
		t.Errorf("historial[1].Votos = %+v, se esperaba nil con la ronda Abierta (D4, R15)", historial[1].Votos)
	}
	if !reflect.DeepEqual(historial[2].Votos, []domain.VotoPoker{
		{Integrante: "Dana Guzmán", Valor: 13},
	}) {
		t.Errorf("historial[2].Votos = %+v, se esperaban los votos de la ronda 3", historial[2].Votos)
	}
	if entorno.votosLlamadas != 2 {
		t.Errorf("Votos fue invocado %d veces, se esperaban 2 (solo las rondas Revelada)", entorno.votosLlamadas)
	}
}

// R4 / E13 / CA2 (US-13): el historial nunca expone valores de una ronda
// Abierta, aunque la ronda tenga votos sembrados.
func TestUS13_Servicio_HistorialRondas_AbiertaNuncaExponeValores(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	historial, err := servicio.HistorialRondas("p1", "HU-01")
	if err != nil {
		t.Fatalf("HistorialRondas() devolvió un error inesperado: %v", err)
	}
	if len(historial) != 1 {
		t.Fatalf("len(historial) = %d, se esperaba 1", len(historial))
	}
	if historial[0].Votos != nil {
		t.Errorf("Votos = %+v, se esperaba nil con la ronda Abierta (D4, R15)", historial[0].Votos)
	}
	if entorno.votosLlamadas != 0 {
		t.Errorf("Votos fue invocado %d veces; con la ronda Abierta NUNCA debe llamarse", entorno.votosLlamadas)
	}
}

// R4 / CA2 (US-13): una historia sin rondas devuelve una lista vacía sin error.
func TestUS13_Servicio_HistorialRondas_HistoriaSinRondasListaVacia(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	historial, err := servicio.HistorialRondas("p1", "HU-01")
	if err != nil {
		t.Fatalf("HistorialRondas() devolvió un error inesperado: %v", err)
	}
	if len(historial) != 0 {
		t.Errorf("len(historial) = %d, se esperaba 0", len(historial))
	}
	if entorno.votosLlamadas != 0 {
		t.Errorf("Votos fue invocado %d veces; sin rondas no debe llamarse", entorno.votosLlamadas)
	}
}

// V1 / R13 (US-13): proyecto inexistente prevalece sobre la historia en el historial.
func TestUS13_Servicio_HistorialRondas_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	historial, err := servicio.HistorialRondas("inexistente", "HU-01")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if historial != nil {
		t.Errorf("historial = %+v, se esperaba nil", historial)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 || entorno.fakeRepositorioRondas.listarLlamadas != 0 {
		t.Errorf("se llamó a backlog/rondas (obtener=%d, listar=%d); el proyecto debía validarse primero",
			entorno.fakeRepositorioBacklog.obtenerLlamadas, entorno.fakeRepositorioRondas.listarLlamadas)
	}
}

// V2 / R13 (US-13): historia inexistente devuelve domain.ErrItemNoEncontrado.
func TestUS13_Servicio_HistorialRondas_HistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	historial, err := servicio.HistorialRondas("p1", "HU-99")
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if err != domain.ErrItemNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
	}
	if historial != nil {
		t.Errorf("historial = %+v, se esperaba nil", historial)
	}
	if entorno.fakeRepositorioRondas.listarLlamadas != 0 {
		t.Errorf("ListarPorHistoria fue invocado %d veces; la historia debía validarse primero", entorno.fakeRepositorioRondas.listarLlamadas)
	}
}

// E11 / R11 / V10 (US-13): los errores de repositorio no conocidos se envuelven
// con %w (prefijo "service:") y conservan la causa con errors.Is.
func TestUS13_Servicio_HistorialRondas_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar las rondas", func(e *entornoPoker) { e.fakeRepositorioRondas.listarErr = errSimulado }},
		{"error al listar los votos", func(e *entornoPoker) { e.votosErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			caso.preparar(entorno)
			servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

			historial, err := servicio.HistorialRondas("p1", "HU-01")
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if err == errSimulado {
				t.Errorf("err = %v, se esperaba un error envuelto, no el error del repositorio tal cual", err)
			}
			if !strings.HasPrefix(err.Error(), "service:") {
				t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
			}
			if historial != nil {
				t.Errorf("historial = %+v, se esperaba nil", historial)
			}
		})
	}
}

// R1 / CA1 (US-13): DiferenciasRonda con la ronda Revelada informa si hay
// diferencia, mínimo, máximo y quiénes votaron cada extremo.
func TestUS13_Servicio_DiferenciasRonda_ReveladaConVotos(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Carlos", Valor: 13},
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 2},
		domain.VotoPoker{Integrante: "Ana", Valor: 8},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("DiferenciasRonda() devolvió un error inesperado: %v", err)
	}
	if !diferencia.HayDiferencia {
		t.Errorf("HayDiferencia = false, se esperaba true")
	}
	if diferencia.Minimo != 2 {
		t.Errorf("Minimo = %d, se esperaba 2", diferencia.Minimo)
	}
	if diferencia.Maximo != 13 {
		t.Errorf("Maximo = %d, se esperaba 13", diferencia.Maximo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMinimo, []string{"Dana Guzmán"}) {
		t.Errorf("QuienesMinimo = %v, se esperaban [Dana Guzmán]", diferencia.QuienesMinimo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMaximo, []string{"Carlos"}) {
		t.Errorf("QuienesMaximo = %v, se esperaban [Carlos]", diferencia.QuienesMaximo)
	}
}

// R2 / CA1 (US-13): con todos los votos iguales no hay diferencia.
func TestUS13_Servicio_DiferenciasRonda_ReveladaTodosIguales(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 8},
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("DiferenciasRonda() devolvió un error inesperado: %v", err)
	}
	if diferencia.HayDiferencia {
		t.Errorf("HayDiferencia = true, se esperaba false")
	}
	if diferencia.Minimo != 8 || diferencia.Maximo != 8 {
		t.Errorf("Minimo/Maximo = %d/%d, se esperaba 8/8", diferencia.Minimo, diferencia.Maximo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMinimo, []string{"Dana Guzmán", "Victoria Caro"}) {
		t.Errorf("QuienesMinimo = %v, se esperaban [Dana Guzmán Victoria Caro]", diferencia.QuienesMinimo)
	}
}

// R2 / CA1 (US-13): con un solo voto no hay diferencia y el votante figura en
// ambos extremos.
func TestUS13_Servicio_DiferenciasRonda_ReveladaUnSoloVoto(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 13},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("DiferenciasRonda() devolvió un error inesperado: %v", err)
	}
	if diferencia.HayDiferencia {
		t.Errorf("HayDiferencia = true, se esperaba false")
	}
	if diferencia.Minimo != 13 || diferencia.Maximo != 13 {
		t.Errorf("Minimo/Maximo = %d/%d, se esperaba 13/13", diferencia.Minimo, diferencia.Maximo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMinimo, []string{"Dana Guzmán"}) {
		t.Errorf("QuienesMinimo = %v, se esperaban [Dana Guzmán]", diferencia.QuienesMinimo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMaximo, []string{"Dana Guzmán"}) {
		t.Errorf("QuienesMaximo = %v, se esperaban [Dana Guzmán]", diferencia.QuienesMaximo)
	}
}

// R1 / R15 / D3 / D8 / V9 / E9 (US-13): con la ronda Abierta no se informa
// nada: domain.ErrRondaNoRevelada; nunca se consultan los votos.
func TestUS13_Servicio_DiferenciasRonda_AbiertaDevuelveErrRondaNoRevelada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
	)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
	if !errors.Is(err, domain.ErrRondaNoRevelada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoRevelada)
	}
	if err != domain.ErrRondaNoRevelada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoRevelada)
	}
	if diferencia != (domain.DiferenciaPoker{}) {
		t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
	}
	if entorno.votosLlamadas != 0 {
		t.Errorf("Votos fue invocado %d veces; con la ronda Abierta NUNCA debe llamarse (R1, R15)", entorno.votosLlamadas)
	}
}

// R3 / V8 / E8 (US-13): con la ronda Revelada y cero votos se devuelve
// domain.ErrRondaSinVotos sin envolver; no se informa mínimo ni máximo.
func TestUS13_Servicio_DiferenciasRonda_ReveladaSinVotosDevuelveErrRondaSinVotos(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
	if !errors.Is(err, domain.ErrRondaSinVotos) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaSinVotos)
	}
	if err != domain.ErrRondaSinVotos {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaSinVotos)
	}
	if diferencia != (domain.DiferenciaPoker{}) {
		t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
	}
}

// V1 / R13 (US-13): proyecto inexistente prevalece sobre la historia y la ronda.
func TestUS13_Servicio_DiferenciasRonda_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("inexistente", "HU-01", 1)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if diferencia != (domain.DiferenciaPoker{}) {
		t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 || entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
		t.Errorf("se llamó a backlog/rondas (obtener=%d, obtener ronda=%d); el proyecto debía validarse primero",
			entorno.fakeRepositorioBacklog.obtenerLlamadas, entorno.fakeRepositorioRondas.obtenerLlamadas)
	}
}

// V2 / R13 (US-13): historia inexistente prevalece sobre la ronda.
func TestUS13_Servicio_DiferenciasRonda_HistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-99", 1)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if err != domain.ErrItemNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
	}
	if diferencia != (domain.DiferenciaPoker{}) {
		t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
	}
	if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
		t.Errorf("Obtener fue invocado %d veces; la historia debía validarse primero", entorno.fakeRepositorioRondas.obtenerLlamadas)
	}
}

// R14 / V4 / E4 (US-13): rondaNumero <= 0 devuelve
// domain.ErrNumeroRondaInvalido en el paso 3, antes de rondas.Obtener.
func TestUS13_Servicio_DiferenciasRonda_NumeroRondaInvalido(t *testing.T) {
	for _, numero := range []int{0, -1} {
		t.Run("numero "+strconv.Itoa(numero), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

			diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", numero)
			if !errors.Is(err, domain.ErrNumeroRondaInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroRondaInvalido)
			}
			if err != domain.ErrNumeroRondaInvalido {
				t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrNumeroRondaInvalido)
			}
			if diferencia != (domain.DiferenciaPoker{}) {
				t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
			}
			if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
				t.Errorf("Obtener fue invocado %d veces; nunca se debía llegar a rondas.Obtener con numero <= 0", entorno.fakeRepositorioRondas.obtenerLlamadas)
			}
		})
	}
}

// V5 / E5 (US-13): ronda inexistente devuelve domain.ErrRondaNoEncontrada.
func TestUS13_Servicio_DiferenciasRonda_RondaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

	diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 7)
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err != domain.ErrRondaNoEncontrada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoEncontrada)
	}
	if diferencia != (domain.DiferenciaPoker{}) {
		t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
	}
}

// E11 / R11 / V10 (US-13): los errores de repositorio no conocidos se envuelven
// con %w (prefijo "service:") y conservan la causa con errors.Is.
func TestUS13_Servicio_DiferenciasRonda_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al obtener la ronda", func(e *entornoPoker) { e.fakeRepositorioRondas.obtenerErr = errSimulado }},
		{"error al listar los votos", func(e *entornoPoker) { e.votosErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			caso.preparar(entorno)
			servicio := entorno.servicioConEstimador(&fakeEstimadorStoryPoints{})

			diferencia, err := servicio.DiferenciasRonda("p1", "HU-01", 1)
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if err == errSimulado {
				t.Errorf("err = %v, se esperaba un error envuelto, no el error del repositorio tal cual", err)
			}
			if !strings.HasPrefix(err.Error(), "service:") {
				t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
			}
			if diferencia != (domain.DiferenciaPoker{}) {
				t.Errorf("diferencia = %+v, se esperaba domain.DiferenciaPoker{}", diferencia)
			}
		})
	}
}

// R5 / R7 / D5 / CA3 (US-13): acordar sobre una ronda Revelada invoca el
// estimador (caso de uso de estimar existente) con exactamente proyectoID,
// itemID y valor, y la historia queda con esos Story Points.
func TestUS13_Servicio_AcordarEstimacion_Exito(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
	)
	estimador := &fakeEstimadorStoryPoints{resultado: domain.ItemBacklog{ID: "HU-01", StoryPoints: 8}}
	servicio := entorno.servicioConEstimador(estimador)

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
	if err != nil {
		t.Fatalf("AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if len(estimador.llamadas) != 1 {
		t.Fatalf("EstimarItem fue invocado %d veces, se esperaba 1", len(estimador.llamadas))
	}
	llamada := estimador.llamadas[0]
	if llamada != (llamadaEstimador{proyectoID: "p1", itemID: "HU-01", storyPoints: 8}) {
		t.Errorf("EstimarItem recibió %+v, se esperaba {p1 HU-01 8}", llamada)
	}
	if historia.StoryPoints != 8 {
		t.Errorf("historia.StoryPoints = %d, se esperaba 8", historia.StoryPoints)
	}
}

// D14 / R7 (US-13): AcordarEstimacion delega en EstimadorStoryPoints.EstimarItem con los argumentos exactos.
func TestUS13_ServicioPlanningPoker_AcordarEstimacionDelegaEnEstimador(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaRevelada})
	estimador := &fakeEstimadorStoryPoints{resultado: domain.ItemBacklog{ID: "HU-01", Titulo: "Historia HU-01", StoryPoints: 13}}
	servicio := entorno.servicioConEstimador(estimador)

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 2, "Victoria Caro", 13)
	if err != nil {
		t.Fatalf("AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if len(estimador.llamadas) != 1 {
		t.Fatalf("EstimarItem fue invocado %d veces, se esperaba 1", len(estimador.llamadas))
	}
	if estimador.llamadas[0] != (llamadaEstimador{proyectoID: "p1", itemID: "HU-01", storyPoints: 13}) {
		t.Errorf("EstimarItem recibió %+v, se esperaba {p1 HU-01 13}", estimador.llamadas[0])
	}
	if historia.StoryPoints != 13 {
		t.Errorf("historia.StoryPoints = %d, se esperaba 13", historia.StoryPoints)
	}
}

// R7 / D5 / CA3 (US-13): el acuerdo guarda Story Points en la historia
// reutilizando el caso de uso existente (ServicioBacklog.EstimarItem).
func TestUS13_Servicio_AcordarEstimacion_PersisteStoryPointsViaCasoDeUsoExistente(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	entorno.items["p1"] = append(entorno.items["p1"], domain.ItemBacklog{
		ID:          "HU-02",
		Titulo:      "Historia HU-02",
		Descripcion: "Descripción HU-02",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: 0,
	})
	sembrarRondaPoker(t, entorno, "p1", "HU-02", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicioConEstimadorReal()

	historia, err := servicio.AcordarEstimacion("p1", "HU-02", 1, "Dana Guzmán", 13)
	if err != nil {
		t.Fatalf("AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if historia.StoryPoints != 13 {
		t.Errorf("historia.StoryPoints = %d, se esperaba 13", historia.StoryPoints)
	}
	if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-02"); obtenidos != 13 {
		t.Errorf("StoryPoints persistidos = %d, se esperaban 13", obtenidos)
	}
}

// R8 / D5 / CA3 (US-13): re-acordar en una ronda posterior reemplaza el valor
// anterior; una historia ya estimada se reestima sin error (R18).
func TestUS13_Servicio_AcordarEstimacion_ReAcuerdoEnRondaPosteriorReemplaza(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicioConEstimadorReal()

	if _, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 5); err != nil {
		t.Fatalf("primer AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-01"); obtenidos != 5 {
		t.Fatalf("StoryPoints tras el primer acuerdo = %d, se esperaban 5", obtenidos)
	}
	if _, err := servicio.AcordarEstimacion("p1", "HU-01", 2, "Victoria Caro", 8); err != nil {
		t.Fatalf("segundo AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-01"); obtenidos != 8 {
		t.Errorf("StoryPoints tras el re-acuerdo = %d, se esperaban 8", obtenidos)
	}
}

// D5 / R6 / CA3 (US-13): el valor acordado no tiene que coincidir con ningún voto.
func TestUS13_Servicio_AcordarEstimacion_ValorQueNadieVotoPermitido(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
	)
	estimador := &fakeEstimadorStoryPoints{}
	servicio := entorno.servicioConEstimador(estimador)

	if _, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 13); err != nil {
		t.Fatalf("AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if len(estimador.llamadas) != 1 || estimador.llamadas[0].storyPoints != 13 {
		t.Errorf("estimador recibió %+v, se esperaba una llamada con storyPoints = 13", estimador.llamadas)
	}
}

// R17 / D18 / E8-bis (US-13): con la ronda Revelada sin votos se puede acordar;
// no se devuelve ErrRondaSinVotos.
func TestUS13_Servicio_AcordarEstimacion_ReveladaSinVotosPermitido(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	estimador := &fakeEstimadorStoryPoints{}
	servicio := entorno.servicioConEstimador(estimador)

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
	if err != nil {
		t.Fatalf("AcordarEstimacion() devolvió un error inesperado: %v", err)
	}
	if errors.Is(err, domain.ErrRondaSinVotos) {
		t.Errorf("err = %v; acordar con ronda Revelada sin votos está permitido (R17, D18)", err)
	}
	if len(estimador.llamadas) != 1 {
		t.Fatalf("EstimarItem fue invocado %d veces, se esperaba 1", len(estimador.llamadas))
	}
	if historia.StoryPoints != 8 {
		t.Errorf("historia.StoryPoints = %d, se esperaba 8", historia.StoryPoints)
	}
}

// R5 / R12 / CA4 / V6 / E6 (US-13): con la ronda Abierta se devuelve
// domain.ErrRondaNoRevelada; no se acuerda ni se persiste nada.
func TestUS13_Servicio_AcordarEstimacion_RondaAbiertaDevuelveErrRondaNoRevelada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	estimador := &fakeEstimadorStoryPoints{}
	servicio := entorno.servicioConEstimador(estimador)
	storyPointsPrevios := historiaStoryPointsPoker(t, entorno, "p1", "HU-01")

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
	if !errors.Is(err, domain.ErrRondaNoRevelada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoRevelada)
	}
	if err != domain.ErrRondaNoRevelada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoRevelada)
	}
	if historia != (domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
	if len(estimador.llamadas) != 0 {
		t.Errorf("EstimarItem fue invocado %d veces; con la ronda Abierta no debe escribirse nada", len(estimador.llamadas))
	}
	if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-01"); obtenidos != storyPointsPrevios {
		t.Errorf("StoryPoints tras acordar sobre ronda Abierta = %d, se esperaban %d (sin cambios)", obtenidos, storyPointsPrevios)
	}
}

// R6 / R12 / V7 / E7 (US-13): un valor fuera de la serie devuelve
// domain.ErrStoryPointsInvalidos; no se persiste nada.
func TestUS13_Servicio_AcordarEstimacion_ValorFueraDeSerie(t *testing.T) {
	for _, valor := range []int{0, 4, 6, 22, -1} {
		t.Run("valor "+strconv.Itoa(valor), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			estimador := &fakeEstimadorStoryPoints{}
			servicio := entorno.servicioConEstimador(estimador)
			storyPointsPrevios := historiaStoryPointsPoker(t, entorno, "p1", "HU-01")

			historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", valor)
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
			}
			if err != domain.ErrStoryPointsInvalidos {
				t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrStoryPointsInvalidos)
			}
			if historia != (domain.ItemBacklog{}) {
				t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
			}
			if len(estimador.llamadas) != 0 {
				t.Errorf("EstimarItem fue invocado %d veces; con valor fuera de la serie no debe escribirse nada", len(estimador.llamadas))
			}
			if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-01"); obtenidos != storyPointsPrevios {
				t.Errorf("StoryPoints = %d, se esperaban %d (sin cambios)", obtenidos, storyPointsPrevios)
			}
		})
	}
}

// CA4 / V3 / E3 (US-13): un integrante que no pertenece al proyecto no puede
// acordar, ni con la ronda Revelada; no se persiste nada.
func TestUS13_Servicio_AcordarEstimacion_IntegranteNoDelProyecto(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	estimador := &fakeEstimadorStoryPoints{}
	servicio := entorno.servicioConEstimador(estimador)

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Intruso", 8)
	if !errors.Is(err, domain.ErrIntegranteNoDelProyecto) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrIntegranteNoDelProyecto)
	}
	if err != domain.ErrIntegranteNoDelProyecto {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrIntegranteNoDelProyecto)
	}
	if historia != (domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
	if len(estimador.llamadas) != 0 {
		t.Errorf("EstimarItem fue invocado %d veces; no debía escribirse nada", len(estimador.llamadas))
	}
}

// E10 / V3 (US-13): nombre vacío o compuesto de espacios se rechaza al inicio
// del paso 3, antes de ListarPorProyecto.
func TestUS13_Servicio_AcordarEstimacion_NombreVacioOEnBlanco(t *testing.T) {
	for _, nombre := range []string{"", "   ", "\t\n"} {
		t.Run("nombre "+strconv.Itoa(len(nombre)), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			estimador := &fakeEstimadorStoryPoints{}
			servicio := entorno.servicioConEstimador(estimador)

			_, err := servicio.AcordarEstimacion("p1", "HU-01", 1, nombre, 8)
			if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
			}
			if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 {
				t.Errorf("ListarPorProyecto fue invocado %d veces; el nombre debía validarse antes", entorno.fakeRepositorioIntegrantes.listarLlamadas)
			}
			if len(estimador.llamadas) != 0 {
				t.Errorf("EstimarItem fue invocado %d veces; no debía escribirse nada", len(estimador.llamadas))
			}
		})
	}
}

// R14 / V4 / E4 (US-13): rondaNumero <= 0 devuelve
// domain.ErrNumeroRondaInvalido en el paso 4, antes de rondas.Obtener.
func TestUS13_Servicio_AcordarEstimacion_NumeroRondaInvalido(t *testing.T) {
	for _, numero := range []int{0, -1} {
		t.Run("numero "+strconv.Itoa(numero), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			estimador := &fakeEstimadorStoryPoints{}
			servicio := entorno.servicioConEstimador(estimador)

			_, err := servicio.AcordarEstimacion("p1", "HU-01", numero, "Dana Guzmán", 8)
			if !errors.Is(err, domain.ErrNumeroRondaInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroRondaInvalido)
			}
			if err != domain.ErrNumeroRondaInvalido {
				t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrNumeroRondaInvalido)
			}
			if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
				t.Errorf("Obtener fue invocado %d veces; nunca se debía llegar a rondas.Obtener con numero <= 0", entorno.fakeRepositorioRondas.obtenerLlamadas)
			}
			if len(estimador.llamadas) != 0 {
				t.Errorf("EstimarItem fue invocado %d veces; no debía escribirse nada", len(estimador.llamadas))
			}
		})
	}
}

// V5 / E5 (US-13): ronda inexistente devuelve domain.ErrRondaNoEncontrada; no
// se persiste nada.
func TestUS13_Servicio_AcordarEstimacion_RondaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	estimador := &fakeEstimadorStoryPoints{}
	servicio := entorno.servicioConEstimador(estimador)

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 7, "Dana Guzmán", 8)
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err != domain.ErrRondaNoEncontrada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoEncontrada)
	}
	if historia != (domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
	if len(estimador.llamadas) != 0 {
		t.Errorf("EstimarItem fue invocado %d veces; con ronda inexistente no debe escribirse nada", len(estimador.llamadas))
	}
}

// R12 / D6 (US-13): orden de evaluación proyecto → historia → integrante →
// ronda → estado → valor; el error de más arriba prevalece sobre el de más abajo.
func TestUS13_Servicio_AcordarEstimacion_Precedencia(t *testing.T) {
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, e *entornoPoker)
		proyectoID  string
		itemID      string
		rondaNumero int
		nombre      string
		valor       int
		esperado    error
	}{
		{
			descripcion: "proyecto inexistente prevalece sobre historia inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-99", rondaNumero: 1, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre integrante ajeno",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 8,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-01", rondaNumero: 7, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre ronda Abierta",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "inexistente", itemID: "HU-01", rondaNumero: 1, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre valor inválido",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-01", rondaNumero: 1, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "historia inexistente prevalece sobre integrante ajeno",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-99", rondaNumero: 1, nombre: "Intruso", valor: 8,
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion: "historia inexistente prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-99", rondaNumero: 7, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion: "historia inexistente prevalece sobre ronda Abierta",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-99", rondaNumero: 1, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion: "historia inexistente prevalece sobre valor inválido",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-99", rondaNumero: 1, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion: "integrante ajeno prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 7, nombre: "Intruso", valor: 8,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "integrante ajeno prevalece sobre ronda Abierta",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 8,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "integrante ajeno prevalece sobre valor inválido",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 4,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "nombre en blanco prevalece sobre rondaNumero <= 0",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 0, nombre: "   ", valor: 8,
			esperado: domain.ErrNombreIntegranteVacio,
		},
		{
			descripcion: "nombre en blanco prevalece sobre valor inválido",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 1, nombre: "", valor: 4,
			esperado: domain.ErrNombreIntegranteVacio,
		},
		{
			descripcion: "rondaNumero <= 0 prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 0, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrNumeroRondaInvalido,
		},
		{
			descripcion: "ronda inexistente prevalece sobre ronda Abierta",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 7, nombre: "Dana Guzmán", valor: 8,
			esperado: domain.ErrRondaNoEncontrada,
		},
		{
			descripcion: "ronda inexistente prevalece sobre valor inválido",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 7, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrRondaNoEncontrada,
		},
		{
			descripcion: "ronda Abierta prevalece sobre valor inválido",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrRondaNoRevelada,
		},
		{
			descripcion: "ronda Abierta prevalece aunque el valor sea válido",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Dana Guzmán", valor: 21,
			esperado: domain.ErrRondaNoRevelada,
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			caso.preparar(t, entorno)
			estimador := &fakeEstimadorStoryPoints{}
			servicio := entorno.servicioConEstimador(estimador)

			_, err := servicio.AcordarEstimacion(caso.proyectoID, caso.itemID, caso.rondaNumero, caso.nombre, caso.valor)
			if !errors.Is(err, caso.esperado) {
				t.Fatalf("err = %v, se esperaba %v", err, caso.esperado)
			}
			if len(estimador.llamadas) != 0 {
				t.Errorf("EstimarItem fue invocado %d veces; con el error de precedencia no debía escribirse nada", len(estimador.llamadas))
			}
		})
	}
}

// R11 / D15 (US-13): un error de repositorio no conocido en la ronda se envuelve
// con %w (prefijo "service:"); los errores de dominio nuevos pasan sin envolver.
func TestUS13_Servicio_AcordarEstimacion_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar los integrantes", func(e *entornoPoker) { e.fakeRepositorioIntegrantes.listarErr = errSimulado }},
		{"error al obtener la ronda", func(e *entornoPoker) { e.fakeRepositorioRondas.obtenerErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			caso.preparar(entorno)
			estimador := &fakeEstimadorStoryPoints{}
			servicio := entorno.servicioConEstimador(estimador)

			historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if err == errSimulado {
				t.Errorf("err = %v, se esperaba un error envuelto, no el error del repositorio tal cual", err)
			}
			if !strings.HasPrefix(err.Error(), "service:") {
				t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
			}
			if historia != (domain.ItemBacklog{}) {
				t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
			}
			if len(estimador.llamadas) != 0 {
				t.Errorf("EstimarItem fue invocado %d veces; no debía escribirse nada", len(estimador.llamadas))
			}
		})
	}
}

// R7 / R11 (US-13): un error del estimador no conocido se envuelve con %w;
// un error de dominio del estimador pasa sin envolver.
func TestUS13_Servicio_AcordarEstimacion_ErrorDelEstimador(t *testing.T) {
	t.Run("error no conocido se envuelve", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
		errSimulado := errors.New("fallo simulado del estimador")
		estimador := &fakeEstimadorStoryPoints{err: errSimulado}
		servicio := entorno.servicioConEstimador(estimador)

		historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
		if !errors.Is(err, errSimulado) {
			t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
		}
		if historia != (domain.ItemBacklog{}) {
			t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
		}
	})
	t.Run("error de dominio pasa sin envolver", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
		estimador := &fakeEstimadorStoryPoints{err: domain.ErrItemNoEncontrado}
		servicio := entorno.servicioConEstimador(estimador)

		historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
		if err != domain.ErrItemNoEncontrado {
			t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
		}
		if historia != (domain.ItemBacklog{}) {
			t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
		}
	})
}

// D14 / R7 (US-13): con el constructor original (estimador nil) AcordarEstimacion
// devuelve un error de configuración envuelto; nunca un pánico ni un éxito falso.
func TestUS13_ServicioPlanningPoker_AcordarEstimacionConstructorOriginal_ErrorDeConfiguracion(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicio()
	storyPointsPrevios := historiaStoryPointsPoker(t, entorno, "p1", "HU-01")

	historia, err := servicio.AcordarEstimacion("p1", "HU-01", 1, "Dana Guzmán", 8)
	if err == nil {
		t.Fatal("err = nil, se esperaba un error de configuración")
	}
	if errors.Is(err, domain.ErrRondaNoRevelada) ||
		errors.Is(err, domain.ErrRondaSinVotos) ||
		errors.Is(err, domain.ErrStoryPointsInvalidos) ||
		errors.Is(err, domain.ErrIntegranteNoDelProyecto) ||
		errors.Is(err, domain.ErrNumeroRondaInvalido) ||
		errors.Is(err, domain.ErrRondaNoEncontrada) ||
		errors.Is(err, domain.ErrProyectoNoEncontrado) ||
		errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, el error de configuración no debe mapear a un error de dominio", err)
	}
	if historia != (domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
	if obtenidos := historiaStoryPointsPoker(t, entorno, "p1", "HU-01"); obtenidos != storyPointsPrevios {
		t.Errorf("StoryPoints = %d, se esperaban %d (sin cambios)", obtenidos, storyPointsPrevios)
	}
}

