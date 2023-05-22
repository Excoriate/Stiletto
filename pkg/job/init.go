package job

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/daggerio"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

const uxPrefix = "JOB-INIT"

type Instance struct {
	InitOptions *InitOptions
	JobName     string
	JobId       string
}

func NewJob(p *pipeline.Config, new InitOptions) (*Job, error) {
	jobId := common.GetUUID()
	p.UXMessage.ShowInfo(uxPrefix, fmt.Sprintf("Initialising job: %s with id: %s",
		new.Name, jobId))

	i := &Instance{
		InitOptions: &new,
		JobName:     new.Name,
		JobId:       jobId,
	}

	// 1. Init the dagger engine.
	c, err := i.InitDagger()
	if err != nil {
		return nil, err
	}

	// 2. Get the container image.
	im, err := i.InitContainerImage()
	if err != nil {
		return nil, err
	}

	// 3. Init the container.
	ct, err := i.InitContainer(c, im)
	if err != nil {
		return nil, err
	}

	// 4. Scan (if applicable) AWS keys environment variables.
	awsEnvVars, err := i.ScanEnvVarsAWSKeys(new.ScanAWSEnvVars)
	if err != nil {
		return nil, err
	}

	// 5. Scan (if applicable) Terraform environment variables.
	terraformEnvVars, err := i.ScanEnvVarsTerraform(new.ScanTerraformEnvVars)
	if err != nil {
		return nil, err
	}

	// 6. Scan (if applicable) custom environment variables.
	customEnvVars, err := i.ScanEnvVarsCustom(new.EnvVarsToScan)
	if err != nil {
		return nil, err
	}

	// 7. Scan (if applicable) env vars from dotenv File
	envVarsFromDotEnv, err := i.ScanEnvVarsFromDotEnvFile(new.DotEnvFile)
	if err != nil {
		return nil, err
	}

	// 8. Scan (if applicable) env vars with prefix
	envVarsFromPrefix, err := i.ScanEnvVarsFromPrefix(new.EnvVarsWithPrefixToScan)
	if err != nil {
		return nil, err
	}

	// 8. Validate and set environment variables.
	envVarsToSet, err := i.ValidatedEnvVarsPassed(new.EnvVarsToSet)
	if err != nil {
		return nil, err
	}

	// 9 RootDir in dagger format.
	rootDir, err := i.BuildRootDir(c)
	if err != nil {
		return nil, err
	}

	// 10. WorkDir in dagger format.
	workDir, err := i.BuildWorkDir(c, new.WorkDir)
	if err != nil {
		return nil, err
	}

	// 11. MountDir in dagger format.
	mountDirPath := p.PipelineOpts.MountDirPath
	mountDir, err := i.BuildMountDir(c, mountDirPath)
	if err != nil {
		return nil, err
	}

	// 12. Target dir in dagger format.
	targetDirPath := p.PipelineOpts.TargetDirPath
	targetDir, err := i.BuildTargetDir(c, targetDirPath)
	if err != nil {
		return nil, err
	}

	var envVarsAllScanned map[string]string

	isAllEnvVarsToBeScanned := p.PipelineOpts.IsAllEnvVarsToScanEnabled
	if isAllEnvVarsToBeScanned {
		envVarsAllScanned, err = i.ScanAllEnvVars()
		if err != nil {
			return nil, err
		}
	}

	//targetDir := p.PipelineOpts.TargetDir
	//mountDir := p.PipelineOpts.MountDir
	//workDir := p.PipelineOpts.WorkDir
	return &Job{
		Id:                jobId,
		Name:              common.NormaliseStringUpper(new.Name),
		Stack:             common.NormaliseStringUpper(new.Stack),
		PipelineCfg:       p,
		Client:            c,
		ContainerImageURL: im,
		ContainerDefault:  ct,

		// Environment variables
		EnvVarsAWSScanned:        awsEnvVars,
		EnvVarsTerraformScanned:  terraformEnvVars,
		EnvVarsCustomScanned:     customEnvVars,
		EnvVarsToSet:             envVarsToSet,
		EnvVarsAllScanned:        envVarsAllScanned,
		EnvVarsFromDotEnvFile:    envVarsFromDotEnv,
		EnvVarsFromPrefixScanned: envVarsFromPrefix,

		// Directories (dagger format).
		RootDir:   rootDir,
		WorkDir:   workDir,
		MountDir:  mountDir,
		TargetDir: targetDir,

		// Paths
		RootDirPath:   ".",
		WorkDirPath:   new.WorkDir,
		MountDirPath:  p.PipelineOpts.MountDirPath,
		TargetDirPath: p.PipelineOpts.TargetDirPath,

		Ctx: new.PipelineCfg.Ctx,
	}, nil
}

// InitDagger 1. Init the job, initialising the Dagger client.
func (i *Instance) InitDagger() (*dagger.Client, error) {
	jobName := i.JobName
	jobId := i.JobId
	ux := i.InitOptions.PipelineCfg.UXMessage
	init := i.InitOptions

	ux.ShowInfo(uxPrefix,
		fmt.Sprintf("Initialising dagger client: job name: %s - job id: %s",
			jobName, jobId))

	c, err := daggerio.NewDaggerClient("", &init.PipelineCfg.Ctx,
		init.PipelineCfg.PipelineOpts.InitDaggerWithWorkDirByDefault)

	if err != nil {
		msg := GetErrMsg(jobName, jobId, "Dagger initialisation failed", nil)
		return nil, errors.NewDaggerEngineError(msg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(jobName, jobId,
		"Dagger client initialised"))

	return c, nil
}

func (i *Instance) ScanAllEnvVars() (map[string]string, error) {
	return filesystem.FetchAllEnvVarsFromHost()
}

// InitContainerImage 2. Get the container image.
func (i *Instance) InitContainerImage() (string, error) {
	stack := common.NormaliseStringUpper(i.InitOptions.Stack)
	init := i.InitOptions

	if stack == "" {
		errMsg := GetErrMsg(i.JobName, i.JobId, "Stack not specified, "+
			"cant create image with an empty stack", nil)
		return "", errors.NewDaggerEngineError(errMsg, nil)
	}

	image, err := daggerio.GetContainerImagePerStack(stack, "")
	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			fmt.Sprintf("Failed to get container image for stack: %s", stack), nil)
		return "", errors.NewDaggerEngineError(errMsg, err)
	}

	init.PipelineCfg.UXMessage.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId,
		fmt.Sprintf("Resolved Job container image: %s", image)))

	return image, nil
}

// InitContainer 3. Get the container.
func (i *Instance) InitContainer(c *dagger.Client, imageURL string) (*dagger.Container, error) {
	jobName := i.JobName
	jobId := i.JobId
	ux := i.InitOptions.PipelineCfg.UXMessage

	ux.ShowInfo(uxPrefix, GetInfoMsg(jobName, jobId,
		fmt.Sprintf("Initialising container with image: %s", imageURL)))

	container, err := daggerio.GetContainer(c, imageURL)

	if err != nil {
		errMsg := GetErrMsg(jobName, jobId,
			fmt.Sprintf("Failed to initialize container with image %s", imageURL), nil)
		return nil, errors.NewDaggerEngineError(errMsg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(jobName, jobId,
		"Container successfully initialised"))

	return container, nil
}

// ScanEnvVarsAWSKeys 4. Scan (if applicable) AWS keys environment variables.
func (i *Instance) ScanEnvVarsAWSKeys(scanAWSVars bool) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	if !scanAWSVars {
		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Skipping AWS env var scan"))
		return map[string]string{}, nil
	}

	envVars, err := filesystem.ScanAWSCredentialsEnvVars()
	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to scan AWS env vars", nil)
		return nil, errors.NewDaggerEngineError(errMsg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "AWS env vars scanned successfully"))

	return envVars, nil
}

// ScanEnvVarsTerraform 5. Scan (if applicable) Terraform environment variables.
func (i *Instance) ScanEnvVarsTerraform(scanAWSVars bool) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	if !scanAWSVars {
		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Skipping Terraform env var scan"))
		return map[string]string{}, nil
	}

	envVars, err := filesystem.ScanTerraformEnvVars()
	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to scan Terraform env vars", nil)
		return nil, errors.NewDaggerEngineError(errMsg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Terraform env vars scanned successfully"))

	return envVars, nil
}

// ScanEnvVarsCustom 6. Scan (if applicable) custom environment variables.
func (i *Instance) ScanEnvVarsCustom(scanCustomVars []string) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	if len(scanCustomVars) == 0 {
		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Skipping custom env var scan"))
		return map[string]string{}, nil
	}

	envVars, err := filesystem.FetchEnvVarsAsMap(scanCustomVars, []string{})
	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to scan custom env vars", nil)
		return nil, errors.NewDaggerEngineError(errMsg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Custom env vars scanned successfully"))

	return envVars, nil
}

func (i *Instance) ScanEnvVarsFromDotEnvFile(dotEnvFile string) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	if i.InitOptions.IsScanEnvVarsFromDotEnv {
		envVars, err := filesystem.GetEnvVarsFromDotFile(dotEnvFile)
		if err != nil {
			errMsg := GetErrMsg(i.JobName, i.JobId,
				"Failed to scan env vars from .env file", nil)
			return nil, errors.NewDaggerEngineError(errMsg, err)
		}

		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Env vars scanned successfully from .env file"))

		return envVars, nil
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Skipping env var scan from .env file"))
	return map[string]string{}, nil
}

func (i *Instance) ScanEnvVarsFromPrefix(prefixes []string) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	if i.InitOptions.IsScanEnvVarsFromPrefix {
		envVars, err := filesystem.ScanEnvVarsFromPrefixes(prefixes)
		if err != nil {
			errMsg := GetErrMsg(i.JobName, i.JobId,
				"Failed to scan env vars from prefix", nil)
			return nil, errors.NewDaggerEngineError(errMsg, err)
		}

		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Env vars scanned successfully from prefix"))

		return envVars, nil
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Skipping env var scan from prefix"))
	return map[string]string{}, nil
}

// ValidatedEnvVarsPassed 7. Validate environment variables to be set.
func (i *Instance) ValidatedEnvVarsPassed(envVarsToSet map[string]string) (map[string]string, error) {
	ux := i.InitOptions.PipelineCfg.UXMessage

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Validating env vars to be set"))

	if common.MapIsNulOrEmpty(envVarsToSet) {
		ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "No env vars to be set"))
		return map[string]string{}, nil
	}

	if err := filesystem.AreEnvVarsConsistent(envVarsToSet); err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to validate env vars to be set", nil)
		return map[string]string{}, errors.NewDaggerEngineError(errMsg, err)
	}

	ux.ShowInfo(uxPrefix, GetInfoMsg(i.JobName, i.JobId, "Env vars to be set validated successfully"))

	return envVarsToSet, nil
}

// BuildRootDir 8. Build root directory.
func (i *Instance) BuildRootDir(client *dagger.Client) (*dagger.Directory, error) {
	dir, err := daggerio.GetDaggerDir(client, "")

	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to get dagger root directory", nil)
		return nil, errors.NewDaggerConfigurationError(errMsg, err)
	}

	return dir, nil
}

func (i *Instance) BuildWorkDir(client *dagger.Client, workDir string) (*dagger.Directory, error) {
	dir, err := daggerio.GetDaggerDirWithEntriesCheck(client, workDir)

	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to get dagger working directory", nil)
		return nil, errors.NewDaggerConfigurationError(errMsg, err)
	}

	return dir, nil
}

func (i *Instance) BuildMountDir(client *dagger.Client, mountDir string) (*dagger.Directory,
	error) {
	dir, err := daggerio.GetDaggerDirWithEntriesCheck(client, mountDir)

	if err != nil {
		errMsg := GetErrMsg(i.JobName, i.JobId,
			"Failed to get dagger working directory", nil)
		return nil, errors.NewDaggerConfigurationError(errMsg, err)
	}

	return dir, nil
}

func (i *Instance) BuildTargetDir(client *dagger.Client, targetDir string) (*dagger.Directory,
	error) {
	// FIXME: Check whether this function will be actually required, or it can be deprecated.
	return nil, nil
}
