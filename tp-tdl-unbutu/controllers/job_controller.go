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

func NewJobController(manager *services.JobManager) JobController {
	return JobController{manager}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	jobId := jc.manager.CreateJob(models.NewJobRequest{})
	c.IndentedJSON(http.StatusOK, jobId)
}

func (jc *JobController) GetJob(c *gin.Context) {
	job, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusOK, job)
	}
}

func (jc *JobController) RegisterRoutes(router *gin.Engine, base string) {
	router.POST(base+"/job", jc.CreateJob)
	router.GET(base+"/job/:jobId", jc.GetJob)
}
