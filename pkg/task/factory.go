package task

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/pkg/job"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

func NewTask(p *pipeline.Config, job *job.Job, actions []string,
	init *InitOptions) *Task {
	taskId := common.GetUUID()
	taskName := "BUILD"
	stackName := "DOCKER"

	p.UXMessage.ShowInfo("TASK-INIT", fmt.Sprintf("Initialising task: %s with id: %s",
		taskName, taskId))

	envVars := filesystem.MergeEnvVars(job.EnvVarsToSet, job.EnvVarsAWSScanned,
		job.EnvVarsCustomScanned, job.EnvVarsTerraformScanned)

	if common.MapIsNulOrEmpty(envVars) {
		infoMsg := fmt.Sprintf("No environment variables are passed from the Job instance, "+
			"skipping the environment variable configuration step for task name: %s with id: %s", taskName, taskId)
		p.UXMessage.ShowInfo("TASK-INIT", infoMsg)
	} else {
		p.UXMessage.ShowInfo("TASK-INIT",
			fmt.Sprintf("Environment variables are passed from the Job instance, "+
				"setting the environment variable configuration step for task name: %s with id: %s", taskName, taskId))
	}

	randomContainerName := common.GenerateRandomStringWithPrefix(3, false, true, false,
		"rand-cont-")

	t := &Task{
		// Identifiers
		Id:    taskId,
		Name:  taskName,
		Stack: stackName,
		// Parent objects
		PipelineCfg:           p,
		JobCfg:                job,
		EnvVarsInheritFromJob: envVars,

		// Default inherited container runtime from the job instantiated.
		ContainerImageDefault: job.ContainerImageURL,
		ContainerDefault:      job.ContainerDefault,
		ContainerNameDefault:  randomContainerName,

		Dirs: Dirs{
			RootDir:         ".",
			WorkDir:         init.WorkDir,
			MountDir:        init.MountDir,
			TargetDir:       init.TargetDir,
			RootDirDagger:   job.RootDir,
			WorkDirDagger:   job.WorkDir,
			MountDirDagger:  job.MountDir,
			TargetDirDagger: job.TargetDir,
		},

		PreReqs: PreRequisites{
			Files: []string{"Dockerfile"},
		},

		Actions: Actions{
			CustomCommands:  actions,
			DefaultCommands: []string{"docker", "build", "-t", randomContainerName, "."},
		},

		Ctx: job.Ctx,
	}

	p.UXMessage.ShowInfo("TASK-INIT",
		fmt.Sprintf("Successfully initialised core task instance: %s with id: %s",
			taskName, taskId))

	return t
}
