package task

import (
	"context"
	"fmt"
	"github.com/Excoriate/stiletto/internal/cloud/adapters/clients"
	"github.com/Excoriate/stiletto/internal/cloud/awscloud"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
)

type AWSECSDeployAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
	Ctx    context.Context
}

type AWSECSDeployActionOptions struct {
	AWSAccessKey              string
	AWSSecretKey              string
	AWSRegion                 string
	ScanAWSCredentialsFromEnv bool
}

type AWSECSDeployActionArgs struct {
	AWSRegion                string
	AWSAccessKey             string
	AWSSecretKey             string
	ClusterName              string
	ServiceName              string
	TaskDefinition           string
	ImageTagOrReleaseVersion string
	Image                    string
}

type AWSECSDeployActions interface {
	DeployNewTask() (Output, error)
}

func getDeployActionArgs(uxLog tui.TUIMessenger) (AWSECSDeployActionArgs, error) {
	awsCredentialsCfg, err := awscloud.GetCredentials()
	actionPrefix := "AWS:ECS:DEPLOY"

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"AWS credentials could not be met")
		uxLog.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	cfg := config.Cfg{}

	ecsService, err := cfg.GetFromAny("ecs-service")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'ecs-service' could not be met")
		uxLog.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	ecsCluster, err := cfg.GetFromAny("ecs-cluster")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'ecs-cluster' could not be met")
		uxLog.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	ecsTaskDefName, err := cfg.GetFromAny("task-definition")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'task-definition' could not be met")
		uxLog.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	imageUrl, err := cfg.GetFromAny("image-url")
	if err != nil {
		uxLog.ShowWarning(actionPrefix, "No 'image-url' found, "+
			"using 'container-image' field that's set in the task definition as the default value")
		imageUrl.Value = "use-task-def"
	}

	tag, err := cfg.GetFromAny("release-version")
	if err != nil {
		uxLog.ShowWarning(actionPrefix, "No 'image-tag' found, the value 'latest' will be used")
		tag.Value = "latest"
	}

	return AWSECSDeployActionArgs{
		AWSRegion:                awsCredentialsCfg.Region,
		AWSAccessKey:             awsCredentialsCfg.AccessKeyID,
		AWSSecretKey:             awsCredentialsCfg.SecretAccessKey,
		ClusterName:              ecsCluster.Value.(string),
		ServiceName:              ecsService.Value.(string),
		TaskDefinition:           ecsTaskDefName.Value.(string),
		ImageTagOrReleaseVersion: tag.Value.(string),
		Image:                    imageUrl.Value.(string),
	}, nil

}

func (a *AWSECSDeployAction) DeployNewTask() (Output, error) {
	// Getting all the requirements.
	uxLog := a.Task.GetPipelineUXLog()
	opts, err := getDeployActionArgs(uxLog)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Getting the AWS Client, to perform the actual deployment.
	ecsClient, err := clients.GetAWSECSClient(opts.AWSRegion)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get AWS ECS client")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Get the target task definition.
	taskDef, err := awscloud.GetECSTaskDefinition(ecsClient, opts.TaskDefinition)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get AWS ECS task definition")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Update the task definition.
	updateTaskARN, err := awscloud.UpdateECSTaskContainerDefinition(ecsClient, taskDef,
		awscloud.ECSTaskDefContainerDefUpdateOptions{
			ImageURL: opts.Image,
			Version:  opts.ImageTagOrReleaseVersion,
		})

	if err != nil {
		errMsg := fmt.Sprintf("Failed to update AWS ECS task definition")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Update the service, and perform the actual deployment.
	err = awscloud.UpdateECSService(ecsClient, awscloud.ECSUpdateServiceOptions{
		Cluster:            opts.ClusterName,
		Service:            opts.ServiceName,
		TaskDefARN:         updateTaskARN,
		ForceNewDeployment: true,
	})

	if err != nil {
		errMsg := fmt.Sprintf("Failed to update AWS ECS service")
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	uxLog.ShowSuccess(a.prefix,
		fmt.Sprintf("Successfully deployed new task to AWS ECS service '%s' - Task definition"+
			" ARN deployed %s", opts.ServiceName, updateTaskARN))

	return Output{}, nil
}

func NewAWSECSAction(task CoreTasker, prefix string) AWSECSDeployActions {
	return &AWSECSDeployAction{
		Task:   task,
		prefix: prefix,
		Id:     common.GetUUID(),
		Name:   "Deploy, or manage ECS configurations, tasks and others",
	}
}
