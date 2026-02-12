package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetProjectRoot devuelve la ruta absoluta de la raíz del proyecto
func GetProjectRoot() string {
	// 1. Intentamos obtener la ruta del archivo actual que se está ejecutando
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	// 2. Subimos por el árbol de directorios buscando go.mod
	for {
		if _, err := os.Stat(filepath.Join(basePath, "go.mod")); err == nil {
			return basePath
		}

		parent := filepath.Dir(basePath)
		if parent == basePath {
			// Hemos llegado a la raíz del sistema de archivos sin encontrar go.mod
			return ""
		}
		basePath = parent
	}
}
