package daggerio

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/logger"
	"strings"
)

//var PlatformToArch = map[dagger.Platform]string{
//	"linux/amd64": "amd64",
//	"linux/arm64": "arm64",
//}

// GetContainerImageCustom DaggerContainerImage represents the container image of the dagger client.
func GetContainerImageCustom(imageURL, version string) (DaggerContainerImage, error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()
	var imageNormalised string
	var versionNormalised string

	if imageURL == "" {
		return DaggerContainerImage{}, errors.NewDaggerEngineError("Unable to fetch container image, "+
			"image URL value is empty",
			nil)
	}

	// If the image has ":", it means that the user has passed the version as well.
	if strings.Contains(imageURL, ":") {
		logPrinter.LogWarn("Dagger Image Configuration", "It seems that you have passed the"+
			" version of the image as"+
			" well along"+
			" with the image name. The version will be ignored and the version passed will be used.")

		imageNormalised = common.NormaliseStringLower(strings.Split(imageURL, ":")[0])
		versionNormalised = common.NormaliseStringLower(strings.Split(imageURL, ":")[1])
	} else {
		imageNormalised = imageURL

		if version == "" {
			logPrinter.LogWarn("Dagger Image Configuration",
				fmt.Sprintf("The 'version' option is empty, "+
					"it will use 'latest' as the image version for the image %s passed", imageURL))
			versionNormalised = "latest"
		} else {
			versionNormalised = version
		}
	}

	return DaggerContainerImage{
		Image:   imageNormalised,
		Version: versionNormalised,
	}, nil
}

// GetContainerImagePerStack returns the container image of the dagger client.
func GetContainerImagePerStack(stack string, version string) (string, error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()

	if stack == "" {
		return "", errors.NewDaggerEngineError("Unable to fetch container image, "+
			"stack value is empty",
			nil)
	}

	stackNormalised := common.NormaliseStringUpper(stack)

	if version == "" {
		logPrinter.LogWarn("Dagger Image Configuration",
			fmt.Sprintf("The 'version' option is empty, "+
				"it will use 'latest' as the image version for the image %s passed", stackNormalised))
		version = "latest"
	}

	if _, ok := StackImagesMap[stackNormalised]; !ok {
		return "", errors.NewDaggerEngineError(fmt.Sprintf("Unable to fetch container image, "+
			"stack %s is not supported or it could not be found", stackNormalised),
			nil)
	}

	imageFromStackConfig := StackImagesMap[stackNormalised]
	// If the image has ":", it means that the user has passed the version as well.
	if strings.Contains(imageFromStackConfig, ":") {
		logPrinter.LogWarn("Dagger Image Configuration", "It seems that you have passed the"+
			" version of the image as"+
			" well along"+
			" with the image name. The version will be ignored and the version passed will be used.")

		imageURLFromDefault := common.NormaliseStringLower(strings.Split(imageFromStackConfig, ":")[0])
		imageVersionFromDefault := common.NormaliseStringLower(strings.Split(imageFromStackConfig, ":")[1])

		return fmt.Sprintf("%s:%s", imageURLFromDefault, imageVersionFromDefault), nil
	}

	return fmt.Sprintf("%s:%s", StackImagesMap[stackNormalised], version), nil
}

// GetContainer returns the container of the dagger client.
func GetContainer(c *dagger.Client, image string) (*dagger.Container, error) {
	if image == "" {
		return nil, errors.NewDaggerEngineError("Unable to fetch container, image value is empty", nil)
	}

	if c == nil {
		return nil, errors.NewDaggerEngineError("Unable to fetch container, dagger client is nil", nil)
	}

	return c.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).From(common.
		NormaliseStringLower(image)), nil
}

// NormaliseDaggerPath will check if the path includes a / at the beginning; if so, just return it; if not, add it.
func NormaliseDaggerPath(path string) string {
	if path == "" {
		return ""
	}

	// Check if the path starts with /build
	if strings.HasPrefix(path, "/build") {
		return path
	}

	return fmt.Sprintf("/build/%s", path)
}

// PushImage pushes the image to the dagger client.
func PushImage(container *dagger.Container, url string, ctx context.Context) (string, error) {
	if container == nil {
		return "", errors.NewDaggerEngineError("Unable to push image, container is nil", nil)
	}

	if url == "" {
		return "", errors.NewDaggerEngineError("Unable to push image, URL is empty", nil)
	}

	url = common.NormaliseStringLower(url)

	addr, err := container.Publish(ctx, url)
	if err != nil {
		return "", errors.NewDaggerEngineError("Unable to push image", err)
	}

	return addr, nil
}

// BuildImage builds the image of the dagger client.
func BuildImage(dockerFilePath string, client *dagger.Client,
	container *dagger.Container) (*dagger.Container, error) {
	if container == nil {
		return nil, errors.NewDaggerEngineError("Unable to build image, container is nil", nil)
	}

	if dockerFilePath == "" {
		return nil, errors.NewDaggerEngineError("Unable to build image, "+
			"docker file directory is empty", nil)
	}

	dockerFileDir, err := GetDaggerDirWithEntriesCheck(client, dockerFilePath)
	if err != nil {
		return nil, errors.NewDaggerEngineError("Unable to build image", err)
	}

	return container.Build(dockerFileDir), nil
}
