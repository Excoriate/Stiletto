package config

type PipelineOptions struct {
	WorkDir               string
	WorkDirPath           string
	MountDir              string
	MountDirPath          string
	TargetDir             string
	TargetDirPath         string
	TaskName              string
	EnvVarsDotEnvFilePath string
	EnvVarsToScanAndSet   []string
	EnvVarsToScanByPrefix []string
	EnvKeyValuePairsToSet map[string]string
	EnvVarsFromDotEnvFile map[string]string
	EnvVarsAWSKeysToScan  map[string]string
	// Automatic discovery of environment variables, for well-known use cases.
	IsAWSEnvVarKeysToScanEnabled   bool
	IsTerraformVarsScanEnabled     bool
	IsAllEnvVarsToScanEnabled      bool
	IsEnvVarsToScanFromDotEnvFile  bool
	IsEnvVarsToScanByPrefix        bool
	InitDaggerWithWorkDirByDefault bool
}

type PipelineDirs struct {
	Dir  string
	Path string
}
