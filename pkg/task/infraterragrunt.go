package task

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
)

var allowedTasks = []string{"PLAN", "APPLY", "DESTROY", "PLAN-ALL", "APPLY-ALL", "DESTROY-ALL"}

func RunTaskInfraTerraGrunt(opt InitOptions) error {
	taskSelector := common.NormaliseStringUpper(opt.Task)

	// Check if the task is allowed
	if !common.IsStringInSlice(taskSelector, allowedTasks) {
		return errors.NewArgumentError(fmt.Sprintf("Task '%s' is not allowed. "+
			"Allowed tasks are: %s", taskSelector, allowedTasks), nil)
	}

	taskPrefix := opt.JobCfg.Stack

	p := opt.PipelineCfg
	j := opt.JobCfg

	actionCMDs := opt.ActionCommands
	actionPrefix := fmt.Sprintf("%s:%s", taskPrefix, taskSelector)

	switch taskSelector {
	case "PLAN":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskInfraTerraGrunt(c, actionCMDs, &opt, actionPrefix)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewInfraTerraGruntAction(t, actionPrefix)

		// Run the action
		_, err := a.Plan()
		if err != nil {
			return err
		}

	case "APPLY":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskInfraTerraGrunt(c, actionCMDs, &opt, actionPrefix)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewInfraTerraGruntAction(t, actionPrefix)

		// Run the action
		_, err := a.Apply()
		if err != nil {
			return err
		}

	case "DESTROY":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskInfraTerraGrunt(c, actionCMDs, &opt, actionPrefix)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewInfraTerraGruntAction(t, actionPrefix)

		// Run the action
		_, err := a.Destroy()
		if err != nil {
			return err
		}

	case "VALIDATE":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskInfraTerraGrunt(c, actionCMDs, &opt, actionPrefix)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewInfraTerraGruntAction(t, actionPrefix)

		// Run the action
		_, err := a.Validate()
		if err != nil {
			return err
		}

	}

	return nil
}
