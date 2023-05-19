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
	GenerateRandomTag bool
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
			"tag could not be met. 'Latest' will be used if the --generate-random-tag option is" +
			" not set.")
		uxLog.ShowWarning("AWS:ECR:PUSH", warnMsg)

		tag.Value = ""
	}

	tagToSet := tag.Value.(string)

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

	generateRandomTag, err := cfg.GetFromViperOrDefault("generate-random-tag", false)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"generate-random-tag could not be met")
		uxLog.ShowError("AWS:ECR:PUSH", errMsg, err)
		return AWSECRPushActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	generateRandomTagValue := generateRandomTag.Value.(bool)

	if generateRandomTagValue && tagToSet != "" {
		errMsg := fmt.Sprintf("Failed to get 'buildTagAndPush' arguments, " +
			"generate-random-tag and tag cannot be used together")
		uxLog.ShowError("AWS:ECR:PUSH", errMsg, err)
		return AWSECRPushActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	// If tag is not set and generate-random-tag is not set, default to 'latest'
	if tagToSet == "" && !generateRandomTagValue {
		tagToSet = "latest"
	}

	// If generate-random-tag is set, and tag is not set, generate a random tag.
	if generateRandomTagValue && tagToSet == "" {
		tagToSet = common.GenerateRandomString(5, true)
	}

	return AWSECRPushActionArgs{
		AWSRegion:         awsCredentialsCfg.Region,
		AWSAccessKey:      awsCredentialsCfg.AccessKeyID,
		AWSSecretKey:      awsCredentialsCfg.SecretAccessKey,
		Repository:        repository.Value.(string),
		Registry:          registry.Value.(string),
		Tag:               tagToSet,
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
	targetDir := a.Task.GetJob().TargetDirPath
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
