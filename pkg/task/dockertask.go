package task

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/stiletto/internal/daggerio"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/job"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

type DockerTask struct {
	Init     *InitOptions
	Cfg      *Task
	Actions  []string
	UXPrefix string
}

func (t *DockerTask) RunCmdInContainer(container *dagger.Container, commands [][]string,
	stdOutEnabled bool, ctx context.Context) error {
	ux := tui.NewTUIMessage()

	if len(commands) != 0 {
		for _, cmd := range commands {
			ux.ShowInfo(t.UXPrefix, fmt.Sprintf("Running command %s", cmd))

			if !stdOutEnabled {
				_, err := container.
					WithExec(cmd).
					ExitCode(ctx)

				if err != nil {
					ux.ShowError(t.UXPrefix, fmt.Sprintf("Failed to run command %s", cmd), err)
					return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run command %s",
						cmd), err)
				}
			} else {
				_, err := container.
					WithExec(cmd).
					Stdout(ctx)

				if err != nil {
					ux.ShowError(t.UXPrefix, fmt.Sprintf("Failed to run command %s", cmd), err)
					return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run command %s",
						cmd), err)
				}
			}
		}
	} else {
		_, err := container.
			WithExec([]string{"ls", "-ltrh"}).
			ExitCode(ctx)

		if err != nil {
			ux.ShowError(t.UXPrefix, fmt.Sprintf("Failed to run command %s", "ls -la"), err)
			return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run command %s",
				"ls -ltrh"), err)
		}
	}

	return nil
}

func (t *DockerTask) SetEnvVarsFromJob(container *dagger.Container) (*dagger.Container, error) {
	ux := t.Cfg.PipelineCfg.UXMessage
	j := t.GetJob()

	awsEnvVars := j.EnvVarsAWSScanned
	tfEnvVars := j.EnvVarsTerraformScanned
	customEnvVars := j.EnvVarsCustomScanned
	envFromHost := j.EnvVarsAllScanned
	specificToSet := j.EnvVarsToSet
	dotEnvEnvVars := j.EnvVarsFromDotEnvFile

	c := container

	var mergedEnvVars map[string]string

	if len(awsEnvVars) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting AWS environment variables from the job")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, awsEnvVars)
	} else {
		ux.ShowInfo(t.UXPrefix, "No AWS environment variables to set from the job")
	}

	if len(tfEnvVars) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting Terraform environment variables from the job")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, tfEnvVars)
	} else {
		ux.ShowInfo(t.UXPrefix, "No Terraform environment variables to set from the job")
	}

	if len(customEnvVars) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting custom environment variables from the job")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, customEnvVars)
	} else {
		ux.ShowInfo(t.UXPrefix, "No custom environment variables to set from the job")
	}

	if len(envFromHost) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting environment variables from the host")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, envFromHost)
	} else {
		ux.ShowInfo(t.UXPrefix, "No environment variables to set from the host")
	}

	if len(specificToSet) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting specific environment variables")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, specificToSet)
	} else {
		ux.ShowInfo(t.UXPrefix, "No specific environment variables to set")
	}

	if len(dotEnvEnvVars) > 0 {
		ux.ShowInfo(t.UXPrefix, "Setting environment variables from .env file")
		mergedEnvVars = filesystem.MergeEnvVars(mergedEnvVars, dotEnvEnvVars)
	} else {
		ux.ShowInfo(t.UXPrefix, "No environment variables to set from .env file")
	}

	finalContainer, err := daggerio.SetEnvVarsInContainer(c, mergedEnvVars)
	if err != nil {
		return nil, err
	}

	return finalContainer, nil
}

func (t *DockerTask) MountDir(targetDir string, client *dagger.Client, container *dagger.
Container,
	filesPreRequisites []string, ctx context.Context) (*dagger.Container, error) {
	ux := tui.NewTUIMessage()

	if targetDir == "" {
		ux.ShowWarning(t.UXPrefix, "An empty directory was passed to be a Target directory ("+
			"also known as Execution path), "+
			"hence the default working directory will be used resolved from the '.' value")

		targetDir = "."
	}

	if targetDir != "." && len(filesPreRequisites) > 0 {
		ux.ShowInfo(t.UXPrefix, "The target directory is not the working directory, "+
			"therefore the files pre-requisites will be verified before mounting the directory")

		if err := daggerio.VerifyFileEntriesInMountedDir(client, targetDir,
			filesPreRequisites, ctx); err != nil {
			ux.ShowError(t.UXPrefix, "Failed to mount the directory", err)
			return nil, err
		}
	}

	workDirDagger, err := daggerio.GetDaggerDir(t.GetClient(), ".")

	if err != nil {
		ux.ShowError(t.UXPrefix,
			fmt.Sprintf("Failed to mount the working directory (with value '.'), failed "+
				"to build a dagger directory from the directory"), err)

		return nil, err
	}

	containerMounted, err := daggerio.MountDir(container, workDirDagger, targetDir)

	if err != nil {
		ux.ShowError(t.UXPrefix,
			fmt.Sprintf("Failed to mount directory %s", targetDir), err)

		return nil, err
	}

	return containerMounted, nil
}

func (t *DockerTask) GetPipelineUXLog() tui.TUIMessenger {
	return t.Cfg.PipelineCfg.UXMessage
}

func (t *DockerTask) GetClient() *dagger.Client {
	return t.Cfg.JobCfg.Client
}

func (t *DockerTask) GetPipeline() *pipeline.Config {
	return t.Cfg.PipelineCfg
}

func (t *DockerTask) GetJob() *job.Job {
	return t.Cfg.JobCfg
}

func (t *DockerTask) ConvertDir(c *dagger.Client, dir string) (*dagger.Directory, error) {
	return daggerio.GetDaggerDir(c, dir)
}

func (t *DockerTask) GetCoreTask() *Task {
	return t.Cfg
}

func (t *DockerTask) GetJobContainerImage() string {
	return t.Cfg.JobCfg.ContainerImageURL
}

func (t *DockerTask) GetJobContainerDefault() *dagger.Container {
	return t.Cfg.JobCfg.ContainerDefault
}

func (t *DockerTask) GetJobEnvVars() map[string]string {
	return t.Cfg.EnvVarsInheritFromJob
}

func (t *DockerTask) PushImage(addr string, container *dagger.
Container, dockerFileDir *dagger.Directory,
	ctx context.Context) error {
	containerBuilt := container.Build(dockerFileDir)
	_, err := daggerio.PushImage(containerBuilt, addr, ctx)

	if err != nil {
		return err
	}

	return nil
}

func (t *DockerTask) BuildImage(dockerFilePath string, container *dagger.Container,
	ctx context.Context) (*dagger.Container, error) {
	return daggerio.BuildImage(dockerFilePath, t.GetClient(), container)
}

func (t *DockerTask) AuthWithRegistry(c *dagger.Client, container *dagger.Container,
	opt daggerio.RegistryAuthOptions) (*dagger.Container, error) {
	return daggerio.AuthWithRegistry(c, container, opt)
}

func (t *DockerTask) SetEnvVars(envVars []map[string]string,
	container *dagger.Container) (*dagger.Container, error) {
	ux := t.Cfg.PipelineCfg.UXMessage

	if len(envVars) == 0 {
		ux.ShowInfo(t.UXPrefix, "There is no environment variables to be set in the container")
		return container, nil
	}

	var envVarsMerged map[string]string

	for _, envVar := range envVars {
		envVarsMerged = filesystem.MergeEnvVars(envVarsMerged, envVar)
	}

	return daggerio.SetEnvVarsInContainer(container, envVarsMerged)
}

func (t *DockerTask) GetContainer(fromImage string) (*dagger.Container,
	error) {
	if fromImage == "" {
		return t.Cfg.JobCfg.ContainerDefault, nil
	}

	return t.Cfg.JobCfg.Client.Container().From(fromImage), nil
}

func NewTaskDocker(coreTask *Task, actions []string,
	init *InitOptions, uxPrefix string) CoreTasker {

	return &DockerTask{
		Init:     init,
		Cfg:      coreTask,
		Actions:  actions,
		UXPrefix: uxPrefix,
	}
}
