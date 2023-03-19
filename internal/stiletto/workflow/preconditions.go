package workflow

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/config"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/logger"
	"github.com/Excoriate/stiletto/pkg/common"
	"github.com/Excoriate/stiletto/pkg/filesystem"
)

// AreWorkflowFileValidationsPassed checks if the workflow file is valid,
//running basic validations from the filesystem package.
func AreWorkflowFileValidationsPassed(l logger.Logger, filepath string) error {
	filepathNormalised := common.NormaliseStringLower(filepath)
	l.LogDebug("Checking if workflow file is valid: %s", filepathNormalised)

	// File doesn't exist.
	if err := filesystem.FileExist(filepathNormalised); err != nil {
		return errors.NewWorkflowFileNotFoundError(filepathNormalised)
	}

	// File isn't a proper stiletto workflow (.yml).
	if err := filesystem.IsValidYAML(filepathNormalised); err != nil {
		return errors.NewWorkflowInvalidExtensionError(filepathNormalised)
	}

	l.LogInfo("Workflow file is valid: %s", filepathNormalised)

	return nil
}

func IsWorkflowSchemaCompliant(l logger.Logger, filepath string) error {
	jsonSchema, err := config.GetSchemaJSONPath()
	l.LogDebug(fmt.Sprintf("JSON schema path resolved is: %s", jsonSchema), "", jsonSchema)

	if err != nil {
		return errors.NewInternalStilettoError(fmt.Sprintf("Could not get schema JSON path: %s", err.Error()))
	}

	if schemaErr := filesystem.IsValidYAMLAgainstSchema(filepath, jsonSchema); schemaErr != nil {
		return errors.NewWorkflowSchemaVerificationError(filepath, schemaErr)
	}

	l.LogInfo("Workflow file is compliant with the schema: %s", filepath)

	return nil
}
