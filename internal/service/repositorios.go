package service

import "github.com/danaguzman22/software-metrics-estimation/internal/domain"

// ProyectoRegistrado une el ID persistido del proyecto con su agregado de
// dominio. El ID no forma parte de domain.Proyecto (decisión D11 de US-02):
// vive en la capa de servicio y en los contratos que lo devuelven.
// Conserva los json tags actuales para mantener el formato de archivos.
type ProyectoRegistrado struct {
	ID       string          `json:"id"`
	Proyecto domain.Proyecto `json:"proyecto"`
}

// RepositorioProyectos persiste el agregado de proyectos.
type RepositorioProyectos interface {
	// Obtener devuelve el registro {ID, Proyecto}.
	// domain.ErrProyectoNoEncontrado si el ID no existe.
	Obtener(id string) (ProyectoRegistrado, error)
	// Existe informa si hay un proyecto con ese ID, sin cargar el agregado.
	// Un repositorio vacío (sin datos) es Existe == false, no un error.
	Existe(id string) (bool, error)
	// Listar devuelve todos los registros de proyectos.
	Listar() ([]ProyectoRegistrado, error)
	// Actualizar reemplaza el proyecto con ese ID (el ID nunca cambia).
	// domain.ErrProyectoNoEncontrado si el ID no existe.
	Actualizar(id string, proyecto domain.Proyecto) error
}

// RepositorioIntegrantes persiste los integrantes por proyecto.
type RepositorioIntegrantes interface {
	// ListarPorProyecto devuelve los integrantes del proyecto en orden de
	// inserción (D7).
	ListarPorProyecto(proyectoID string) ([]domain.Integrante, error)
	// Agregar registra el integrante en el proyecto.
	// No valida duplicados: el chequeo es del servicio (ListarPorProyecto +
	// domain.EsDuplicado, D3). Ante una violación de unicidad del motor
	// (SQLSTATE 23505, defensa ante carreras) devuelve
	// domain.ErrIntegranteDuplicado (nombre del integrante dentro del
	// proyecto, sin distinguir mayúsculas).
	Agregar(proyectoID string, integrante domain.Integrante) error
}

// RepositorioBacklog persiste los ítems del Product Backlog por proyecto.
type RepositorioBacklog interface {
	// AgregarItem asigna el ID "HU-NN" de forma atómica (siguiente al mayor
	// número del proyecto, D4) y devuelve el ítem persistido con su ID.
	AgregarItem(proyectoID string, item domain.ItemBacklog) (domain.ItemBacklog, error)
	// ObtenerItem devuelve el ítem identificado por itemID dentro del proyecto.
	// domain.ErrItemNoEncontrado si no existe en ESE proyecto.
	ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error)
	// ListarPorProyecto devuelve los ítems del proyecto ordenados por número
	// (HU-01, HU-02, ..., D7).
	ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error)
	// ActualizarItem reemplaza el ítem identificado por item.ID dentro del
	// proyecto. domain.ErrItemNoEncontrado si no existe en ese proyecto.
	ActualizarItem(proyectoID string, item domain.ItemBacklog) error
}
