package domain

import "errors"

// Errores de dominio de integrantes (D6).
var (
	// ErrNombreIntegranteVacio indica que el nombre del integrante está
	// vacío o compuesto únicamente por espacios.
	ErrNombreIntegranteVacio = errors.New("el nombre del integrante no puede estar vacío")
	// ErrRolIntegranteVacio indica que el rol del integrante está vacío o
	// compuesto únicamente por espacios.
	ErrRolIntegranteVacio = errors.New("el rol del integrante no puede estar vacío")
	// ErrIntegranteDuplicado indica que ya hay un integrante con ese nombre
	// en el mismo proyecto.
	ErrIntegranteDuplicado = errors.New("el integrante ya está registrado en el proyecto")
)

// Integrante representa un miembro del equipo de un proyecto, con su nombre
// y su rol (D4).
type Integrante struct {
	Nombre string
	Rol    string
}

// NuevoIntegrante crea un integrante a partir de su nombre y su rol.
//
// El nombre no puede estar vacío ni compuesto únicamente por espacios (C3)
// y el rol tampoco (C4); el nombre se valida antes que el rol (C7). Ambos se
// guardan sin espacios al inicio y al final, respetando mayúsculas (C6).
// Ante cualquier error se devuelve Integrante{}.
//
// Stub de la fase RED de US-03: sin lógica; devuelve Integrante{} y un error
// "no implementado" para que los tests fallen por separado.
func NuevoIntegrante(nombre, rol string) (Integrante, error) {
	return Integrante{}, errors.New("no implementado")
}

// EsDuplicado indica si el nombre ya está entre los integrantes existentes
// de un mismo proyecto, ignorando mayúsculas y espacios al inicio y al final
// (C2). El rol no participa en la comparación y los espacios internos sí
// distinguen (C8).
//
// Stub de la fase RED de US-03: sin lógica; devuelve siempre false.
func EsDuplicado(existentes []Integrante, nombre string) bool {
	return false
}
