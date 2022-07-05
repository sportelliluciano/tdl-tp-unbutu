package repositories

import (
	"context"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(collection *mongo.Collection) JobRepository {
	return JobRepository{collection: collection}
}

func (jr *JobRepository) CreateJob(newJob models.NewJob) (*models.Job, models.JobError) {
	job := models.Job{
		JobId:    models.JobId(newJob.JobId),
		Progress: models.JobProgress(make(map[string]interface{})),
		Status:   models.StatusQueued,
		Output:   models.JobOutput(""),
		Format:   models.Format(newJob.Format),
	}

	_, err := jr.collection.InsertOne(context.TODO(), job)
	if err != nil {
		return nil, models.DatabaseError
	}
	return &job, models.NoError
}

func (jr *JobRepository) FindJob(jobId models.JobId) (*models.Job, models.JobError) {
	var job models.Job
	err := jr.collection.FindOne(context.TODO(), bson.D{{Key: "jobId", Value: jobId}}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.JobNotFound
		}
		return nil, models.DatabaseError
	}
	
	return &job, models.NoError
}

func (jr *JobRepository) UpdateJobProgress(jobId models.JobId, newProgress models.JobProgress) {
	jr.collection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{Key: "jobId", Value: jobId}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "progress", Value: newProgress}}}},
	)
}

func (jr *JobRepository) UpdateJobStatus(jobId models.JobId, newStatus models.JobStatus) {
	jr.collection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{Key: "jobId", Value: jobId}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: newStatus}}}},
	)
}

func (jr *JobRepository) SaveJobOutput(jobId models.JobId, output models.JobOutput) {
	jr.collection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{Key: "jobId", Value: jobId}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "output", Value: output}}}},
	)
}