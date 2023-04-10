package daggerio

import (
	"dagger.io/dagger"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/logger"
)

const ContainerMountPathPrefix = "/build"

func SetEnvVarsInContainer(c *dagger.Container, envVars map[string]string) (*dagger.Container,
	error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()

	if common.MapIsNulOrEmpty(envVars) {
		logPrinter.LogWarn("Dagger Container Configuration",
			"No environment variables are passed, skipping the environment variable configuration step")
		return c, nil
	}

	for k, v := range envVars {
		c = c.WithEnvVariable(k, v)
	}

	return c, nil
}
