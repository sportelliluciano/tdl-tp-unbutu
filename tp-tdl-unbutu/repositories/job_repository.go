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
	_, err := jr.collection.InsertOne(context.TODO(), serializeNewJobToBson(string(newJob.JobId), newJob))
	if err != nil {
		return nil, models.DatabaseError
	}
	return &models.Job{
		JobId:    models.JobId(newJob.JobId),
		Progress: models.JobProgress(make(map[string]interface{})),
		Status:   models.StatusQueued,
		Output:   models.JobOutput(""),
		Format:   models.Format(newJob.Format),
	}, models.NoError
}

func (jr *JobRepository) FindJob(jobId models.JobId) (*models.Job, models.JobError) {
	var job bson.M
	err := jr.collection.FindOne(context.TODO(), bson.D{{Key: "jobId", Value: jobId}}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.JobNotFound
		}
		return nil, models.DatabaseError
	}
	jobDeserialized, isValid := deserializeJobFromBson(job)
	if isValid {
		return jobDeserialized, models.NoError
	} else {
		return nil, models.JobMetadataCorrupted
	}
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

func serializeNewJobToBson(newId string, newJob models.NewJob) bson.D {
	return bson.D{
		{Key: "jobId", Value: newId},
		{Key: "progress", Value: make(map[string]string)},
		{Key: "status", Value: models.StatusQueued},
		{Key: "output", Value: ""},
		{Key: "format", Value: newJob.Format},
	}
}

func deserializeJobFromBson(job bson.M) (*models.Job, bool) {
	jobId, ok_0 := job["jobId"].(string)
	progress, ok_1 := job["progress"].(bson.M)
	status, ok_2 := job["status"].(string)
	output, ok_3 := job["output"].(string)
	format, ok_4 := job["format"].(string)

	if !ok_0 || !ok_1 || !ok_2 || !ok_3 || !ok_4 {
		println(ok_0)
		println(ok_1)
		println(ok_2)
		println(ok_3)
		println(ok_4)
		return nil, false
	}

	return &models.Job{
		JobId:    models.JobId(jobId),
		Progress: models.JobProgress(progress),
		Status:   models.JobStatus(status),
		Output:   models.JobOutput(output),
		Format:   models.Format(format),
	}, true
}
