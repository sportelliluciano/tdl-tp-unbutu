package services

import (
	"io"
	"log"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"

	"os/exec"

	"strings"
)

func spawn(newJob models.NewJob, output_channel chan models.JobResult, progress_channel chan models.JobProgressReport) {
	inputFile := "./input/" + newJob.JobId
	outputFile := "./output/" + string(newJob.JobId) + "." + string(newJob.Format)

	log.Println("input: " + inputFile)
	log.Println("output: " + outputFile)

	dateCmd := exec.Command("bash", "./date-with-sleep.sh")
	pipeReader, err := dateCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go reportJobProgress(pipeReader, progress_channel, newJob.JobId)
	waitForJobOutput(dateCmd, output_channel, newJob.JobId)
}

func reportJobProgress(pipeReader io.ReadCloser, progress_channel chan models.JobProgressReport, job_id models.JobId) {
	for {
		var buffer = make([]byte, 8192)
		n, err := pipeReader.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		lines := strings.Split(string(buffer[0:n]), "\n")
		currentProgress := make(map[string]interface{})
		for _, line := range lines {
			values := strings.Split(line, "=")
			if len(values) == 2 {
				currentProgress[values[0]] = values[1]
			}
		}
		progress_channel <- models.JobProgressReport{JobId: job_id, Progress: currentProgress}
	}
}

func waitForJobOutput(cmd *exec.Cmd, output_channel chan models.JobResult, job_id models.JobId) {
	cmd.Start()
	err := cmd.Wait()
	output := models.NoError
	if err != nil {
		output = models.JobFail
	}
	output_channel <- models.JobResult{JobId: job_id, Output: output}
}
