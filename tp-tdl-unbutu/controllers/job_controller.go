package controllers

import (
	"net/http"
	"strconv"
	"tp-tdl-unbutu/tp-tdl-unbutu/services"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	manager *services.JobManager
}

func NewJobController(manager *services.JobManager) JobController {
	return JobController{manager}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	jobId := jc.manager.CreateJob(services.NewJobRequest{})
	c.IndentedJSON(http.StatusOK, jobId)
}

func (jc *JobController) GetJobStatus(c *gin.Context) {
	jobId, err := strconv.ParseInt(c.Param("jobId"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		jobStatus, err := jc.manager.GetJobStatus(services.JobId(jobId))
		if err != services.NoError {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			c.IndentedJSON(http.StatusOK, jobStatus)
		}
	}
}

func (jc *JobController) GetJobOutput(c *gin.Context) {
	jobId, err := strconv.ParseInt(c.Param("jobId"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		jobStatus, err := jc.manager.GetJobOutput(services.JobId(jobId))
		if err != services.NoError {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			c.IndentedJSON(http.StatusOK, jobStatus)
		}
	}
}

func (jc *JobController) RegisterRoutes(router *gin.Engine, base string) {
	router.GET(base+"/date", jc.CreateJob)
	router.GET(base+"/date/:jobId/status", jc.GetJobStatus)
	router.GET(base+"/date/:jobId/output", jc.GetJobOutput)
}
