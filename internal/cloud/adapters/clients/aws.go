package clients

import (
	"context"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func GetAWS(region string) (aws.Config, error) {
	uxLog := tui.NewTUIMessage()
	awsAuth, err := awsCfg.LoadDefaultConfig(context.TODO(), awsCfg.WithRegion(region))
	if err != nil {
		uxLog.ShowError("AWS", "Failed to get AWS credentials. Cannot initialise AWS SDK", err)
		return aws.Config{}, err
	}

	return awsAuth, nil
}

func GetAWSECSClient(region string) (*ecs.Client, error) {
	awsAuth, err := GetAWS(region)
	if err != nil {
		return nil, err
	}

	return ecs.NewFromConfig(awsAuth), nil
}
