package services

import (
	"log"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"
)

type Worker struct {
	input_channel    chan models.Job
	output_channel   chan models.JobResult
	progress_channel chan models.JobProgressReport
}

func NewWorker(result chan models.JobResult, job chan models.Job, progress chan models.JobProgressReport) *Worker {
	return &Worker{
		output_channel:   result,
		input_channel:    job,
		progress_channel: progress,
	}
}

func (w *Worker) executeJob(newJob models.Job) {
	spawn(models.NewJob{JobId: newJob.JobId}, w.output_channel, w.progress_channel)
}

func (w *Worker) Run() {
	for {
		log.Println("waiting for job")
		msg := <-w.input_channel
		log.Println("starting job")
		w.executeJob(msg)
		log.Println("job finished")
	}
}
