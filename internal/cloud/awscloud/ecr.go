package awscloud

import (
	"bytes"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/tui"
	"os/exec"
)

func GetImageURL(repository, tag string) string {
	repoNormalised := common.NormaliseNoSpaces(repository)
	tagNormalised := common.NormaliseNoSpaces(tag)

	if tagNormalised == "" {
		tagNormalised = "latest"
	}

	return fmt.Sprintf("%s:%s", repoNormalised, tagNormalised)
}

func GetECRPublishAddress(registry, repository string) string {
	repoNormalised := common.NormaliseNoSpaces(repository)
	registryNormalised := common.NormaliseNoSpaces(registry)

	return fmt.Sprintf("%s/%s", registryNormalised, repoNormalised)
}

func AWSECRLogin(registry string, credentials AWSCredentials) error {
	var out bytes.Buffer
	uxLog := tui.NewTUIMessage()
	uxLogPrefix := "ECR-LOGIN"

	cmd := exec.Command("aws", "ecr", "get-login-password", "--region", credentials.Region)
	cmd.Stdout = &out

	uxLog.ShowInfo(uxLogPrefix, fmt.Sprintf("Getting ECR login password for %s", registry))

	if err := cmd.Run(); err != nil {
		errMsg := fmt.Sprintf("Failed to get ECR login password, with region %s", credentials.Region)
		uxLog.ShowError(uxLogPrefix, errMsg, err)
		return errors.NewInternalPipelineError(errMsg)
	}

	token := out.String()

	out = bytes.Buffer{}
	cmd = exec.Command("docker", "login", "--username", "AWS", "--password", token, registry)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		errMsg := fmt.Sprintf("Failed to login to ECR, with region %s", credentials.Region)
		uxLog.ShowError(uxLogPrefix, errMsg, err)
		return errors.NewInternalPipelineError(errMsg)
	}

	return nil
}
