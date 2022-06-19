package main

type JobError string

const (
	NoError        = "no_error"
	JobNotFound    = "job_not_found"
	JobNotFinished = "job_not_finished"
	JobFinished    = "job_finished"
)

type JobId int64

type Job struct {
	job_id          JobId
	currentProgress string
	output          string
	hasFinished     bool
}

type NewJobRequest struct{}

type JobManager struct {
	jobs             map[JobId]Job
	output_channel   chan JobResult
	progress_channel chan JobProgress
	last_job_id      JobId
}

func NewJobManager() *JobManager {
	return &JobManager{
		output_channel:   make(chan JobResult),
		progress_channel: make(chan JobProgress),
		last_job_id:      0,
		jobs:             make(map[JobId]Job),
	}
}

func (jm *JobManager) CreateJob(newJob NewJobRequest) JobId {
	job := Job{job_id: jm.last_job_id + 1, hasFinished: false}
	jm.last_job_id += 1
	jm.jobs[jm.last_job_id] = job
	spawn(job.job_id, jm.output_channel, jm.progress_channel)
	return job.job_id
}

func (jm *JobManager) GetJobStatus(jobId JobId) (string, JobError) {
	job, exists := jm.jobs[jobId]
	if !exists {
		return "", JobNotFound
	} else if job.hasFinished {
		return job.currentProgress, JobFinished
	} else {
		return job.currentProgress, NoError
	}
}

func (jm *JobManager) GetJobOutput(jobId JobId) (string, JobError) {
	job, exists := jm.jobs[jobId]
	if !exists {
		return "", JobNotFound
	} else if job.hasFinished {
		return job.output, NoError
	} else {
		return "", JobNotFinished
	}
}

func (jm *JobManager) Run() {
	for {
		select {
		case msg := <-jm.progress_channel:
			job := jm.jobs[msg.JobId]
			job.currentProgress = msg.Progress
			jm.jobs[msg.JobId] = job
		case msg := <-jm.output_channel:
			job := jm.jobs[msg.JobId]
			job.output = msg.Output
			job.hasFinished = true
			jm.jobs[msg.JobId] = job
		}
	}
}
