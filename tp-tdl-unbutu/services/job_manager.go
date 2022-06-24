package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobError string

const (
	NoError        = "no_error"
	JobNotFound    = "job_not_found"
	JobNotFinished = "job_not_finished"
	JobFinished    = "job_finished"
)

type JobId int64

type Job struct {
	job_id          JobId  `bson:"job_id"`
	currentProgress string `bson:"currentProgress"`
	output          string `bson:"output"`
	hasFinished     bool   `bson:"hasFinished"`
}

type NewJobRequest struct{}

type JobManager struct {
	input_channel    chan Job
	output_channel   chan JobResult
	progress_channel chan JobProgress
	last_job_id      JobId
	jobs             *mongo.Collection
}

func NewJobManager(collection *mongo.Collection) *JobManager {
	return &JobManager{
		output_channel:   make(chan JobResult),
		input_channel:    make(chan Job),
		progress_channel: make(chan JobProgress),
		last_job_id:      0,
		jobs:             collection,
	}
}

func (jm *JobManager) CreateJob(newJob NewJobRequest) JobId {

	job := Job{job_id: jm.last_job_id + 1, hasFinished: false}
	jm.last_job_id += 1
	jm.input_channel <- job
	//TODO: chequear si ahy metodo para confirmar encolado
	return job.job_id
}

func (jm *JobManager) spawnJob(newJob Job) {
	doc := bson.D{{"job_id", newJob.job_id}, {"currentProgress", ""}, {"output", ""}, {"hasFinished", false}}
	_, err := jm.jobs.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
	spawn(NewJob{JobId: newJob.job_id}, jm.output_channel, jm.progress_channel)
}

func (jm *JobManager) GetJobStatus(jobId JobId) (string, JobError) {
	var job bson.M
	err := jm.jobs.FindOne(context.TODO(), bson.D{{Key: "job_id", Value: jobId}}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return "", JobNotFound
		}
		panic(err)
	}
	hasFinished, ok_1 := job["hasFinished"].(bool)
	currentProgress, ok_2 := job["currentProgress"].(string)

	if !(ok_1 && ok_2) {
		panic("Esto es una cagada")
	}

	if hasFinished {
		return currentProgress, JobFinished
	} else {
		return currentProgress, NoError
	}
}

func (jm *JobManager) GetJobOutput(jobId JobId) (string, JobError) {
	var job bson.M
	err := jm.jobs.FindOne(context.TODO(), bson.D{{Key: "job_id", Value: jobId}}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return "", JobNotFound
		}
		panic(err)
	}

	hasFinished, ok_1 := job["hasFinished"].(bool)
	output, ok_2 := job["output"].(string)

	if !(ok_1 && ok_2) {
		panic("Esto es una cagada")
	}

	if hasFinished {
		return output, NoError
	} else {
		return "", JobNotFinished
	}
}

func (jm *JobManager) Run() {
	for {
		select {
		case msg := <-jm.input_channel:
			jm.spawnJob(msg)
		case msg := <-jm.progress_channel:
			update := bson.D{{Key: "$set", Value: bson.D{{Key: "currentProgress", Value: msg.Progress}}}}
			jm.jobs.FindOneAndUpdate(context.TODO(), bson.D{{Key: "job_id", Value: msg.JobId}}, update)
		case msg := <-jm.output_channel:

			update := bson.D{{Key: "$set", Value: bson.D{{Key: "output", Value: msg.Output}, {Key: "hasFinished", Value: true}}}}
			jm.jobs.FindOneAndUpdate(context.TODO(), bson.D{{Key: "job_id", Value: msg.JobId}}, update)
		}
	}
}
