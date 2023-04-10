package config

type PipelineOptions struct {
	WorkDir               string
	WorkDirPath           string
	MountDir              string
	MountDirPath          string
	TargetDir             string
	TargetDirPath         string
	TaskName              string
	EnvVarsToScanAndSet   []string
	EnvKeyValuePairsToSet map[string]string
	EnvVarsAWSKeysToScan  map[string]string
	// Automatic discovery of environment variables, for well-known use cases.
	IsAWSEnvVarKeysToScanEnabled   bool
	IsTerraformVarsScanEnabled     bool
	InitDaggerWithWorkDirByDefault bool
}
