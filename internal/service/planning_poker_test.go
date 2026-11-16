package service

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// claveHistoriaPoker identifica una historia del backlog dentro de un proyecto
// para las rondas de Planning Poker (clave compuesta, como
// (proyecto_id, item_numero)).
type claveHistoriaPoker struct {
	proyectoID string
	itemID     string
}

// claveRondaPoker identifica una ronda de Planning Poker dentro de una
// historia (clave compuesta, como (proyecto_id, item_numero, numero)).
type claveRondaPoker struct {
	proyectoID string
	itemID     string
	numero     int
}

// fakeRepositorioRondas es el doble en memoria de RepositorioRondas para los
// tests de servicio de US-12 (D1, D9 de ADR-001: solo PostgreSQL; el doble
// sigue el patrón US-07/US-09). No depende de internal/storage. Permite
// inyectar errores por operación y registra las llamadas recibidas.
type fakeRepositorioRondas struct {
	rondas                map[claveHistoriaPoker][]domain.RondaPoker
	votos                 map[claveRondaPoker][]domain.VotoPoker
	iniciarErr            error
	obtenerErr            error
	listarErr             error
	registrarVotoErr      error
	yaVotaronErr          error
	votosErr              error
	revelarErr            error
	iniciarLlamadas       int
	obtenerLlamadas       int
	listarLlamadas        int
	registrarVotoLlamadas int
	yaVotaronLlamadas     int
	votosLlamadas         int
	revelarLlamadas       int
	votosRegistrados      []domain.VotoPoker
	ultimoNumeroIniciado  int
}

func nuevoFakeRondas() *fakeRepositorioRondas {
	return &fakeRepositorioRondas{
		rondas: map[claveHistoriaPoker][]domain.RondaPoker{},
		votos:  map[claveRondaPoker][]domain.VotoPoker{},
	}
}

func (f *fakeRepositorioRondas) Iniciar(proyectoID, itemID string, numero int) error {
	f.iniciarLlamadas++
	f.ultimoNumeroIniciado = numero
	if f.iniciarErr != nil {
		return f.iniciarErr
	}
	clave := claveHistoriaPoker{proyectoID, itemID}
	for _, ronda := range f.rondas[clave] {
		if ronda.Estado == domain.EstadoRondaAbierta {
			return domain.ErrRondaAbiertaExistente
		}
	}
	f.rondas[clave] = append(f.rondas[clave], domain.RondaPoker{Numero: numero, Estado: domain.EstadoRondaAbierta})
	return nil
}

func (f *fakeRepositorioRondas) Obtener(proyectoID, itemID string, numero int) (domain.RondaPoker, error) {
	f.obtenerLlamadas++
	if f.obtenerErr != nil {
		return domain.RondaPoker{}, f.obtenerErr
	}
	for _, ronda := range f.rondas[claveHistoriaPoker{proyectoID, itemID}] {
		if ronda.Numero == numero {
			return ronda, nil
		}
	}
	return domain.RondaPoker{}, domain.ErrRondaNoEncontrada
}

func (f *fakeRepositorioRondas) ListarPorHistoria(proyectoID, itemID string) ([]domain.RondaPoker, error) {
	f.listarLlamadas++
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	rondas := append([]domain.RondaPoker(nil), f.rondas[claveHistoriaPoker{proyectoID, itemID}]...)
	sort.Slice(rondas, func(i, j int) bool { return rondas[i].Numero < rondas[j].Numero })
	return rondas, nil
}

func (f *fakeRepositorioRondas) RegistrarVoto(proyectoID, itemID string, numero int, voto domain.VotoPoker) error {
	f.registrarVotoLlamadas++
	f.votosRegistrados = append(f.votosRegistrados, voto)
	if f.registrarVotoErr != nil {
		return f.registrarVotoErr
	}
	clave := claveHistoriaPoker{proyectoID, itemID}
	rondas := f.rondas[clave]
	indice := -1
	for i := range rondas {
		if rondas[i].Numero == numero {
			indice = i
			break
		}
	}
	if indice == -1 {
		return domain.ErrRondaNoEncontrada
	}
	if rondas[indice].Estado == domain.EstadoRondaRevelada {
		return domain.ErrRondaRevelada
	}
	claveVoto := claveRondaPoker{proyectoID, itemID, numero}
	votos := append([]domain.VotoPoker(nil), f.votos[claveVoto]...)
	for i, existente := range votos {
		if existente.Integrante == voto.Integrante {
			votos[i] = voto
			f.votos[claveVoto] = votos
			return nil
		}
	}
	f.votos[claveVoto] = append(votos, voto)
	return nil
}

func (f *fakeRepositorioRondas) YaVotaron(proyectoID, itemID string, numero int) ([]string, error) {
	f.yaVotaronLlamadas++
	if f.yaVotaronErr != nil {
		return nil, f.yaVotaronErr
	}
	votos := f.votos[claveRondaPoker{proyectoID, itemID, numero}]
	nombres := make([]string, 0, len(votos))
	for _, voto := range votos {
		nombres = append(nombres, voto.Integrante)
	}
	ordenarNombresPoker(nombres)
	return nombres, nil
}

func (f *fakeRepositorioRondas) Votos(proyectoID, itemID string, numero int) ([]domain.VotoPoker, error) {
	f.votosLlamadas++
	if f.votosErr != nil {
		return nil, f.votosErr
	}
	votos := append([]domain.VotoPoker(nil), f.votos[claveRondaPoker{proyectoID, itemID, numero}]...)
	ordenarVotosPoker(votos)
	return votos, nil
}

func (f *fakeRepositorioRondas) Revelar(proyectoID, itemID string, numero int) error {
	f.revelarLlamadas++
	if f.revelarErr != nil {
		return f.revelarErr
	}
	clave := claveHistoriaPoker{proyectoID, itemID}
	rondas := f.rondas[clave]
	for i := range rondas {
		if rondas[i].Numero == numero {
			if rondas[i].Estado == domain.EstadoRondaRevelada {
				return domain.ErrRondaRevelada
			}
			rondas[i].Estado = domain.EstadoRondaRevelada
			f.rondas[clave] = rondas
			return nil
		}
	}
	return domain.ErrRondaNoEncontrada
}

// ordenarNombresPoker ordena los nombres sin distinguir mayúsculas y
// desempatando por el nombre exacto (D20).
func ordenarNombresPoker(nombres []string) {
	sort.Slice(nombres, func(i, j int) bool {
		ni, nj := strings.ToLower(nombres[i]), strings.ToLower(nombres[j])
		if ni == nj {
			return nombres[i] < nombres[j]
		}
		return ni < nj
	})
}

// ordenarVotosPoker ordena los votos por integrante sin distinguir
// mayúsculas, desempatando por el nombre exacto (D20).
func ordenarVotosPoker(votos []domain.VotoPoker) {
	sort.Slice(votos, func(i, j int) bool {
		ni, nj := strings.ToLower(votos[i].Integrante), strings.ToLower(votos[j].Integrante)
		if ni == nj {
			return votos[i].Integrante < votos[j].Integrante
		}
		return ni < nj
	})
}

// entornoPoker agrupa los fakes de proyectos, integrantes, backlog y rondas
// para sembrar datos y verificar el estado sin depender de internal/storage.
type entornoPoker struct {
	*fakeRepositorioProyectos
	*fakeRepositorioIntegrantes
	*fakeRepositorioBacklog
	*fakeRepositorioRondas
}

func nuevoEntornoPoker() *entornoPoker {
	return &entornoPoker{
		fakeRepositorioProyectos:   nuevoFakeProyectos(),
		fakeRepositorioIntegrantes: nuevoFakeIntegrantes(),
		fakeRepositorioBacklog:     nuevoFakeBacklog(),
		fakeRepositorioRondas:      nuevoFakeRondas(),
	}
}

func (e *entornoPoker) servicio() *ServicioPlanningPoker {
	return NuevoServicioPlanningPoker(
		e.fakeRepositorioProyectos,
		e.fakeRepositorioIntegrantes,
		e.fakeRepositorioBacklog,
		e.fakeRepositorioRondas,
	)
}

func sembrarProyectoPoker(t *testing.T, entorno *entornoPoker, ids ...string) {
	t.Helper()
	for _, id := range ids {
		entorno.proyectos[id] = proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	}
}

func sembrarHistoriaPoker(t *testing.T, entorno *entornoPoker, proyectoID string, ids ...string) {
	t.Helper()
	for _, id := range ids {
		entorno.items[proyectoID] = append(entorno.items[proyectoID], domain.ItemBacklog{
			ID:          id,
			Titulo:      "Historia " + id,
			Descripcion: "Descripción " + id,
			Prioridad:   "Alta",
			Estado:      "Pendiente",
			StoryPoints: 5,
		})
	}
}

func sembrarIntegrantesPoker(t *testing.T, entorno *entornoPoker, proyectoID string, nombres ...string) {
	t.Helper()
	for _, nombre := range nombres {
		entorno.integrantes[proyectoID] = append(entorno.integrantes[proyectoID], domain.Integrante{Nombre: nombre, Rol: "Desarrollador"})
	}
}

func sembrarRondaPoker(t *testing.T, entorno *entornoPoker, proyectoID, itemID string, ronda domain.RondaPoker) {
	t.Helper()
	clave := claveHistoriaPoker{proyectoID, itemID}
	entorno.rondas[clave] = append(entorno.rondas[clave], ronda)
}

func sembrarVotoPoker(t *testing.T, entorno *entornoPoker, proyectoID, itemID string, numero int, votos ...domain.VotoPoker) {
	t.Helper()
	clave := claveRondaPoker{proyectoID, itemID, numero}
	entorno.votos[clave] = append(entorno.votos[clave], votos...)
}

func nuevoEntornoPokerSembrado(t *testing.T) *entornoPoker {
	t.Helper()
	entorno := nuevoEntornoPoker()
	sembrarProyectoPoker(t, entorno, "p1")
	sembrarHistoriaPoker(t, entorno, "p1", "HU-01")
	sembrarIntegrantesPoker(t, entorno, "p1", "Dana Guzmán", "Victoria Caro")
	return entorno
}

// verificarConsultaPokerVacia comprueba que una consulta fallida devuelve la
// consulta en cero (R17: ConsultaRondaPoker no es comparable por contener
// slices; se comparan los campos).
func verificarConsultaPokerVacia(t *testing.T, consulta domain.ConsultaRondaPoker) {
	t.Helper()
	if consulta.Ronda != (domain.RondaPoker{}) {
		t.Errorf("Ronda = %+v, se esperaba domain.RondaPoker{}", consulta.Ronda)
	}
	if consulta.YaVotaron != nil {
		t.Errorf("YaVotaron = %+v, se esperaba nil", consulta.YaVotaron)
	}
	if consulta.Votos != nil {
		t.Errorf("Votos = %+v, se esperaba nil", consulta.Votos)
	}
}

func votosPersistidosPoker(t *testing.T, entorno *entornoPoker, proyectoID, itemID string, numero int) []domain.VotoPoker {
	t.Helper()
	return append([]domain.VotoPoker(nil), entorno.votos[claveRondaPoker{proyectoID, itemID, numero}]...)
}

// R9 / D7 / CA1 (US-12): la primera ronda de una historia es la número 1 y se
// persiste en estado Abierta.
func TestUS12_Servicio_IniciarRonda_PrimeraRondaEsUno(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("p1", "HU-01")
	if err != nil {
		t.Fatalf("IniciarRonda() devolvió un error inesperado: %v", err)
	}
	if ronda.Numero != 1 {
		t.Errorf("Numero = %d, se esperaba %d", ronda.Numero, 1)
	}
	if ronda.Estado != domain.EstadoRondaAbierta {
		t.Errorf("Estado = %q, se esperaba %q", ronda.Estado, domain.EstadoRondaAbierta)
	}
	if entorno.iniciarLlamadas != 1 {
		t.Errorf("Iniciar del repositorio fue invocado %d veces, se esperaba 1", entorno.iniciarLlamadas)
	}
	if entorno.ultimoNumeroIniciado != 1 {
		t.Errorf("Iniciar fue invocado con numero = %d, se esperaba %d", entorno.ultimoNumeroIniciado, 1)
	}
	rondas := entorno.rondas[claveHistoriaPoker{"p1", "HU-01"}]
	if len(rondas) != 1 {
		t.Fatalf("se esperaba 1 ronda persistida, hay %d", len(rondas))
	}
	if rondas[0] != (domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}) {
		t.Errorf("ronda persistida = %+v, se esperaba {Numero:1 Estado:Abierta}", rondas[0])
	}
}

// R9 / R16 / D7 (US-12): con la ronda anterior Revelada, la nueva lleva
// max(numero)+1.
func TestUS12_Servicio_IniciarRonda_SiguienteNumeroTrasRevelada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 3, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("p1", "HU-01")
	if err != nil {
		t.Fatalf("IniciarRonda() devolvió un error inesperado: %v", err)
	}
	if ronda.Numero != 4 {
		t.Errorf("Numero = %d, se esperaba %d (max+1)", ronda.Numero, 4)
	}
	if ronda.Estado != domain.EstadoRondaAbierta {
		t.Errorf("Estado = %q, se esperaba %q", ronda.Estado, domain.EstadoRondaAbierta)
	}
	if entorno.ultimoNumeroIniciado != 4 {
		t.Errorf("Iniciar fue invocado con numero = %d, se esperaba %d", entorno.ultimoNumeroIniciado, 4)
	}
}

// R9 / R16 / V7 / E8 (US-12): con una ronda Abierta pendiente no se crea la
// siguiente y se devuelve domain.ErrRondaAbiertaExistente.
func TestUS12_Servicio_IniciarRonda_ConAbiertaDevuelveErrRondaAbiertaExistente(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("p1", "HU-01")
	if !errors.Is(err, domain.ErrRondaAbiertaExistente) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaAbiertaExistente)
	}
	if err != domain.ErrRondaAbiertaExistente {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaAbiertaExistente)
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
	if entorno.iniciarLlamadas != 0 {
		t.Errorf("Iniciar fue invocado %d veces; no debía crearse ninguna ronda", entorno.iniciarLlamadas)
	}
	if len(entorno.rondas[claveHistoriaPoker{"p1", "HU-01"}]) != 2 {
		t.Errorf("se esperaban 2 rondas persistidas, hay %d", len(entorno.rondas[claveHistoriaPoker{"p1", "HU-01"}]))
	}
}

// V1 / R1 (US-12): proyecto inexistente prevalece sobre cualquier otro error;
// no se consulta el backlog ni las rondas.
func TestUS12_Servicio_IniciarRonda_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("inexistente", "HU-01")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 || entorno.fakeRepositorioRondas.listarLlamadas != 0 || entorno.iniciarLlamadas != 0 {
		t.Errorf("se llamó a backlog/rondas (obtener=%d, listar=%d, iniciar=%d); el proyecto debía validarse primero",
			entorno.fakeRepositorioBacklog.obtenerLlamadas, entorno.fakeRepositorioRondas.listarLlamadas, entorno.iniciarLlamadas)
	}
}

// V2 / R1 (US-12): historia inexistente devuelve domain.ErrItemNoEncontrado;
// no se persiste nada.
func TestUS12_Servicio_IniciarRonda_HistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("p1", "HU-99")
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if err != domain.ErrItemNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
	if entorno.fakeRepositorioRondas.listarLlamadas != 0 || entorno.iniciarLlamadas != 0 {
		t.Errorf("se llamó a rondas (listar=%d, iniciar=%d); la historia debía validarse primero",
			entorno.fakeRepositorioRondas.listarLlamadas, entorno.iniciarLlamadas)
	}
}

// E11 / R15 / V10 (US-12): los errores de repositorio no conocidos se
// envuelven con %w; los errores de dominio conocidos pasan sin envolver.
func TestUS12_Servicio_IniciarRonda_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar las rondas", func(e *entornoPoker) { e.fakeRepositorioRondas.listarErr = errSimulado }},
		{"error al iniciar la ronda", func(e *entornoPoker) { e.iniciarErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			caso.preparar(entorno)
			servicio := entorno.servicio()

			ronda, err := servicio.IniciarRonda("p1", "HU-01")
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrItemNoEncontrado) ||
				errors.Is(err, domain.ErrRondaAbiertaExistente) {
				t.Errorf("err = %v, no debería mapearse a un error de dominio", err)
			}
			if ronda != (domain.RondaPoker{}) {
				t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
			}
		})
	}
}

// R10 / E10 (US-12): si el repositorio informa una violación de la unicidad
// de ronda Abierta (carrera), el error pasa sin envolver.
func TestUS12_Servicio_IniciarRonda_ErrRondaAbiertaExistenteDelRepositorioPasaSinEnvolver(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	entorno.iniciarErr = domain.ErrRondaAbiertaExistente
	servicio := entorno.servicio()

	ronda, err := servicio.IniciarRonda("p1", "HU-01")
	if err != domain.ErrRondaAbiertaExistente {
		t.Fatalf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaAbiertaExistente)
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
}

// CA1 / R1 (US-12): un voto válido se registra en la ronda Abierta.
func TestUS12_Servicio_Votar_VotoValido(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	if err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 5); err != nil {
		t.Fatalf("Votar() devolvió un error inesperado: %v", err)
	}
	if entorno.registrarVotoLlamadas != 1 {
		t.Errorf("RegistrarVoto fue invocado %d veces, se esperaba 1", entorno.registrarVotoLlamadas)
	}
	if len(entorno.votosRegistrados) != 1 {
		t.Fatalf("se esperaba 1 voto registrado por el doble, hay %d", len(entorno.votosRegistrados))
	}
	if entorno.votosRegistrados[0] != (domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5}) {
		t.Errorf("voto registrado = %+v, se esperaba {Integrante:\"Dana Guzmán\" Valor:5}", entorno.votosRegistrados[0])
	}
	votos := votosPersistidosPoker(t, entorno, "p1", "HU-01", 1)
	if len(votos) != 1 {
		t.Fatalf("se esperaba 1 voto persistido, hay %d", len(votos))
	}
	if votos[0] != (domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5}) {
		t.Errorf("voto persistido = %+v, se esperaba {Integrante:\"Dana Guzmán\" Valor:5}", votos[0])
	}
}

// CA3 / R3 / V8 (US-12): un segundo voto del mismo integrante en la ronda
// Abierta reemplaza el anterior (upsert).
func TestUS12_Servicio_Votar_CambioDeVotoReemplaza(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	if err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 5); err != nil {
		t.Fatalf("primer Votar() devolvió un error inesperado: %v", err)
	}
	if err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 13); err != nil {
		t.Fatalf("segundo Votar() devolvió un error inesperado: %v", err)
	}
	votos := votosPersistidosPoker(t, entorno, "p1", "HU-01", 1)
	if len(votos) != 1 {
		t.Fatalf("se esperaba 1 voto persistido (upsert), hay %d", len(votos))
	}
	if votos[0].Valor != 13 {
		t.Errorf("Valor = %d, se esperaba %d (el último escrito gana)", votos[0].Valor, 13)
	}
}

// CA3 / R4 / V5 / E6 (US-12): votar con la ronda Revelada devuelve
// domain.ErrRondaRevelada y el voto previo permanece intacto.
func TestUS12_Servicio_Votar_RondaRevelada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1, domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5})
	servicio := entorno.servicio()

	err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 8)
	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaRevelada)
	}
	if err != domain.ErrRondaRevelada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaRevelada)
	}
	if entorno.registrarVotoLlamadas != 0 {
		t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
	}
	votos := votosPersistidosPoker(t, entorno, "p1", "HU-01", 1)
	if len(votos) != 1 || votos[0].Valor != 5 {
		t.Errorf("voto previo = %+v, se esperaba intacto {Integrante:\"Dana Guzmán\" Valor:5}", votos)
	}
}

// R4 / E6 (US-12): si el repositorio informa la ronda Revelada dentro de la
// transacción (carrera con Revelar), el error pasa sin envolver.
func TestUS12_Servicio_Votar_ErrRondaReveladaDelRepositorioPasaSinEnvolver(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	entorno.registrarVotoErr = domain.ErrRondaRevelada
	servicio := entorno.servicio()

	err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 8)
	if err != domain.ErrRondaRevelada {
		t.Fatalf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaRevelada)
	}
}

// R5 / V6 / E7 (US-12): un valor fuera de la serie devuelve
// domain.ErrValorVotoPokerInvalido; no se persiste nada.
func TestUS12_Servicio_Votar_ValorInvalido(t *testing.T) {
	for _, valor := range []int{0, 4, 6, 22, -1} {
		t.Run(strconv.Itoa(valor), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			servicio := entorno.servicio()

			err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", valor)
			if !errors.Is(err, domain.ErrValorVotoPokerInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrValorVotoPokerInvalido)
			}
			if entorno.registrarVotoLlamadas != 0 {
				t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
			}
		})
	}
}

// R13 / D17 / CA4 (US-12): se persiste el nombre canónico almacenado del
// integrante, no el texto tal como llega.
func TestUS12_Servicio_Votar_NombreCanonicoPersistido(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	for _, entrada := range []string{"dana guzmán", "DANA GUZMÁN", "  Dana Guzmán  "} {
		if err := servicio.Votar("p1", "HU-01", 1, entrada, 5); err != nil {
			t.Fatalf("Votar(%q) devolvió un error inesperado: %v", entrada, err)
		}
	}
	votos := votosPersistidosPoker(t, entorno, "p1", "HU-01", 1)
	if len(votos) != 1 {
		t.Fatalf("se esperaba 1 voto persistido (misma identidad), hay %d", len(votos))
	}
	if votos[0].Integrante != "Dana Guzmán" {
		t.Errorf("Integrante persistido = %q, se esperaba el nombre canónico %q", votos[0].Integrante, "Dana Guzmán")
	}
}

// CA4 / V3 / E3 (US-12): un integrante que no pertenece al proyecto no puede
// votar; no se persiste nada.
func TestUS12_Servicio_Votar_IntegranteNoDelProyecto(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	err := servicio.Votar("p1", "HU-01", 1, "Intruso", 5)
	if !errors.Is(err, domain.ErrIntegranteNoDelProyecto) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrIntegranteNoDelProyecto)
	}
	if err != domain.ErrIntegranteNoDelProyecto {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrIntegranteNoDelProyecto)
	}
	if entorno.registrarVotoLlamadas != 0 {
		t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
	}
}

// E12 / V3 / D15 (US-12): nombre vacío o compuesto de espacios se rechaza al
// inicio del paso 3, antes de ListarPorProyecto; no se consulta la base.
func TestUS12_Servicio_Votar_NombreVacioOEnBlanco(t *testing.T) {
	for _, nombre := range []string{"", "   ", "\t\n"} {
		t.Run("nombre "+strconv.Itoa(len(nombre)), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			servicio := entorno.servicio()

			err := servicio.Votar("p1", "HU-01", 1, nombre, 5)
			if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
			}
			if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 {
				t.Errorf("ListarPorProyecto fue invocado %d veces; el nombre debía validarse antes", entorno.fakeRepositorioIntegrantes.listarLlamadas)
			}
			if entorno.registrarVotoLlamadas != 0 {
				t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
			}
		})
	}
}

// D14 / R14 / V4 / E4 (US-12): rondaNumero <= 0 devuelve
// domain.ErrNumeroRondaInvalido en el paso 4, antes de rondas.Obtener.
func TestUS12_Servicio_Votar_NumeroRondaInvalido(t *testing.T) {
	for _, numero := range []int{0, -1} {
		t.Run("numero "+strconv.Itoa(numero), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			servicio := entorno.servicio()

			err := servicio.Votar("p1", "HU-01", numero, "Dana Guzmán", 5)
			if !errors.Is(err, domain.ErrNumeroRondaInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroRondaInvalido)
			}
			if err != domain.ErrNumeroRondaInvalido {
				t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrNumeroRondaInvalido)
			}
			if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
				t.Errorf("Obtener fue invocado %d veces; nunca se debía llegar a rondas.Obtener con numero <= 0", entorno.fakeRepositorioRondas.obtenerLlamadas)
			}
			if entorno.registrarVotoLlamadas != 0 {
				t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
			}
		})
	}
}

// V4 / E5 (US-12): ronda inexistente devuelve domain.ErrRondaNoEncontrada; no
// se persiste nada.
func TestUS12_Servicio_Votar_RondaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	err := servicio.Votar("p1", "HU-01", 7, "Dana Guzmán", 5)
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err != domain.ErrRondaNoEncontrada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoEncontrada)
	}
	if entorno.registrarVotoLlamadas != 0 {
		t.Errorf("RegistrarVoto fue invocado %d veces; no debía escribirse nada", entorno.registrarVotoLlamadas)
	}
}

// V1 (US-12): proyecto inexistente prevalece sobre historia, integrante,
// ronda y valor.
func TestUS12_Servicio_Votar_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	err := servicio.Votar("inexistente", "HU-01", 1, "Dana Guzmán", 5)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 || entorno.fakeRepositorioRondas.obtenerLlamadas != 0 || entorno.registrarVotoLlamadas != 0 {
		t.Errorf("se llamó a integrantes/rondas; el proyecto debía validarse primero")
	}
}

// V2 / R1 (US-12): historia inexistente prevalece sobre integrante, ronda y
// valor.
func TestUS12_Servicio_Votar_HistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	err := servicio.Votar("p1", "HU-99", 1, "Dana Guzmán", 5)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if err != domain.ErrItemNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
	}
	if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 || entorno.fakeRepositorioRondas.obtenerLlamadas != 0 || entorno.registrarVotoLlamadas != 0 {
		t.Errorf("se llamó a integrantes/rondas; la historia debía validarse primero")
	}
}

// E11 / R15 / V10 (US-12): los errores de repositorio no conocidos se
// envuelven con %w.
func TestUS12_Servicio_Votar_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar los integrantes", func(e *entornoPoker) { e.fakeRepositorioIntegrantes.listarErr = errSimulado }},
		{"error al obtener la ronda", func(e *entornoPoker) {
			e.fakeRepositorioRondas.obtenerErr = errSimulado
		}},
		{"error al registrar el voto", func(e *entornoPoker) { e.registrarVotoErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			caso.preparar(entorno)
			servicio := entorno.servicio()

			err := servicio.Votar("p1", "HU-01", 1, "Dana Guzmán", 5)
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
		})
	}
}

// R1 / D13 (US-12): orden de evaluación proyecto → historia → integrante →
// ronda → valor; el error de más arriba prevalece sobre el de más abajo.
func TestUS12_Servicio_Votar_Precedencia(t *testing.T) {
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
			proyectoID:  "inexistente", itemID: "HU-99", rondaNumero: 1, nombre: "Dana Guzmán", valor: 5,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre integrante ajeno",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 5,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "proyecto inexistente prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "inexistente", itemID: "HU-01", rondaNumero: 7, nombre: "Dana Guzmán", valor: 5,
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
			proyectoID:  "p1", itemID: "HU-99", rondaNumero: 1, nombre: "Intruso", valor: 5,
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion: "historia inexistente prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-99", rondaNumero: 7, nombre: "Dana Guzmán", valor: 5,
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
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 7, nombre: "Intruso", valor: 5,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "integrante ajeno prevalece sobre ronda Revelada",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 5,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "integrante ajeno prevalece sobre valor inválido",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Intruso", valor: 4,
			esperado: domain.ErrIntegranteNoDelProyecto,
		},
		{
			descripcion: "nombre en blanco prevalece sobre rondaNumero <= 0",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 0, nombre: "   ", valor: 5,
			esperado: domain.ErrNombreIntegranteVacio,
		},
		{
			descripcion: "nombre en blanco prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 7, nombre: "", valor: 5,
			esperado: domain.ErrNombreIntegranteVacio,
		},
		{
			descripcion: "rondaNumero <= 0 prevalece sobre ronda inexistente",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 0, nombre: "Dana Guzmán", valor: 5,
			esperado: domain.ErrNumeroRondaInvalido,
		},
		{
			descripcion: "ronda inexistente prevalece sobre valor inválido",
			preparar:    func(t *testing.T, e *entornoPoker) {},
			proyectoID:  "p1", itemID: "HU-01", rondaNumero: 7, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrRondaNoEncontrada,
		},
		{
			descripcion: "ronda Revelada prevalece sobre valor inválido",
			preparar: func(t *testing.T, e *entornoPoker) {
				sembrarRondaPoker(t, e, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
			},
			proyectoID: "p1", itemID: "HU-01", rondaNumero: 1, nombre: "Dana Guzmán", valor: 4,
			esperado: domain.ErrRondaRevelada,
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			caso.preparar(t, entorno)
			servicio := entorno.servicio()

			err := servicio.Votar(caso.proyectoID, caso.itemID, caso.rondaNumero, caso.nombre, caso.valor)
			if !errors.Is(err, caso.esperado) {
				t.Fatalf("err = %v, se esperaba %v", err, caso.esperado)
			}
		})
	}
}

// D6 / R7 / CA2 (US-12): cualquier integrante del proyecto puede revelar la
// ronda, incluso uno que no votó.
func TestUS12_Servicio_RevelarRonda_CualquierIntegranteDelProyecto(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1, domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5})
	servicio := entorno.servicio()

	if err := servicio.RevelarRonda("p1", "HU-01", 1, "Victoria Caro"); err != nil {
		t.Fatalf("RevelarRonda() devolvió un error inesperado: %v", err)
	}
	if entorno.revelarLlamadas != 1 {
		t.Errorf("Revelar fue invocado %d veces, se esperaba 1", entorno.revelarLlamadas)
	}
	rondas := entorno.rondas[claveHistoriaPoker{"p1", "HU-01"}]
	if len(rondas) != 1 || rondas[0].Estado != domain.EstadoRondaRevelada {
		t.Errorf("ronda = %+v, se esperaba estado %q", rondas, domain.EstadoRondaRevelada)
	}
}

// D9 / E13 / R7 (US-12): revelar una ronda con cero votos está permitido.
func TestUS12_Servicio_RevelarRonda_SinVotosPermitido(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	if err := servicio.RevelarRonda("p1", "HU-01", 1, "Dana Guzmán"); err != nil {
		t.Fatalf("RevelarRonda() devolvió un error inesperado: %v", err)
	}
	rondas := entorno.rondas[claveHistoriaPoker{"p1", "HU-01"}]
	if len(rondas) != 1 || rondas[0].Estado != domain.EstadoRondaRevelada {
		t.Errorf("ronda = %+v, se esperaba estado %q", rondas, domain.EstadoRondaRevelada)
	}
}

// D6 / R7 / V9 / E9 (US-12): revelar una ronda ya Revelada devuelve
// domain.ErrRondaRevelada; irreversible.
func TestUS12_Servicio_RevelarRonda_YaRevelada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicio()

	err := servicio.RevelarRonda("p1", "HU-01", 1, "Dana Guzmán")
	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaRevelada)
	}
	if err != domain.ErrRondaRevelada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaRevelada)
	}
	if entorno.revelarLlamadas != 0 {
		t.Errorf("Revelar fue invocado %d veces; no debía cambiarse nada", entorno.revelarLlamadas)
	}
}

// CA4 / V3 / E3 (US-12): un integrante que no pertenece al proyecto no puede
// revelar.
func TestUS12_Servicio_RevelarRonda_IntegranteNoDelProyecto(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	err := servicio.RevelarRonda("p1", "HU-01", 1, "Intruso")
	if !errors.Is(err, domain.ErrIntegranteNoDelProyecto) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrIntegranteNoDelProyecto)
	}
	if err != domain.ErrIntegranteNoDelProyecto {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrIntegranteNoDelProyecto)
	}
	if entorno.revelarLlamadas != 0 {
		t.Errorf("Revelar fue invocado %d veces; no debía cambiarse nada", entorno.revelarLlamadas)
	}
}

// E12 / V3 / D15 (US-12): nombre vacío o compuesto de espacios se rechaza al
// inicio del paso 3, antes de ListarPorProyecto.
func TestUS12_Servicio_RevelarRonda_NombreVacioOEnBlanco(t *testing.T) {
	for _, nombre := range []string{"", "   "} {
		t.Run("nombre "+strconv.Itoa(len(nombre)), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			servicio := entorno.servicio()

			err := servicio.RevelarRonda("p1", "HU-01", 1, nombre)
			if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
			}
			if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 {
				t.Errorf("ListarPorProyecto fue invocado %d veces; el nombre debía validarse antes", entorno.fakeRepositorioIntegrantes.listarLlamadas)
			}
			if entorno.revelarLlamadas != 0 {
				t.Errorf("Revelar fue invocado %d veces; no debía cambiarse nada", entorno.revelarLlamadas)
			}
		})
	}
}

// D14 / R14 / E4 (US-12): rondaNumero <= 0 devuelve
// domain.ErrNumeroRondaInvalido en el paso 4, antes de rondas.Obtener.
func TestUS12_Servicio_RevelarRonda_NumeroRondaInvalido(t *testing.T) {
	for _, numero := range []int{0, -2} {
		t.Run("numero "+strconv.Itoa(numero), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			servicio := entorno.servicio()

			err := servicio.RevelarRonda("p1", "HU-01", numero, "Dana Guzmán")
			if !errors.Is(err, domain.ErrNumeroRondaInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroRondaInvalido)
			}
			if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
				t.Errorf("Obtener fue invocado %d veces; nunca se debía llegar a rondas.Obtener con numero <= 0", entorno.fakeRepositorioRondas.obtenerLlamadas)
			}
			if entorno.revelarLlamadas != 0 {
				t.Errorf("Revelar fue invocado %d veces; no debía cambiarse nada", entorno.revelarLlamadas)
			}
		})
	}
}

// V4 / E5 (US-12): ronda inexistente devuelve domain.ErrRondaNoEncontrada.
func TestUS12_Servicio_RevelarRonda_RondaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	err := servicio.RevelarRonda("p1", "HU-01", 7, "Dana Guzmán")
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err != domain.ErrRondaNoEncontrada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoEncontrada)
	}
	if entorno.revelarLlamadas != 0 {
		t.Errorf("Revelar fue invocado %d veces; no debía cambiarse nada", entorno.revelarLlamadas)
	}
}

// V1 / V2 (US-12): proyecto y historia inexistentes prevalecen sobre integrante
// y ronda.
func TestUS12_Servicio_RevelarRonda_ProyectoEHistoriaNoEncontrados(t *testing.T) {
	t.Run("proyecto inexistente", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		servicio := entorno.servicio()

		err := servicio.RevelarRonda("inexistente", "HU-01", 1, "Dana Guzmán")
		if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
		}
		if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 || entorno.fakeRepositorioRondas.obtenerLlamadas != 0 || entorno.revelarLlamadas != 0 {
			t.Errorf("se llamó a integrantes/rondas; el proyecto debía validarse primero")
		}
	})
	t.Run("historia inexistente", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		servicio := entorno.servicio()

		err := servicio.RevelarRonda("p1", "HU-99", 1, "Dana Guzmán")
		if !errors.Is(err, domain.ErrItemNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
		}
		if entorno.fakeRepositorioIntegrantes.listarLlamadas != 0 || entorno.fakeRepositorioRondas.obtenerLlamadas != 0 || entorno.revelarLlamadas != 0 {
			t.Errorf("se llamó a integrantes/rondas; la historia debía validarse primero")
		}
	})
}

// E11 / R15 / V10 (US-12): los errores de repositorio no conocidos se
// envuelven con %w.
func TestUS12_Servicio_RevelarRonda_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar los integrantes", func(e *entornoPoker) { e.fakeRepositorioIntegrantes.listarErr = errSimulado }},
		{"error al obtener la ronda", func(e *entornoPoker) {
			e.fakeRepositorioRondas.obtenerErr = errSimulado
		}},
		{"error al revelar la ronda", func(e *entornoPoker) { e.revelarErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			caso.preparar(entorno)
			servicio := entorno.servicio()

			err := servicio.RevelarRonda("p1", "HU-01", 1, "Dana Guzmán")
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
		})
	}
}

// R6 / R17 / CA2 (US-12): con la ronda Abierta, ConsultarRonda devuelve
// YaVotaron (nombres) y Votos nil; nunca llama a Votos.
func TestUS12_Servicio_ConsultarRonda_AbiertaDevuelveYaVotaronSinValores(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
	)
	servicio := entorno.servicio()

	consulta, err := servicio.ConsultarRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("ConsultarRonda() devolvió un error inesperado: %v", err)
	}
	if consulta.Ronda != (domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}) {
		t.Errorf("Ronda = %+v, se esperaba {Numero:1 Estado:Abierta}", consulta.Ronda)
	}
	if len(consulta.YaVotaron) != 2 {
		t.Fatalf("len(YaVotaron) = %d, se esperaban 2 nombres", len(consulta.YaVotaron))
	}
	if consulta.YaVotaron[0] != "Dana Guzmán" || consulta.YaVotaron[1] != "Victoria Caro" {
		t.Errorf("YaVotaron = %v, se esperaban los nombres ordenados [Dana Guzmán Victoria Caro]", consulta.YaVotaron)
	}
	if consulta.Votos != nil {
		t.Errorf("Votos = %+v, se esperaba nil con la ronda Abierta (R17)", consulta.Votos)
	}
	if entorno.votosLlamadas != 0 {
		t.Errorf("Votos fue invocado %d veces; con la ronda Abierta NUNCA debe llamarse (R6)", entorno.votosLlamadas)
	}
	if entorno.yaVotaronLlamadas != 1 {
		t.Errorf("YaVotaron fue invocado %d veces, se esperaba 1", entorno.yaVotaronLlamadas)
	}
}

// R6 / R17 / CA2 (US-12): con la ronda Revelada, ConsultarRonda devuelve Votos
// y YaVotaron nil; nunca llama a YaVotaron.
func TestUS12_Servicio_ConsultarRonda_ReveladaDevuelveVotosSinYaVotaron(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarVotoPoker(t, entorno, "p1", "HU-01", 1,
		domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8},
		domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5},
	)
	servicio := entorno.servicio()

	consulta, err := servicio.ConsultarRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("ConsultarRonda() devolvió un error inesperado: %v", err)
	}
	if consulta.Ronda != (domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada}) {
		t.Errorf("Ronda = %+v, se esperaba {Numero:1 Estado:Revelada}", consulta.Ronda)
	}
	if len(consulta.Votos) != 2 {
		t.Fatalf("len(Votos) = %d, se esperaban 2 votos", len(consulta.Votos))
	}
	if consulta.Votos[0] != (domain.VotoPoker{Integrante: "Dana Guzmán", Valor: 5}) ||
		consulta.Votos[1] != (domain.VotoPoker{Integrante: "Victoria Caro", Valor: 8}) {
		t.Errorf("Votos = %+v, se esperaban [{Dana Guzmán 5} {Victoria Caro 8}]", consulta.Votos)
	}
	if consulta.YaVotaron != nil {
		t.Errorf("YaVotaron = %+v, se esperaba nil con la ronda Revelada (R17)", consulta.YaVotaron)
	}
	if entorno.yaVotaronLlamadas != 0 {
		t.Errorf("YaVotaron fue invocado %d veces; con la ronda Revelada NUNCA debe llamarse (R6)", entorno.yaVotaronLlamadas)
	}
	if entorno.votosLlamadas != 1 {
		t.Errorf("Votos fue invocado %d veces, se esperaba 1", entorno.votosLlamadas)
	}
}

// V4 / E5 (US-12): consultar una ronda inexistente devuelve
// domain.ErrRondaNoEncontrada.
func TestUS12_Servicio_ConsultarRonda_RondaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	consulta, err := servicio.ConsultarRonda("p1", "HU-01", 7)
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err != domain.ErrRondaNoEncontrada {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaNoEncontrada)
	}
	verificarConsultaPokerVacia(t, consulta)
	if entorno.yaVotaronLlamadas != 0 || entorno.votosLlamadas != 0 {
		t.Errorf("se llamó a YaVotaron (%d) o Votos (%d); la ronda debía validarse primero",
			entorno.yaVotaronLlamadas, entorno.votosLlamadas)
	}
}

// D14 / R14 / E4 (US-12): rondaNumero <= 0 devuelve
// domain.ErrNumeroRondaInvalido en el paso 4, antes de rondas.Obtener.
func TestUS12_Servicio_ConsultarRonda_NumeroRondaInvalido(t *testing.T) {
	for _, numero := range []int{0, -3} {
		t.Run("numero "+strconv.Itoa(numero), func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			servicio := entorno.servicio()

			consulta, err := servicio.ConsultarRonda("p1", "HU-01", numero)
			if !errors.Is(err, domain.ErrNumeroRondaInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroRondaInvalido)
			}
			if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
				t.Errorf("Obtener fue invocado %d veces; nunca se debía llegar a rondas.Obtener con numero <= 0", entorno.fakeRepositorioRondas.obtenerLlamadas)
			}
			verificarConsultaPokerVacia(t, consulta)
		})
	}
}

// V1 / V2 / R1 (US-12): proyecto y historia inexistentes prevalecen sobre la
// ronda.
func TestUS12_Servicio_ConsultarRonda_ProyectoEHistoriaNoEncontrados(t *testing.T) {
	t.Run("proyecto inexistente", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		servicio := entorno.servicio()

		consulta, err := servicio.ConsultarRonda("inexistente", "HU-01", 1)
		if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
		}
		if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
			t.Errorf("Obtener fue invocado %d veces; el proyecto debía validarse primero", entorno.fakeRepositorioRondas.obtenerLlamadas)
		}
		verificarConsultaPokerVacia(t, consulta)
	})
	t.Run("historia inexistente", func(t *testing.T) {
		entorno := nuevoEntornoPokerSembrado(t)
		servicio := entorno.servicio()

		consulta, err := servicio.ConsultarRonda("p1", "HU-99", 1)
		if !errors.Is(err, domain.ErrItemNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
		}
		if entorno.fakeRepositorioRondas.obtenerLlamadas != 0 {
			t.Errorf("Obtener fue invocado %d veces; la historia debía validarse primero", entorno.fakeRepositorioRondas.obtenerLlamadas)
		}
		verificarConsultaPokerVacia(t, consulta)
	})
}

// D9 / R17 (US-12): una ronda Revelada con cero votos devuelve Votos vacío (no
// nil necesariamente, pero sin elementos) y YaVotaron nil.
func TestUS12_Servicio_ConsultarRonda_ReveladaSinVotos(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	servicio := entorno.servicio()

	consulta, err := servicio.ConsultarRonda("p1", "HU-01", 1)
	if err != nil {
		t.Fatalf("ConsultarRonda() devolvió un error inesperado: %v", err)
	}
	if len(consulta.Votos) != 0 {
		t.Errorf("Votos = %+v, se esperaban 0 votos", consulta.Votos)
	}
	if consulta.YaVotaron != nil {
		t.Errorf("YaVotaron = %+v, se esperaba nil con la ronda Revelada (R17)", consulta.YaVotaron)
	}
}

// E11 / R15 / V10 (US-12): los errores de repositorio no conocidos se
// envuelven con %w.
func TestUS12_Servicio_ConsultarRonda_ErrorDeRepositorioEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		estado      string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", domain.EstadoRondaAbierta, func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", domain.EstadoRondaAbierta, func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al obtener la ronda", domain.EstadoRondaAbierta, func(e *entornoPoker) {
			e.fakeRepositorioRondas.obtenerErr = errSimulado
		}},
		{"error al listar quiénes votaron", domain.EstadoRondaAbierta, func(e *entornoPoker) { e.yaVotaronErr = errSimulado }},
		{"error al listar los votos", domain.EstadoRondaRevelada, func(e *entornoPoker) { e.votosErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: caso.estado})
			caso.preparar(entorno)
			servicio := entorno.servicio()

			consulta, err := servicio.ConsultarRonda("p1", "HU-01", 1)
			if !errors.Is(err, errSimulado) {
				t.Fatalf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			verificarConsultaPokerVacia(t, consulta)
		})
	}
}
