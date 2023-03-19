package config

import (
	"errors"
	"path/filepath"
	"runtime"
)

// GetSchemaJSONPath returns the path to the schema.json file.
func GetSchemaJSONPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get current file path")
	}

	// Path to the schema.json file
	schemaPath := filepath.Join(filepath.Dir(currentFilePath), "workflow", "schema.json")

	return schemaPath, nil
}
