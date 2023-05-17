package task

import (
	"context"
	"github.com/Excoriate/stiletto/internal/common"
)

type DockerBuildAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
	Ctx    context.Context
}

type DockerBuildActions interface {
	BuildTagAndPush(dockerFile string) (Output, error)
}

func (a *DockerBuildAction) BuildTagAndPush(dockerFile string) (Output, error) {
	ctx := a.Task.GetJob().Ctx

	container := a.Task.GetJobContainerDefault()
	client := a.Task.GetClient()
	targetDir := a.Task.GetCoreTask().Dirs.TargetDir
	preRequiredFiles := []string{"Dockerfile"}

	mountedContainer, err := a.Task.MountDir(targetDir, client, container, preRequiredFiles, ctx)
	if err != nil {
		return Output{}, err
	}

	targetDirDagger, _ := a.Task.ConvertDir(client, targetDir)

	_, err = mountedContainer.
		WithExec([]string{"ls", "-ltrh"}).
		WithExec([]string{"cat", "Dockerfile"}).
		Build(targetDirDagger).
		ExitCode(ctx)

	return Output{}, nil
}

func NewDockerAction(task CoreTasker) DockerBuildActions {
	return &DockerBuildAction{
		Task:   task,
		prefix: "DOCKER:ACTION-BUILD",
		Id:     common.GetUUID(),
		Name:   "Build Docker Image",
	}
}
