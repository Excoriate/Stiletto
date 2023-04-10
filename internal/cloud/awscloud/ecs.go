package awscloud

import (
	"context"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSTaskDefContainerDefUpdateOptions struct {
	ImageURL string
	Version  string
}

type ECSUpdateServiceOptions struct {
	Service            string
	Cluster            string
	ForceNewDeployment bool
	TaskDefARN         string
}

func GetECSTaskDefinition(client *ecs.Client, taskDefName string) (*ecs.DescribeTaskDefinitionOutput, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefName,
	}

	taskDef, err := client.DescribeTaskDefinition(context.TODO(), input)
	if err != nil {
		return &ecs.DescribeTaskDefinitionOutput{}, err
	}

	return taskDef, nil
}

func UpdateECSTaskContainerDefinition(client *ecs.Client,
	taskDef *ecs.DescribeTaskDefinitionOutput, opt ECSTaskDefContainerDefUpdateOptions) (string,
	error) {
	imageURL := common.NormaliseNoSpaces(opt.ImageURL)
	version := common.NormaliseNoSpaces(opt.Version)

	if !common.IsImageURLIncludesTag(imageURL) && version != "" {
		imageURL = fmt.Sprintf("%s:%s", imageURL, version)
	} else {
		imageURL = fmt.Sprintf("%s:latest", imageURL)
	}

	for i := range taskDef.TaskDefinition.ContainerDefinitions {
		taskDef.TaskDefinition.ContainerDefinitions[i].Image = &imageURL
	}

	newTaskDefInput := &ecs.RegisterTaskDefinitionInput{
		Family:                  taskDef.TaskDefinition.Family,
		TaskRoleArn:             taskDef.TaskDefinition.TaskRoleArn,
		ExecutionRoleArn:        taskDef.TaskDefinition.ExecutionRoleArn,
		NetworkMode:             taskDef.TaskDefinition.NetworkMode,
		ContainerDefinitions:    taskDef.TaskDefinition.ContainerDefinitions, // updated value.
		RequiresCompatibilities: taskDef.TaskDefinition.RequiresCompatibilities,
		Cpu:                     taskDef.TaskDefinition.Cpu,
		Memory:                  taskDef.TaskDefinition.Memory}

	updatedTask, err := client.RegisterTaskDefinition(context.TODO(), newTaskDefInput)
	if err != nil {
		return "", err
	}

	// TaskARN that'll be used to perform an 'updateECSService' sort of operation,
	//in a new deployment action.

	updatedTaskDefARN := *updatedTask.TaskDefinition.TaskDefinitionArn

	return updatedTaskDefARN, nil
}

func UpdateECSService(client *ecs.Client, opt ECSUpdateServiceOptions) error {
	input := &ecs.UpdateServiceInput{
		Cluster:            aws.String(opt.Cluster),
		Service:            aws.String(opt.Service),
		TaskDefinition:     aws.String(opt.TaskDefARN),
		ForceNewDeployment: opt.ForceNewDeployment,
	}

	_, err := client.UpdateService(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
