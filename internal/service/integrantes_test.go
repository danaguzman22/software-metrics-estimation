package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func nuevoServicioIntegrantesConDatos(registros ...ProyectoRegistrado) (*ServicioProyectos, *fakeRepositorioProyectos, *fakeRepositorioIntegrantes) {
	return nuevoServicioProyectosConDatos(registros...)
}

func buscarIntegranteEnFake(t *testing.T, integrantes *fakeRepositorioIntegrantes, proyectoID, nombre string) domain.Integrante {
	t.Helper()
	for _, it := range integrantes.integrantes[proyectoID] {
		if it.Nombre == nombre {
			return it
		}
	}
	t.Fatalf("no se encontró el integrante %q del proyecto %q", nombre, proyectoID)
	return domain.Integrante{}
}

func verificarIntegrante(t *testing.T, obtenido, esperado domain.Integrante) {
	t.Helper()
	if obtenido.Nombre != esperado.Nombre {
		t.Errorf("Nombre = %q, se esperaba %q", obtenido.Nombre, esperado.Nombre)
	}
	if obtenido.Rol != esperado.Rol {
		t.Errorf("Rol = %q, se esperaba %q", obtenido.Rol, esperado.Rol)
	}
}

func verificarIntegranteVacio(t *testing.T, integrante domain.Integrante) {
	t.Helper()
	if integrante != (domain.Integrante{}) {
		t.Errorf("se esperaba domain.Integrante{} ante un error, se obtuvo %+v", integrante)
	}
}

func TestUS03_AgregarIntegrante_Exitoso(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	esperado := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
	verificarIntegrante(t, agregado, esperado)
	if len(integrantes.integrantes["p1"]) != 1 {
		t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(integrantes.integrantes["p1"]))
	}
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), esperado)
}

func TestUS03_AgregarIntegrante_SegundoIntegrante(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})
	integrantes.integrantes["p1"] = []domain.Integrante{{Nombre: "Dana Guzmán", Rol: "Desarrollador"}}

	agregado, err := servicio.AgregarIntegrante("p1", "Victoria Caro Boldrini", "Tester")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregado, domain.Integrante{Nombre: "Victoria Caro Boldrini", Rol: "Tester"})
	if len(integrantes.integrantes["p1"]) != 2 {
		t.Fatalf("se esperaban 2 integrantes persistidos, hay %d", len(integrantes.integrantes["p1"]))
	}
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Victoria Caro Boldrini"), domain.Integrante{Nombre: "Victoria Caro Boldrini", Rol: "Tester"})
}

func TestUS03_AgregarIntegrante_NombreConEspaciosExternos(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})

	agregado, err := servicio.AgregarIntegrante("p1", "  Ana Paz  ", "Desarrollador Backend")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	esperado := domain.Integrante{Nombre: "Ana Paz", Rol: "Desarrollador Backend"}
	verificarIntegrante(t, agregado, esperado)
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Ana Paz"), esperado)
}

func TestUS03_AgregarIntegrante_PrimerAltaSinArchivoIntegrantes(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})
	if len(integrantes.integrantes["p1"]) != 0 {
		t.Fatalf("no debería haber integrantes antes del alta")
	}

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregado, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
}

func TestUS03_AgregarIntegrante_MismoNombreEnOtroProyecto(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		ProyectoRegistrado{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	)

	agregadoP1, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante(p1) devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregadoP1, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})

	agregadoP2, err := servicio.AgregarIntegrante("p2", "Dana Guzmán", "Tester")
	if err != nil {
		t.Fatalf("AgregarIntegrante(p2) devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregadoP2, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Tester"})

	if len(integrantes.integrantes["p1"])+len(integrantes.integrantes["p2"]) != 2 {
		t.Fatalf("se esperaban 2 integrantes persistidos")
	}
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p2", "Dana Guzmán"), domain.Integrante{Nombre: "Dana Guzmán", Rol: "Tester"})
}

func TestUS03_AgregarIntegrante_DuplicadoExacto(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})
	existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
	integrantes.integrantes["p1"] = []domain.Integrante{existente}

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un duplicado exacto debería devolver un error")
	}
	if !errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrIntegranteDuplicado)
	}
	verificarIntegranteVacio(t, agregado)
	if len(integrantes.integrantes["p1"]) != 1 {
		t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(integrantes.integrantes["p1"]))
	}
	verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), existente)
}

func TestUS03_AgregarIntegrante_RolVacioNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		rol         string
	}{
		{descripcion: "rol vacío", rol: ""},
		{descripcion: "rol compuesto solo por espacios", rol: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})
			existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
			integrantes.integrantes["p1"] = []domain.Integrante{existente}

			agregado, err := servicio.AgregarIntegrante("p1", "Victoria Caro Boldrini", caso.rol)
			if err == nil {
				t.Fatal("AgregarIntegrante() con rol vacío debería devolver un error")
			}
			if !errors.Is(err, domain.ErrRolIntegranteVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrRolIntegranteVacio)
			}
			verificarIntegranteVacio(t, agregado)
			if integrantes.agregados != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", integrantes.agregados)
			}
			if len(integrantes.integrantes["p1"]) != 1 {
				t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(integrantes.integrantes["p1"]))
			}
			verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), existente)
		})
	}
}

func TestUS03_AgregarIntegrante_ErrorNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		rol           string
		errorEsperado error
	}{
		{descripcion: "nombre vacío", nombre: "", rol: "Desarrollador", errorEsperado: domain.ErrNombreIntegranteVacio},
		{descripcion: "rol vacío", nombre: "Victoria Caro Boldrini", rol: "", errorEsperado: domain.ErrRolIntegranteVacio},
		{descripcion: "duplicado exacto", nombre: "Dana Guzmán", rol: "Desarrollador", errorEsperado: domain.ErrIntegranteDuplicado},
		{descripcion: "mismo nombre con distinto rol", nombre: "Dana Guzmán", rol: "Tester", errorEsperado: domain.ErrIntegranteDuplicado},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			servicio, _, integrantes := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})
			existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
			integrantes.integrantes["p1"] = []domain.Integrante{existente}

			agregado, err := servicio.AgregarIntegrante("p1", caso.nombre, caso.rol)
			if err == nil {
				t.Fatal("AgregarIntegrante() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarIntegranteVacio(t, agregado)
			if integrantes.agregados != 0 {
				t.Errorf("Agregar fue invocado %d veces ante un error; no debía escribirse nada", integrantes.agregados)
			}
			if len(integrantes.integrantes["p1"]) != 1 {
				t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(integrantes.integrantes["p1"]))
			}
			verificarIntegrante(t, buscarIntegranteEnFake(t, integrantes, "p1", "Dana Guzmán"), existente)
		})
	}
}

func TestUS03_AgregarIntegrante_ProyectoNoEncontrado(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		ProyectoRegistrado{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	)

	agregado, err := servicio.AgregarIntegrante("ausente", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarIntegranteVacio(t, agregado)
	if integrantes.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", integrantes.agregados)
	}
}

func TestUS03_AgregarIntegrante_ProyectoInexistenteAntesQueNombreVacio(t *testing.T) {
	servicio, _, _ := nuevoServicioIntegrantesConDatos(ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)})

	agregado, err := servicio.AgregarIntegrante("ausente", "", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrNombreIntegranteVacio) {
		t.Errorf("err = %v, no debería devolverse ErrNombreIntegranteVacio si el proyecto no existe", err)
	}
	verificarIntegranteVacio(t, agregado)
}

func TestUS03_AgregarIntegrante_ArchivoProyectosInexistente(t *testing.T) {
	servicio, _, integrantes := nuevoServicioIntegrantesConDatos()

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() sin archivo de proyectos debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarIntegranteVacio(t, agregado)
	if integrantes.agregados != 0 {
		t.Errorf("no debía escribirse nada; agregados = %d", integrantes.agregados)
	}
}

var errInfraCorruptoUS03 = errors.New("contenido inválido")
var errFalloEscrituraUS03 = errors.New("fallo simulado de escritura")

func TestUS03_AgregarIntegrante_JSONCorruptoNoEsNoEncontrado(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.existeErr = errInfraCorruptoUS03
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con JSON corrupto debería devolver un error")
	}
	if !errors.Is(err, errInfraCorruptoUS03) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errInfraCorruptoUS03)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarIntegranteVacio(t, agregado)
}

func TestUS03_AgregarIntegrante_FalloDeEscrituraNoInformaExito(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)
	integrantes := nuevoFakeIntegrantes()
	integrantes.agregarErr = errFalloEscrituraUS03
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con fallo de escritura debería devolver un error")
	}
	if !errors.Is(err, errFalloEscrituraUS03) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errFalloEscrituraUS03)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería ser ErrProyectoNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrNombreIntegranteVacio) || errors.Is(err, domain.ErrRolIntegranteVacio) || errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Errorf("err = %v, no debería ser un error de validación de dominio", err)
	}
	verificarIntegranteVacio(t, agregado)
}

func TestUS03_AgregarIntegrante_JSONCorruptoIntegrantes(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)
	integrantes := nuevoFakeIntegrantes()
	integrantes.listarErr = errInfraCorruptoUS03
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con integrantes.json corrupto debería devolver un error")
	}
	if !errors.Is(err, errInfraCorruptoUS03) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errInfraCorruptoUS03)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarIntegranteVacio(t, agregado)
	if integrantes.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", integrantes.agregados)
	}
}
