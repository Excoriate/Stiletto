package job

import "fmt"

func GetErrMsg(jobName, jobId, msg string, err error) string {
	errMsgPrefix := fmt.Sprintf("Job %s - id: %s failed with error: ", jobName, jobId)

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

func GetInfoMsg(jobName, jobId, msg string) string {
	infoMsgPrefix := fmt.Sprintf("Job %s - id: %s info: ", jobName, jobId)

	return fmt.Sprintf("%s%s", infoMsgPrefix, msg)
}
