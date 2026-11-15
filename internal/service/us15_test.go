package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

type fakeRepositorioDefectosUS15 struct {
	defectos        map[string][]domain.Defecto
	crearErr        error
	listarErr       error
	obtenerErr      error
	actualizarErr   error
	crearLlamadas   int
	listarLlamadas  int
	actualizaciones int
}

func nuevoFakeDefectosUS15() *fakeRepositorioDefectosUS15 {
	return &fakeRepositorioDefectosUS15{defectos: map[string][]domain.Defecto{}}
}

func (f *fakeRepositorioDefectosUS15) Crear(proyectoID string, defecto domain.Defecto) (domain.Defecto, error) {
	f.crearLlamadas++
	if f.crearErr != nil {
		return domain.Defecto{}, f.crearErr
	}
	defecto.ID = int64(len(f.defectos[proyectoID]) + 1)
	f.defectos[proyectoID] = append(f.defectos[proyectoID], defecto)
	return defecto, nil
}

func (f *fakeRepositorioDefectosUS15) ListarPorProyecto(proyectoID string) ([]domain.Defecto, error) {
	f.listarLlamadas++
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return append([]domain.Defecto{}, f.defectos[proyectoID]...), nil
}

func (f *fakeRepositorioDefectosUS15) Obtener(proyectoID string, defectoID int64) (domain.Defecto, error) {
	if f.obtenerErr != nil {
		return domain.Defecto{}, f.obtenerErr
	}
	for _, defecto := range f.defectos[proyectoID] {
		if defecto.ID == defectoID {
			return defecto, nil
		}
	}
	return domain.Defecto{}, domain.ErrDefectoNoEncontrado
}

func (f *fakeRepositorioDefectosUS15) Actualizar(proyectoID string, defecto domain.Defecto) error {
	f.actualizaciones++
	if f.actualizarErr != nil {
		return f.actualizarErr
	}
	for i := range f.defectos[proyectoID] {
		if f.defectos[proyectoID][i].ID == defecto.ID {
			f.defectos[proyectoID][i] = defecto
			return nil
		}
	}
	return domain.ErrDefectoNoEncontrado
}

func prepararServicioDefectosUS15(t *testing.T) (
	*fakeRepositorioProyectos,
	*fakeRepositorioBacklog,
	*fakeRepositorioSprints,
	*fakeRepositorioDefectosUS15,
	*ServicioDefectos,
) {
	t.Helper()
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	backlog := nuevoFakeBacklog()
	backlog.items["p1"] = []domain.ItemBacklog{
		{ID: "HU-01", Titulo: "Historia de prueba", Estado: domain.EstadoItemPendiente, StoryPoints: 3},
	}
	sprints := nuevoFakeSprints()
	sprints.sprints["p1"] = []domain.Sprint{
		sprintUS07Servicio(1, domain.EstadoSprintActivo),
		sprintUS07Servicio(2, domain.EstadoSprintCerrado),
	}
	defectos := nuevoFakeDefectosUS15()
	servicio := NuevoServicioDefectos(proyectos, backlog, sprints, defectos)
	return proyectos, backlog, sprints, defectos, servicio
}

func defectoEntradaUS15() domain.Defecto {
	return domain.Defecto{
		Descripcion:     "No guarda los cambios",
		Severidad:       domain.SeveridadDefectoAlta,
		Estado:          domain.EstadoDefectoAbierto,
		ItemID:          "HU-01",
		SprintDeteccion: 1,
	}
}

func TestUS15_RegistrarDefectoPersisteYDevuelveID(t *testing.T) {
	_, _, _, repo, servicio := prepararServicioDefectosUS15(t)
	entrada := defectoEntradaUS15()

	creado, err := servicio.RegistrarDefecto("p1", entrada)
	if err != nil {
		t.Fatalf("RegistrarDefecto() devolvió error: %v", err)
	}
	if creado.ID != 1 || creado.Descripcion != entrada.Descripcion ||
		creado.ItemID != entrada.ItemID || creado.SprintDeteccion != 1 {
		t.Fatalf("defecto creado = %+v", creado)
	}
	if repo.crearLlamadas != 1 || len(repo.defectos["p1"]) != 1 {
		t.Fatalf("se esperaba un registro, llamadas=%d defectos=%d", repo.crearLlamadas, len(repo.defectos["p1"]))
	}
}

func TestUS15_RegistrarDefectoValidaProyectoReferenciasYCampos(t *testing.T) {
	casos := []struct {
		nombre   string
		preparar func(*testing.T, *fakeRepositorioProyectos, *fakeRepositorioBacklog, *fakeRepositorioSprints)
		proyecto string
		defecto  domain.Defecto
		esperado error
	}{
		{
			nombre:   "proyecto inexistente",
			proyecto: "p404",
			defecto:  defectoEntradaUS15(),
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			nombre:   "historia inexistente",
			proyecto: "p1",
			defecto:  func() domain.Defecto { d := defectoEntradaUS15(); d.ItemID = "HU-99"; return d }(),
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			nombre:   "Sprint detección inexistente",
			proyecto: "p1",
			defecto:  func() domain.Defecto { d := defectoEntradaUS15(); d.SprintDeteccion = 99; return d }(),
			esperado: domain.ErrSprintNoEncontrado,
		},
		{
			nombre:   "Sprint resolución inexistente",
			proyecto: "p1",
			defecto:  func() domain.Defecto { d := defectoEntradaUS15(); n := 99; d.SprintResolucion = &n; return d }(),
			esperado: domain.ErrSprintNoEncontrado,
		},
		{
			nombre:   "severidad inválida",
			proyecto: "p1",
			defecto:  func() domain.Defecto { d := defectoEntradaUS15(); d.Severidad = "Crítica"; return d }(),
			esperado: domain.ErrSeveridadDefectoInvalida,
		},
		{
			nombre:   "Sprint de resolución anterior",
			proyecto: "p1",
			defecto: func() domain.Defecto {
				d := defectoEntradaUS15()
				d.SprintDeteccion = 2
				n := 1
				d.SprintResolucion = &n
				return d
			}(),
			esperado: domain.ErrSprintResolucionAnteriorADeteccion,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog, sprints, repo, servicio := prepararServicioDefectosUS15(t)
			if caso.preparar != nil {
				caso.preparar(t, proyectos, backlog, sprints)
			}
			_, err := servicio.RegistrarDefecto(caso.proyecto, caso.defecto)
			if !errors.Is(err, caso.esperado) {
				t.Fatalf("err = %v, se esperaba %v", err, caso.esperado)
			}
			if repo.crearLlamadas != 0 {
				t.Errorf("se persistió el defecto inválido; llamadas=%d", repo.crearLlamadas)
			}
		})
	}
}

func TestUS15_CambiarEstadoDefectoRequiereSprintYActualiza(t *testing.T) {
	_, _, _, repo, servicio := prepararServicioDefectosUS15(t)
	defecto, err := servicio.RegistrarDefecto("p1", defectoEntradaUS15())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := servicio.CambiarEstadoDefecto("p1", defecto.ID, domain.EstadoDefectoResuelto, nil); !errors.Is(err, domain.ErrSprintResolucionDefectoObligatorio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintResolucionDefectoObligatorio)
	}
	if repo.actualizaciones != 0 {
		t.Fatal("no debía actualizarse el defecto al faltar Sprint de resolución")
	}

	resolucion := 2
	actualizado, err := servicio.CambiarEstadoDefecto("p1", defecto.ID, domain.EstadoDefectoResuelto, &resolucion)
	if err != nil {
		t.Fatalf("CambiarEstadoDefecto() devolvió error: %v", err)
	}
	if actualizado.Estado != domain.EstadoDefectoResuelto ||
		actualizado.SprintResolucion == nil ||
		*actualizado.SprintResolucion != 2 {
		t.Fatalf("defecto actualizado = %+v", actualizado)
	}
}

func TestUS15_ReabrirDefectoLimpiaSprintResolucion(t *testing.T) {
	_, _, _, repo, servicio := prepararServicioDefectosUS15(t)
	defecto, err := servicio.RegistrarDefecto("p1", defectoEntradaUS15())
	if err != nil {
		t.Fatal(err)
	}
	resolucion := 2
	defecto, err = servicio.CambiarEstadoDefecto("p1", defecto.ID, domain.EstadoDefectoResuelto, &resolucion)
	if err != nil {
		t.Fatalf("resolver defecto: %v", err)
	}

	reabierto, err := servicio.CambiarEstadoDefecto("p1", defecto.ID, domain.EstadoDefectoEnCurso, nil)
	if err != nil {
		t.Fatalf("reabrir defecto: %v", err)
	}
	if reabierto.Estado != domain.EstadoDefectoEnCurso || reabierto.SprintResolucion != nil {
		t.Fatalf("defecto reabierto = %+v; se esperaba En curso sin Sprint de resolución", reabierto)
	}
	if repo.actualizaciones != 2 {
		t.Fatalf("el repositorio recibió %d actualizaciones; se esperaban 2", repo.actualizaciones)
	}
	if persistido := repo.defectos["p1"][0]; persistido.Estado != domain.EstadoDefectoEnCurso || persistido.SprintResolucion != nil {
		t.Fatalf("defecto persistido tras reapertura = %+v", persistido)
	}
}

func TestUS15_ListarDefectosDevuelveListaVaciaYEnOrden(t *testing.T) {
	_, _, _, repo, servicio := prepararServicioDefectosUS15(t)
	lista, err := servicio.ListarDefectos("p1")
	if err != nil {
		t.Fatalf("ListarDefectos() devolvió error: %v", err)
	}
	if lista == nil || len(lista) != 0 {
		t.Fatalf("lista = %#v; se esperaba lista vacía", lista)
	}
	for _, descripcion := range []string{"Primero", "Segundo"} {
		_, err := servicio.RegistrarDefecto("p1", func() domain.Defecto {
			d := defectoEntradaUS15()
			d.Descripcion = descripcion
			return d
		}())
		if err != nil {
			t.Fatal(err)
		}
	}
	lista, err = servicio.ListarDefectos("p1")
	if err != nil || len(lista) != 2 || lista[0].ID != 1 || lista[1].ID != 2 {
		t.Fatalf("lista = %+v, err=%v; se esperaba orden por ID ascendente", lista, err)
	}
	if repo.listarLlamadas != 2 {
		t.Errorf("ListarPorProyecto se llamó %d veces, se esperaban 2", repo.listarLlamadas)
	}
}
