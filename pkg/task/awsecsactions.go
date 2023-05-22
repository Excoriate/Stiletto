package task

import (
	"context"
	"fmt"
	"github.com/Excoriate/stiletto/internal/cloud/adapters/clients"
	"github.com/Excoriate/stiletto/internal/cloud/awscloud"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
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

	EnvVarsToSetInContainerDef map[string]string
}

type AWSECSDeployActions interface {
	DeployTask() (Output, error)
}

func getDeployActionArgs(log tui.TUIMessenger) (AWSECSDeployActionArgs, error) {
	awsCredentialsCfg, err := awscloud.GetCredentials()
	actionPrefix := "AWS:ECS:DEPLOY"

	if err != nil {
		msg := "Failed to execute ECS action. Pre-requirements could not be satisfied"
		log.ShowError(actionPrefix, msg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(msg, err)
	}

	cfg := config.Cfg{}

	ecsService, err := cfg.GetFromAny("ecs-service")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'ecs-service' could not be met")
		log.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	ecsCluster, err := cfg.GetFromAny("ecs-cluster")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'ecs-service' could not be met")
		log.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	ecsTaskDefName, err := cfg.GetFromAny("task-definition")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get 'ecsDeployAction' arguments, " +
			"'ecs-service' could not be met")
		log.ShowError(actionPrefix, errMsg, err)
		return AWSECSDeployActionArgs{}, errors.NewActionCfgError(errMsg, err)
	}

	imageUrl, err := cfg.GetFromAny("image-url")
	if err != nil {
		log.ShowWarning(actionPrefix, "No 'image-url' found, "+
			"using 'container-image' field that's set in the task definition as the default value")
		imageUrl.Value = "use-task-def"
	}

	tag, err := cfg.GetFromAny("release-version")
	if err != nil {
		log.ShowWarning(actionPrefix, "No 'image-tag' found, the value 'latest' will be used")
		tag.Value = "latest"
	}

	// Env specific options.
	var contDefVarsTotal map[string]string
	// Allowed scan/set options to update the task definition/container def.
	//with environment variables.
	var contDefEnvVarsScannedFromHost map[string]string
	var contDefEnvVarsScannedFromKeys map[string]string
	var contDefEnvVarsSetCustom map[string]string
	var contDefEnvVarScannedByPrefix map[string]string

	// 1. Scan from host.
	envVarsScanFromHostCfg, err := cfg.GetBoolFromViper("set-env-from-host")
	if err != nil {
		log.ShowInfo(actionPrefix, "No 'set-env-from-host' found, "+
			"no environment variables will be scanned from host")
	} else {
		isEnvVarsScanFromHostEnabled := envVarsScanFromHostCfg.Value.(bool)
		if isEnvVarsScanFromHostEnabled {
			log.ShowWarning(actionPrefix, "The 'set-env-from-host' is set. "+
				"All the host environment variables will be scanned and set in the task definition/container def.")

			contDefEnvVarsScannedFromHost, err = filesystem.FetchAllEnvVarsFromHost()
			if err != nil {
				log.ShowError(actionPrefix, "Failed to scan the host environment variables", err)
				return AWSECSDeployActionArgs{}, errors.NewActionCfgError("Failed to scan the host environment variables", err)
			}
		} else {
			log.ShowInfo(actionPrefix, "The option 'set-env-from-host' is disabled, "+
				"no environment variables will be scanned from host")
		}
	}

	// 2. Scan from specific keys passed.
	envVarsScanFromKeysCfg, err := cfg.GetStringSliceFromViper("set-env-from-keys")
	if err != nil {
		log.ShowInfo(actionPrefix, "No 'set-env-from-keys' found, "+
			"no environment variables will be scanned from keys")
	} else {
		keysToScan := envVarsScanFromKeysCfg.Value.([]string)
		if len(keysToScan) == 0 {
			log.ShowInfo(actionPrefix, "The option 'set-env-from-keys' is set, "+
				"however no keys were passed, no environment variables will be scanned from keys")
		} else {
			scannedFromKeys, err := filesystem.FetchEnvVarsAsMap(keysToScan, []string{})
			if err != nil {
				log.ShowError(actionPrefix, "Failed to scan the environment variables from keys", err)
				return AWSECSDeployActionArgs{}, errors.NewActionCfgError("Failed to scan the environment variables from keys", err)
			}

			contDefEnvVarsScannedFromKeys = scannedFromKeys
		}
	}

	// 3. Scan from prefix.
	envVarsScanFromPrefixCfg, err := cfg.GetStringFromViper("set-env-vars-with-prefix")
	if err != nil {
		log.ShowInfo(actionPrefix, "The option 'set-env-vars-with-prefix' is not set, no environment variables will be scanned from prefix")
	} else {
		prefix := envVarsScanFromPrefixCfg.Value.(string)
		prefix = common.NormaliseNoSpaces(prefix)

		if prefix == "" {
			log.ShowInfo(actionPrefix, "The option 'set-env-vars-with-prefix' is set, "+
				"however the prefix is empty")
		} else {
			log.ShowInfo(actionPrefix, fmt.Sprintf("Scanning the environment variables with the prefix '%s'", prefix))
			scannedEnvVarsWithPrefix, err := filesystem.FetchEnvVarsWithPrefix(prefix)
			if err != nil {
				log.ShowError(actionPrefix, "Failed to scan the environment variables with the prefix", err)
				return AWSECSDeployActionArgs{}, errors.NewActionCfgError("Failed to scan the environment variables with the prefix", err)
			}

			contDefEnvVarScannedByPrefix = scannedEnvVarsWithPrefix
		}
	}

	// 4. Set custom and directly passed environment variables
	envVarsSetCustomCfg, err := cfg.GetStringMapFromViper("set-env-vars-custom")
	if err != nil {
		log.ShowInfo(actionPrefix, "No 'set-env-vars-custom' found, so no custom environment variables will be set")
	} else {
		envVarsSetCustomValue := envVarsSetCustomCfg.Value
		if len(envVarsSetCustomValue.(map[string]interface{})) == 0 {
			log.ShowInfo(actionPrefix, "The option 'set-env-vars-custom' is set, "+
				"however no custom environment variables were passed, no environment variables will be set")
		} else {
			log.ShowInfo(actionPrefix, "Setting the custom environment variables")
			contDefEnvVarsSetCustom = envVarsSetCustomValue.(map[string]string)
		}
	}

	contDefVarsTotal = filesystem.MergeEnvVars(contDefEnvVarsScannedFromHost,
		contDefEnvVarsScannedFromKeys, contDefEnvVarScannedByPrefix, contDefEnvVarsSetCustom)

	return AWSECSDeployActionArgs{
		AWSRegion:                  awsCredentialsCfg.Region,
		AWSAccessKey:               awsCredentialsCfg.AccessKeyID,
		AWSSecretKey:               awsCredentialsCfg.SecretAccessKey,
		ClusterName:                ecsCluster.Value.(string),
		ServiceName:                ecsService.Value.(string),
		TaskDefinition:             ecsTaskDefName.Value.(string),
		ImageTagOrReleaseVersion:   tag.Value.(string),
		Image:                      imageUrl.Value.(string),
		EnvVarsToSetInContainerDef: contDefVarsTotal,
	}, nil

}

func (a *AWSECSDeployAction) DeployTask() (Output, error) {
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
			ImageURL:             opts.Image,
			Version:              opts.ImageTagOrReleaseVersion,
			EnvironmentVariables: opts.EnvVarsToSetInContainerDef,
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
