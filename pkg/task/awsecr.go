package task

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
)

func RunTaskAWSECR(opt InitOptions) error {
	taskSelector := common.NormaliseStringUpper(opt.Task)
	taskPrefix := "AWS:ECR"

	p := opt.PipelineCfg
	j := opt.JobCfg

	actionCMDs := opt.ActionCommands

	switch taskSelector {
	case "PUSH":
		actionPrefix := fmt.Sprintf("%s:%s", taskPrefix, taskSelector)
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskAWSECR(c, actionCMDs, &opt, actionPrefix)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewAWSECRPushAction(t, actionPrefix)

		// Run the action
		_, err := a.DeployNewTask()
		if err != nil {
			return err
		}
	}
	return nil
}
