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
		input_channel:    make(chan models.Job),
		progress_channel: make(chan models.JobProgressReport),
		jobsRepository:   jobsRepository,
	}
}

func (jm *JobManager) CreateJob(newJob models.NewJobRequest) models.JobId {
	// parse newjobrequest into a NewJob struct
	job, _ := jm.jobsRepository.CreateJob(models.NewJob{})
	jm.input_channel <- *job
	//TODO: chequear si ahy metodo para confirmar encolado
	return job.JobId
}

func (jm *JobManager) spawnJob(newJob models.Job) {
	jm.jobsRepository.UpdateJobStatus(newJob.JobId, models.StatusRunning)
	spawn(models.NewJob{JobId: newJob.JobId}, jm.output_channel, jm.progress_channel)
}

func (jm *JobManager) FindJob(jobId models.JobId) (*models.Job, models.JobError) {
	return jm.jobsRepository.FindJob(jobId)
}

func (jm *JobManager) Run() {
	for {
		select {
		case msg := <-jm.input_channel:
			jm.spawnJob(msg)
		case msg := <-jm.progress_channel:
			jm.jobsRepository.UpdateJobProgress(msg.JobId, models.JobProgress(msg.Progress))
		case msg := <-jm.output_channel:
			jm.jobsRepository.UpdateJobStatus(msg.JobId, models.StatusFinished)
			jm.jobsRepository.SaveJobOutput(msg.JobId, models.JobOutput(msg.Output))
		}
	}
}
