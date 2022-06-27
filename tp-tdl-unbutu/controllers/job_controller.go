package controllers

import (
	"io"
	"net/http"
	"os"
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

func (jc *JobController) GetJob(c *gin.Context) {
	job, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusCreated, job)
	}
}

func (jc *JobController) GetFile(c *gin.Context) {
	_, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	var format string = ".txt"
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=hola.txt")
		c.Header("Content-Type", "application/octet-stream")
		c.File("./output/" + format)
		c.IndentedJSON(http.StatusOK, "Download successful ")
	}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	newId := uuid.New().String()
	jobId, err := jc.manager.CreateJob(models.NewJobRequest{JobId: models.JobId(newId), Format: models.Format(c.Query("outputFormat"))})
	if err != models.NoError {
		c.IndentedJSON(http.StatusTooManyRequests, jobId)
	} else {
		// Source
		file, err := c.FormFile("file")
		if err == nil {
			src, err := file.Open()
			if err != nil {
				defer src.Close()

				// Destination
				dst, err := os.Create("./input/" + newId + ".txt")
				if err != nil {
					defer dst.Close()

					// Copy
					io.Copy(dst, src)
					c.IndentedJSON(http.StatusOK, "Upload successful ")
				} else {
					c.IndentedJSON(http.StatusBadRequest, err)
				}

			} else {
				c.IndentedJSON(http.StatusBadRequest, err)
			}

		} else {
			c.IndentedJSON(http.StatusBadRequest, err)
		}
		c.IndentedJSON(http.StatusCreated, jobId)
	}
}

func (jc *JobController) RegisterRoutes(router *gin.Engine, base string) {
	router.POST(base+"/job", jc.CreateJob)
	router.GET(base+"/job/:jobId", jc.GetJob)
	router.GET(base+"/job/:jobId/file", jc.GetFile)
}
