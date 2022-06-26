package services

import (
	"io"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"

	"os/exec"

	"strings"
)

func spawn(newJob models.NewJob, output_channel chan models.JobResult, progress_channel chan models.JobProgressReport) {
	dateCmd := exec.Command("bash", "./date-with-sleep.sh")
	pipeReader, err := dateCmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	go reportJobProgress(pipeReader, progress_channel, newJob.JobId)
	waitForJobOutput(dateCmd, output_channel, newJob.JobId)
}

func reportJobProgress(pipeReader io.ReadCloser, progress_channel chan models.JobProgressReport, job_id models.JobId) {
	for {
		var buffer = make([]byte, 100)
		n, err := pipeReader.Read(buffer)
		if n == 0 || err != nil {
			break
		}

		currentProgress := strings.Trim(string(buffer[0:n]), "\n")
		progress_channel <- models.JobProgressReport{JobId: job_id, Progress: currentProgress}
	}
}

func waitForJobOutput(cmd *exec.Cmd, output_channel chan models.JobResult, job_id models.JobId) {
	dateOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	output := strings.Trim(string(dateOut), "\n")
	output_channel <- models.JobResult{JobId: job_id, Output: output}
}
