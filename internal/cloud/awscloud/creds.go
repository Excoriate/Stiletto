package awscloud

import (
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/pkg/config"
)

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func GetAWSRegionSet() (string, error) {
	cfg := config.Cfg{}
	awsRegion, viperErr := cfg.GetFromViper("aws-region")
	awsRegionEnvVars, envVarErr := cfg.GetFromEnvVars("AWS_REGION")

	if viperErr != nil && envVarErr != nil {
		return "", errors.NewAWSCfgError("AWS_REGION is not set ("+
			"check the flags passed or exported env vars)", nil)
	}

	if viperErr == nil {
		return awsRegion.Value.(string), nil
	}

	return awsRegionEnvVars.Value.(string), nil
}

func GetAWSAccessKeyID() (string, error) {
	cfg := config.Cfg{}
	awsAccessKeyID, viperErr := cfg.GetFromViper("aws-access-key-id")
	awsAccessKeyIDEnvVars, envVarErr := cfg.GetFromEnvVars("AWS_ACCESS_KEY_ID")

	if viperErr != nil && envVarErr != nil {
		return "", errors.NewAWSCfgError("AWS_ACCESS_KEY_ID is not set ("+
			"check the flags passed or exported env vars)", nil)
	}

	if viperErr == nil {
		return awsAccessKeyID.Value.(string), nil
	}

	return awsAccessKeyIDEnvVars.Value.(string), nil

}

func GetAWSSecretAccessKey() (string, error) {
	cfg := config.Cfg{}
	secretAccessKey, viperErr := cfg.GetFromViper("aws-secret-access-key")
	secretAccessKeyEnvVars, envVarErr := cfg.GetFromEnvVars("AWS_SECRET_ACCESS_KEY")

	if viperErr != nil && envVarErr != nil {
		return "", errors.NewAWSCfgError("AWS_SECRET_ACCESS_KEY is not set ("+
			"check the flags passed or exported env vars)", nil)
	}

	if viperErr == nil {
		return secretAccessKey.Value.(string), nil
	}

	return secretAccessKeyEnvVars.Value.(string), nil
}

func GetCredentials() (AWSCredentials, error) {
	awsRegion, err := GetAWSRegionSet()
	if err != nil {
		return AWSCredentials{}, err
	}

	awsAccessKeyID, err := GetAWSAccessKeyID()
	if err != nil {
		return AWSCredentials{}, err
	}

	awsSecretAccessKey, err := GetAWSSecretAccessKey()
	if err != nil {
		return AWSCredentials{}, err
	}

	return AWSCredentials{
		AccessKeyID:     awsAccessKeyID,
		SecretAccessKey: awsSecretAccessKey,
		Region:          awsRegion,
	}, nil
}

func GetCredentialsAsEnvVarsMap(cred AWSCredentials) map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     cred.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY": cred.SecretAccessKey,
		"AWS_REGION":            cred.Region,
	}
}
