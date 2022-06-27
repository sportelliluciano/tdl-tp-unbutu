package services

import (
	"tp-tdl-unbutu/tp-tdl-unbutu/models"
	"tp-tdl-unbutu/tp-tdl-unbutu/repositories"
)

type JobManager struct {
	input_channel    chan models.Job
	output_channel   chan models.JobResult
	progress_channel chan models.JobProgressReport
	jobsRepository   *repositories.JobRepository
}

func NewJobManager(jobsRepository *repositories.JobRepository) *JobManager {
	// W1(J1), W2(J2), BACKLOG(channel size=1)=J3; J4 falla porque size = 1 y lleno
	return &JobManager{
		output_channel:   make(chan models.JobResult),
		input_channel:    make(chan models.Job, 4),
		progress_channel: make(chan models.JobProgressReport),
		jobsRepository:   jobsRepository,
	}
}

func (jm *JobManager) CreateJob(newJob models.NewJobRequest) (*models.JobId, models.JobError) {
	// parse newjobrequest into a NewJob struct
	job, _ := jm.jobsRepository.CreateJob(models.NewJob(newJob))
	select {
	case jm.input_channel <- *job:
	default:
		return nil, models.QueueFull
	}
	return &job.JobId, models.NoError
}

func (jm *JobManager) FindJob(jobId models.JobId) (*models.Job, models.JobError) {
	return jm.jobsRepository.FindJob(jobId)
}

func (jm *JobManager) createWorkerPool(n_workers int) {
	for i := 0; i < n_workers; i++ {
		worker := NewWorker(jm.output_channel, jm.input_channel, jm.progress_channel)
		go worker.Run(jm.jobsRepository)
	}
}

func (jm *JobManager) Run() {
	jm.createWorkerPool(3)
	for {
		select {
		case msg := <-jm.progress_channel:
			jm.jobsRepository.UpdateJobProgress(msg.JobId, models.JobProgress(msg.Progress))
		case msg := <-jm.output_channel:
			if msg.Output == models.JobFail {
				jm.jobsRepository.UpdateJobStatus(msg.JobId, models.StatusFailed)
			} else {
				jm.jobsRepository.UpdateJobStatus(msg.JobId, models.StatusFinished)
			}
			jm.jobsRepository.SaveJobOutput(msg.JobId, models.JobOutput(msg.Output))
		}
	}
}
