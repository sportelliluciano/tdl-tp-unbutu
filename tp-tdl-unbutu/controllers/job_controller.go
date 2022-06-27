package controllers

import (
	"net/http"
	"tp-tdl-unbutu/tp-tdl-unbutu/models"
	"tp-tdl-unbutu/tp-tdl-unbutu/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobController struct {
	manager *services.JobManager
}

func NewJobController(manager *services.JobManager) JobController {
	return JobController{manager}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	newId := uuid.New().String()
	jobId, err := jc.manager.CreateJob(models.NewJobRequest{JobId: models.JobId(newId), Format: models.Format("mp3")})
	if err != models.NoError {
		c.IndentedJSON(http.StatusTooManyRequests, jobId)
	} else {
		c.IndentedJSON(http.StatusCreated, jobId)
	}
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
