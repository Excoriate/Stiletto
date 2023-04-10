package task

import (
	"github.com/Excoriate/stiletto/internal/common"
)

// RunTaskDocker is the entry point for all Docker tasks.
func RunTaskDocker(opt InitOptions) error {
	taskSelector := common.NormaliseStringUpper(opt.Task)

	p := opt.PipelineCfg
	j := opt.JobCfg

	actionCMDs := opt.ActionCommands

	switch taskSelector {
	case "BUILD":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskDocker(c, actionCMDs, &opt, "DOCKER-BUILD")

		// New action to execute (mapped to the --task passed from the command line)
		a := NewDockerBuildAction(t)

		// Run the action
		_, err := a.BuildTagAndPush("Dockerfile")
		if err != nil {
			return err
		}
	}
	return nil
}
