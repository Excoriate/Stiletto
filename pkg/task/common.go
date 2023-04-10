package task

import (
	"fmt"
)

func GetErrMsg(t *Task, msg string, err error) string {
	taskName := t.Name
	taskId := t.Id
	errMsgPrefix := fmt.Sprintf("Task %s - id: %s failed with error: ", taskName, taskId)

	// If msg is passed empty, it'll use the error message.
	if msg == "" {
		return fmt.Sprintf("%s%s", errMsgPrefix, err.Error())
	}

	// if the msg is passed, and also the error
	if err != nil {
		return fmt.Sprintf("%s%s: %s", errMsgPrefix, msg, err.Error())
	}

	// if the msg is passed, but no error
	return fmt.Sprintf("%s%s", errMsgPrefix, msg)
}

func GetInfoMsg(t *Task, msg string) string {
	taskName := t.Name
	taskId := t.Id
	infoMsgPrefix := fmt.Sprintf("Task %s - id: %s info: ", taskName, taskId)

	return fmt.Sprintf("%s%s", infoMsgPrefix, msg)
}
