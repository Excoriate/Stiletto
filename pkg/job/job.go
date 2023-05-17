package job

import (
	"context"
	"dagger.io/dagger"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

type InitOptions struct {
	Name        string
	Stack       string
	PipelineCfg *pipeline.Config

	// Directories that the task will use.
	WorkDir   string
	MountDir  string
	TargetDir string

	// Scanned Environment variables to resolve, and set.
	ScanAWSEnvVars       bool
	ScanTerraformEnvVars bool
	EnvVarsToSet         map[string]string
	EnvVarsToScan        []string
	DotEnvFile           string
}

type Job struct {
	// Identifiers.
	Id    string
	Name  string
	Stack string

	// PipelineCfg client.
	PipelineCfg *pipeline.Config
	Client      *dagger.Client

	// Dagger directories
	RootDir   *dagger.Directory // Normally should be the same as the workDir
	WorkDir   *dagger.Directory
	MountDir  *dagger.Directory
	TargetDir *dagger.Directory

	RootDirPath   string
	WorkDirPath   string
	MountDirPath  string
	TargetDirPath string

	// Container configuration.
	ContainerImageURL string
	ContainerDefault  *dagger.Container

	// Scanned Environment variables to resolve, and set.
	EnvVarsAWSScanned       map[string]string
	EnvVarsTerraformScanned map[string]string
	EnvVarsCustomScanned    map[string]string
	EnvVarsAllScanned       map[string]string
	EnvVarsToSet            map[string]string
	EnvVarsFromDotEnvFile   map[string]string

	Ctx context.Context
}

type Runner interface {
	InitDagger() (*dagger.Client, error)
	InitContainerImage() (string, error)
	InitContainer(c *dagger.Client, imageURL string) (*dagger.Container, error)
	ScanEnvVarsAWSKeys(scanAWSVars bool) (map[string]string, error)
	ScanEnvVarsTerraform(scanTerraformVars bool) (map[string]string, error)
	ScanEnvVarsCustom(scanCustomVars []string) (map[string]string, error)
	ScanAllEnvVars() (map[string]string, error)
	ScanEnvVarsFromDotEnvFile(dotEnvFile string) (map[string]string, error)
	ValidatedEnvVarsPassed(envVarsToSet map[string]string) (map[string]string, error)
	BuildRootDir(client *dagger.Client) (*dagger.Directory, error)
	BuildWorkDir(client *dagger.Client, workDir string) (*dagger.Directory, error)
	BuildMountDir(client *dagger.Client, mountDir string) (*dagger.Directory, error)
	BuildTargetDir(client *dagger.Client, targetDir string) (*dagger.Directory, error)
}
