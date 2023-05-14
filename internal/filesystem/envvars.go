package filesystem

import (
	"errors"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"os"
	"strings"
)

type EnvVars map[string]string

// AreEnvVarsExportedAndSet IsEnvVarSetOrExported checks if the environment variable exist,
//and also if it is exported or set.
func AreEnvVarsExportedAndSet(keys []string) error {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); !ok || value == "" {
			return errors.New(fmt.Sprintf("Environment variable %s is not set or has an empty value", key))
		}
	}
	return nil
}

// DoEnvVarsExist checks if the environment variables exist.
func DoEnvVarsExist(keys []string) error {
	for _, key := range keys {
		_, ok := os.LookupEnv(key)
		if !ok {
			return errors.New(fmt.Sprintf("Environment variable %s does not exist", key))
		}
	}
	return nil
}

// AreEnvVarsSet checks if the environment variables have non-empty values.
func AreEnvVarsSet(keys []string) error {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok || value == "" {
			return errors.New(fmt.Sprintf("Environment variable %s is not set or has an empty value", key))
		}
	}
	return nil
}

// FetchEnvVarsAsMap checks if the environment variables exist and returns them as a map.
func FetchEnvVarsAsMap(keys []string, optionalKeys []string) (EnvVars, error) {
	result := make(EnvVars)

	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok {
			if common.IsKeyInMapOptional(key, optionalKeys) {
				continue
			} else {
				return nil, errors.New(fmt.Sprintf("Environment variable %s does not exist", key))
			}
		}
		result[key] = value
	}

	return result, nil
}

// ScanAWSCredentialsEnvVars scans the environment variables for AWS credentials.
func ScanAWSCredentialsEnvVars() (EnvVars, error) {
	keys := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
		"AWS_SECURITY_TOKEN",
		"AWS_DEFAULT_REGION",
		"AWS_REGION",
		"AWS_PROFILE",
	}

	mandatoryKeys := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}

	if err := AreEnvVarsSet(mandatoryKeys); err != nil {
		return nil, fmt.Errorf("AWS credentials are not set as environment variables: %w", err)
	}

	// Here' it's scanned the required environment variables
	envs, err := FetchEnvVarsAsMap(keys, []string{"AWS_PROFILE", "AWS_SECURITY_TOKEN",
		"AWS_SESSION_TOKEN", "AWS_REGION", "AWS_DEFAULT_REGION"})

	if err != nil {
		return nil, err
	}

	// Precedence rules
	if defaultRegion, ok := envs["AWS_DEFAULT_REGION"]; ok {
		envs["AWS_REGION"] = defaultRegion
	}

	if accessKey, ok := envs["AWS_ACCESS_KEY_ID"]; ok {
		envs["AWS_SESSION_TOKEN"] = accessKey
	}

	if secretKey, ok := envs["AWS_SECRET_ACCESS_KEY"]; ok {
		envs["AWS_SECURITY_TOKEN"] = secretKey
	}

	// If keys are set, but neither of the regions are set, set the region to us-east-1
	if _, ok := envs["AWS_REGION"]; !ok {
		if _, ok := envs["AWS_DEFAULT_REGION"]; !ok {
			envs["AWS_REGION"] = "us-east-1"
		}
	}

	return envs, nil
}

// FetchEnvVarsWithPrefix fetches environment variables that start with the specified prefix
// and returns an error if any of the variables either do not exist or have an empty value.
func FetchEnvVarsWithPrefix(prefix string) (EnvVars, error) {
	result := make(EnvVars)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]

		if strings.HasPrefix(key, prefix) {
			value := pair[1]
			if value == "" {
				return nil, errors.New(fmt.Sprintf("Environment variable %s has an empty value", key))
			}
			result[key] = value
		}
	}

	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("No environment variables with the prefix %s found", prefix))
	}

	return result, nil
}

// ScanTerraformEnvVars fetches environment variables that start with the prefix "TF_VAR_"
func ScanTerraformEnvVars() (EnvVars, error) {
	return FetchEnvVarsWithPrefix("TF_VAR_")
}

// FetchAWSEnvVars fetches environment variables that start with the prefix "AWS_"
func FetchAWSEnvVars() (EnvVars, error) {
	return FetchEnvVarsWithPrefix("AWS_")
}

// AreEnvVarsConsistent checks if the environment variables have non-empty values.
func AreEnvVarsConsistent(envVars EnvVars) error {
	for key, value := range envVars {
		if value == "" {
			return errors.New(fmt.Sprintf("The value for the environment variable %s is empty", key))
		}
	}
	return nil
}

func MergeEnvVars(envVars ...EnvVars) EnvVars {
	result := make(EnvVars)

	for _, env := range envVars {
		for key, value := range env {
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}

	return result
}
