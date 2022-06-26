package models

type JobId string
type JobStatus string
type JobProgress map[string]interface{}
type JobOutput string

type JobError string

const (
	NoError              = "no_error"
	JobNotFound          = "job_not_found"
	JobNotFinished       = "job_not_finished"
	JobFinished          = "job_finished"
	JobMetadataCorrupted = "job_metadata_corrupted"
	DatabaseError        = "database_error"
	QueueFull            = "queue_full"
)

const (
	StatusQueued   = "queued"
	StatusRunning  = "running"
	StatusFinished = "finished"
	StatusFailed   = "failed"
)

type NewJob struct {
	JobId JobId
}

type JobResult struct {
	JobId  JobId
	Output string
}

type JobProgressReport struct {
	JobId    JobId
	Progress map[string]interface{}
}

type Job struct {
	JobId    JobId       `bson:"jobId" json:"jobId"`
	Progress JobProgress `bson:"progress" json:"progress"`
	Status   JobStatus   `bson:"status" json:"status"`
	Output   JobOutput   `bson:"output" json:"output"`
}

func (j *Job) HasFinished() bool {
	return j.Status == StatusFinished
}

type NewJobRequest struct{}
