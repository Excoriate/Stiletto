package task

import (
	"context"
	"fmt"
	"github.com/Excoriate/stiletto/internal/cloud/awscloud"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
)

type AWSECRPushAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
	Ctx    context.Context
}

type AWSECRPushActionOptions struct {
	AWSAccessKey              string
	AWSSecretKey              string
	AWSRegion                 string
	ScanAWSCredentialsFromEnv bool
}

type AWSECRPushActionArgs struct {
	AWSRegion         string
	AWSAccessKey      string
	AWSSecretKey      string
	Repository        string
	Registry          string
	Tag               string
	RunECRLoginInHost bool
}

type AWSECRPushActions interface {
	DeployNewTask() (Output, error)
}

func getBuildTagAndPushActionArgs(uxLog tui.TUIMessenger) (AWSECRPushActionArgs, error) {
	awsCredentialsCfg, err := awscloud.GetCredentials()

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"AWS credentials could not be met")
		uxLog.ShowError("AWS:ECR:PUSH", errMsg, err)
		return AWSECRPushActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	cfg := config.Cfg{}
	registry, err := cfg.GetFromAny("ecr-registry")

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"ECR repository could not be met")
		uxLog.ShowError("AWS:ECR:PUSH", errMsg, err)
		return AWSECRPushActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	tag, err := cfg.GetFromAny("tag")
	if err != nil {
		warnMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"tag could not be met. 'Latest' will be used.")
		uxLog.ShowWarning("AWS:ECR:PUSH", warnMsg)

		tag.Value = "latest"
	}

	repository, err := cfg.GetFromAny("ecr-repository")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"ECR repository could not be met")
		uxLog.ShowError("AWS:ECR:PUSH", errMsg, err)
		return AWSECRPushActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	runInVendor := cfg.IsRunningInVendorAutomation()
	if runInVendor {
		uxLog.ShowWarning("AWS:ECR:PUSH", "Running in vendor automation. "+
			"ECR Login should be performed using the 'automation' mechanism that the vendor"+
			" provides (E.g.: GitHub action)")
	}

	return AWSECRPushActionArgs{
		AWSRegion:         awsCredentialsCfg.Region,
		AWSAccessKey:      awsCredentialsCfg.AccessKeyID,
		AWSSecretKey:      awsCredentialsCfg.SecretAccessKey,
		Repository:        repository.Value.(string),
		Registry:          registry.Value.(string),
		Tag:               tag.Value.(string),
		RunECRLoginInHost: !runInVendor,
	}, nil
}

func (a *AWSECRPushAction) DeployNewTask() (Output, error) {
	// Getting all the requirements.
	uxLog := a.Task.GetPipelineUXLog()
	opts, err := getBuildTagAndPushActionArgs(uxLog)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Specific container/runtime requirements.
	ctx := a.Task.GetJob().Ctx
	container := a.Task.GetJobContainerDefault()
	client := a.Task.GetClient()
	targetDir := a.Task.GetCoreTask().Dirs.TargetDir
	preRequiredFiles := []string{"Dockerfile"}

	// Mounting dir.
	containerToUse, err := a.Task.MountDir(targetDir, client, container, preRequiredFiles, ctx)
	if err != nil {
		return Output{}, err
	}

	// Resolving publish address, repository URL, and other parameters to pass to AWS ECR.
	repositoryURL := awscloud.GetImageURL(opts.Repository, opts.Tag)
	publishAddress := awscloud.GetECRPublishAddress(opts.Registry, repositoryURL)
	uxLog.ShowInfo(a.prefix, fmt.Sprintf("Pushing image to %s", publishAddress))

	// Logging into AWS ECR.
	if opts.RunECRLoginInHost {
		err = awscloud.AWSECRLogin(opts.Registry, awscloud.AWSCredentials{
			AccessKeyID:     opts.AWSAccessKey,
			SecretAccessKey: opts.AWSSecretKey,
			Region:          opts.AWSRegion,
		})

		if err != nil {
			errMsg := fmt.Sprintf("Failed to login to AWS ECR")
			uxLog.ShowError(a.prefix, errMsg, err)
			return Output{}, errors.NewActionCfgError(errMsg, err)
		}
	}

	// Publishing the image into ECR.
	dockerFileDir, _ := a.Task.ConvertDir(client, targetDir)
	err = a.Task.PushImage(publishAddress, containerToUse, dockerFileDir, ctx)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to push image to AWS ECR: %s", publishAddress)
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	uxLog.ShowSuccess(a.prefix, fmt.Sprintf("Pushed image to %s", publishAddress))

	return Output{}, nil
}

func NewAWSECRAction(task CoreTasker, prefix string) AWSECRPushActions {
	return &AWSECRPushAction{
		Task:   task,
		prefix: prefix,
		Id:     common.GetUUID(),
		Name:   "Build, tag and push Docker Image to AWS ECR",
	}
}
