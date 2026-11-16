// Package storage ofrece persistencia genérica en archivos JSON.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNotFound    = errors.New("archivo no encontrado")
	ErrInvalidName = errors.New("nombre inválido")
	ErrCorrupt     = errors.New("contenido JSON inválido")
)

// JSONStore guarda y carga valores como archivos <nombre>.json dentro de un directorio.
type JSONStore struct {
	dir string
}

// NewJSONStore crea el almacenamiento y el directorio base si no existe.
func NewJSONStore(dir string) (*JSONStore, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("storage: el directorio no puede estar vacío")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("storage: crear directorio: %w", err)
	}
	return &JSONStore{dir: dir}, nil
}

func (s *JSONStore) path(name string) (string, error) {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", ErrInvalidName
	}
	return filepath.Join(s.dir, name+".json"), nil
}

// Save serializa v y lo escribe de forma atómica (archivo temporal + rename).
func (s *JSONStore) Save(name string, v any) error {
	path, err := s.path(name)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("storage: serializar: %w", err)
	}
	tmp, err := os.CreateTemp(s.dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("storage: crear temporal: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("storage: escribir: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("storage: cerrar: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("storage: renombrar: %w", err)
	}
	return nil
}

// Load lee <name>.json y lo deserializa en v (debe ser un puntero).
func (s *JSONStore) Load(name string, v any) error {
	path, err := s.path(name)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return fmt.Errorf("storage: leer: %w", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		var invalid *json.InvalidUnmarshalError
		if errors.As(err, &invalid) {
			return fmt.Errorf("storage: destino inválido: %w", err)
		}
		return fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	return nil
}
