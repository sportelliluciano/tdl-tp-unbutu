package controllers

import (
	"net/http"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"
	"tp-tdl-unbutu/tp-tdl-unbutu/services"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	manager *services.JobManager
}

type jobStatusReport struct {
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

func NewJobController(manager *services.JobManager) JobController {
	return JobController{manager}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	jobId := jc.manager.CreateJob(models.NewJobRequest{})
	c.IndentedJSON(http.StatusOK, jobId)
}

func (jc *JobController) GetJobStatus(c *gin.Context) {
	job, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusOK, jobStatusReport{
			Status:   (string)(job.Status),
			Progress: (string)(job.Progress),
		})
	}

}

func (jc *JobController) GetJobOutput(c *gin.Context) {
	job, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusOK, job.Output)
	}
}

func (jc *JobController) RegisterRoutes(router *gin.Engine, base string) {
	router.GET(base+"/date", jc.CreateJob)
	router.GET(base+"/date/:jobId/status", jc.GetJobStatus)
	router.GET(base+"/date/:jobId/output", jc.GetJobOutput)
}
