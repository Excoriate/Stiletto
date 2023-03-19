package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func FileExist(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", filePath)
	}

	return nil
}

func IsValidYAML(filePath string) error {
	extension := filepath.Ext(filePath)

	if extension != ".yml" && extension != ".yaml" {
		return errors.New("the file is not a valid YAML file, allowed extensions are .yml and .yaml")
	}

	return nil
}

func IsValidYAMLAgainstSchema(yamlFilePath, schemaJSONPath string) error {
	// Read the YAML file
	yamlContent, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Convert the YAML content to JSON
	var yamlData interface{}
	err = yaml.Unmarshal(yamlContent, &yamlData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	jsonContent, err := json.Marshal(yamlData)
	if err != nil {
		return fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	// Read the schema.json file
	schemaContent, err := os.ReadFile(schemaJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Load the JSON schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaContent)
	documentLoader := gojsonschema.NewBytesLoader(jsonContent)

	// Validate the YAML (as JSON) against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate YAML against JSON schema: %w", err)
	}

	if !result.Valid() {
		var validationErrors []string
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.String())
		}
		return errors.New("validation failed:\n" + strings.Join(validationErrors, "\n"))
	}

	return nil
}
