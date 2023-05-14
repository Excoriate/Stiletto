package task

import (
	"context"
	"dagger.io/dagger"
	"github.com/Excoriate/stiletto/internal/daggerio"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/job"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

type CoreTasker interface {
	GetClient() *dagger.Client
	GetPipeline() *pipeline.Config
	GetPipelineUXLog() tui.TUIMessenger
	ConvertDir(c *dagger.Client, dir string) (*dagger.Directory, error)
	GetJob() *job.Job
	GetCoreTask() *Task
	GetJobContainerImage() string
	GetJobContainerDefault() *dagger.Container
	GetJobEnvVars() map[string]string
	SetEnvVars(envVars []map[string]string, container *dagger.Container) (*dagger.Container, error)
	AuthWithRegistry(c *dagger.Client, container *dagger.Container,
		opt daggerio.RegistryAuthOptions) (*dagger.Container,
		error)
	GetContainer(fromImage string) (*dagger.Container, error)
	BuildImage(dockerFilePath string, container *dagger.Container, ctx context.Context) (*dagger.Container, error)
	PushImage(addr string, container *dagger.Container,
		dockerFileDir *dagger.Directory, ctx context.Context) error
	MountDir(targetDir string, client *dagger.Client, container *dagger.
	Container,
		filesPreRequisites []string, ctx context.Context) (*dagger.Container, error)

	RunCmdInContainer(container *dagger.Container, commands [][]string,
		stdOutEnabled bool, ctx context.Context) error
}

type Runner struct {
	Init *InitOptions
	Cfg  *Task
}

type Task struct {
	// Identifiers.
	Id    string
	Name  string
	Stack string

	// Configuration
	PipelineCfg *pipeline.Config
	JobCfg      *job.Job

	// Specific attributes
	EnvVarsInheritFromJob map[string]string
	Dirs                  Dirs
	ContainerImageDefault string
	ContainerNameDefault  string
	ContainerDefault      *dagger.Container

	PreReqs PreRequisites
	Actions Actions

	// Output
	Result Output

	Ctx context.Context
}

type Dirs struct {
	RootDir         string
	WorkDir         string
	MountDir        string
	TargetDir       string
	RootDirDagger   *dagger.Directory
	WorkDirDagger   *dagger.Directory
	MountDirDagger  *dagger.Directory
	TargetDirDagger *dagger.Directory
}

type Output struct {
	Files        []*dagger.File
	Directories  []*dagger.Directory
	ExitCode     int
	DaggerOutput interface{}
	IsError      bool
}

type Actions struct {
	CustomCommands  []string
	DefaultCommands []string
}

type PreRequisites struct {
	Files       []string
	Directories []string
}
