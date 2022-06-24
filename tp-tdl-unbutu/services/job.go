package services

import (
	"io"

	"os/exec"

	"strings"
)

type NewJob struct {
	JobId JobId
}

type JobResult struct {
	JobId  JobId
	Output string
}

type JobProgress struct {
	JobId    JobId
	Progress string
}

func spawn(newJob NewJob, output_channel chan JobResult, progress_channel chan JobProgress) {
	dateCmd := exec.Command("bash", "./date-with-sleep.sh")
	pipeReader, err := dateCmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	go reportJobProgress(pipeReader, progress_channel, newJob.JobId)
	go waitForJobOutput(dateCmd, output_channel, newJob.JobId)
}

func reportJobProgress(pipeReader io.ReadCloser, progress_channel chan JobProgress, job_id JobId) {
	for {
		var buffer = make([]byte, 100)
		n, err := pipeReader.Read(buffer)
		if n == 0 || err != nil {
			break
		}

		currentProgress := strings.Trim(string(buffer[0:n]), "\n")
		progress_channel <- JobProgress{JobId: job_id, Progress: currentProgress}
	}
}

func waitForJobOutput(cmd *exec.Cmd, output_channel chan JobResult, job_id JobId) {
	dateOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	output := strings.Trim(string(dateOut), "\n")
	output_channel <- JobResult{JobId: job_id, Output: output}
}
