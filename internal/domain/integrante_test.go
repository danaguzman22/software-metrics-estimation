package domain_test

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// verificarIntegranteVacio comprueba que un error de validación devuelve el
// valor cero de domain.Integrante.
func verificarIntegranteVacio(t *testing.T, integrante domain.Integrante) {
	t.Helper()
	if integrante != (domain.Integrante{}) {
		t.Errorf("se esperaba domain.Integrante{} ante un error, se obtuvo %+v", integrante)
	}
}

// CA3 / R4 / V1: un nombre vacío o compuesto solo por espacios devuelve
// ErrNombreIntegranteVacio.
func TestUS03_NombreVacio(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "nombre vacío",
			nombre:      "",
		},
		{
			// C3: un nombre compuesto únicamente por espacios se considera vacío.
			descripcion: "nombre compuesto solo por espacios",
			nombre:      "   ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			integrante, err := domain.NuevoIntegrante(caso.nombre, "Desarrollador")
			if err == nil {
				t.Fatal("NuevoIntegrante() con nombre vacío o solo espacios debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
			}
			verificarIntegranteVacio(t, integrante)
		})
	}
}

// CA4 / R7 / V4: un rol vacío o compuesto solo por espacios devuelve
// ErrRolIntegranteVacio.
func TestUS03_RolVacio(t *testing.T) {
	casos := []struct {
		descripcion string
		rol         string
	}{
		{
			descripcion: "rol vacío",
			rol:         "",
		},
		{
			// C4: un rol compuesto únicamente por espacios se considera vacío.
			descripcion: "rol compuesto solo por espacios",
			rol:         "   ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			integrante, err := domain.NuevoIntegrante("Dana Guzmán", caso.rol)
			if err == nil {
				t.Fatal("NuevoIntegrante() con rol vacío o solo espacios debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrRolIntegranteVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrRolIntegranteVacio)
			}
			verificarIntegranteVacio(t, integrante)
		})
	}
}

// CA1 / C6: el nombre y el rol se guardan sin espacios al inicio y al final,
// respetando mayúsculas.
func TestUS03_NuevoIntegrante_RecortaEspaciosExternos(t *testing.T) {
	casos := []struct {
		descripcion    string
		nombre         string
		rol            string
		nombreGuardado string
		rolGuardado    string
	}{
		{
			descripcion:    "nombre y rol con espacios externos",
			nombre:         "  Ana Paz  ",
			rol:            "  Revisor  ",
			nombreGuardado: "Ana Paz",
			rolGuardado:    "Revisor",
		},
		{
			descripcion:    "mayúsculas respetadas tras el recorte",
			nombre:         "  Dana Guzmán  ",
			rol:            "  Desarrollador Backend  ",
			nombreGuardado: "Dana Guzmán",
			rolGuardado:    "Desarrollador Backend",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			integrante, err := domain.NuevoIntegrante(caso.nombre, caso.rol)
			if err != nil {
				t.Fatalf("NuevoIntegrante() devolvió un error inesperado: %v", err)
			}
			if integrante.Nombre != caso.nombreGuardado {
				t.Errorf("Nombre = %q, se esperaba %q", integrante.Nombre, caso.nombreGuardado)
			}
			if integrante.Rol != caso.rolGuardado {
				t.Errorf("Rol = %q, se esperaba %q", integrante.Rol, caso.rolGuardado)
			}
		})
	}
}

// CA2 / R2 / C2: la comparación de duplicado ignora mayúsculas.
func TestUS03_Duplicado_IgnoraMayusculas(t *testing.T) {
	existentes := []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}

	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "candidato en minúsculas",
			nombre:      "dana guzmán",
		},
		{
			descripcion: "candidato en mayúsculas",
			nombre:      "DANA GUZMÁN",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			if !domain.EsDuplicado(existentes, caso.nombre) {
				t.Errorf("EsDuplicado(existentes, %q) = false, se esperaba true", caso.nombre)
			}
		})
	}
}

// CA2 / R2 / C2: la comparación de duplicado ignora espacios al inicio y al
// final.
func TestUS03_Duplicado_IgnoraEspaciosExtremos(t *testing.T) {
	existentes := []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
	}

	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "candidato con espacios al inicio y al final",
			nombre:      "  Dana Guzmán  ",
		},
		{
			descripcion: "candidato con espacios externos y mayúsculas distintas",
			nombre:      "  dana guzmán  ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			if !domain.EsDuplicado(existentes, caso.nombre) {
				t.Errorf("EsDuplicado(existentes, %q) = false, se esperaba true", caso.nombre)
			}
		})
	}
}

// CA2 / R2 / C8: el rol no participa en la comparación de duplicado.
func TestUS03_Duplicado_IgnoraRol(t *testing.T) {
	existentes := []domain.Integrante{
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
		{Nombre: "Victoria Caro", Rol: "Revisor"},
	}

	casos := []struct {
		descripcion string
		nombre      string
		esperado    bool
	}{
		{
			// C8: el mismo nombre con distinto rol en el alta también es duplicado.
			descripcion: "nombre ya registrado",
			nombre:      "Dana Guzmán",
			esperado:    true,
		},
		{
			// C8: un nombre no registrado no es duplicado; el rol no participa.
			descripcion: "nombre no registrado",
			nombre:      "Ana Paz",
			esperado:    false,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			if obtenido := domain.EsDuplicado(existentes, caso.nombre); obtenido != caso.esperado {
				t.Errorf("EsDuplicado(existentes, %q) = %v, se esperaba %v", caso.nombre, obtenido, caso.esperado)
			}
		})
	}
}

// CA2 / R2 / C8: los espacios internos sí distinguen; nombres que difieren
// solo en espacios internos no son duplicados.
func TestUS03_Duplicado_EspaciosInternosDistintos(t *testing.T) {
	existentes := []domain.Integrante{
		{Nombre: "Ana Paz", Rol: "Revisor"},
	}

	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			// C8: "AnaPaz" y "Ana Paz" son nombres distintos.
			descripcion: "sin espacio interno",
			nombre:      "AnaPaz",
		},
		{
			descripcion: "con doble espacio interno",
			nombre:      "Ana  Paz",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			if domain.EsDuplicado(existentes, caso.nombre) {
				t.Errorf("EsDuplicado(existentes, %q) = true, se esperaba false", caso.nombre)
			}
		})
	}
}
