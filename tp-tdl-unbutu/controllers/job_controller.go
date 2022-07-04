package controllers

import (
	"io"
	"mime/multipart"
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
		c.IndentedJSON(http.StatusOK, job)
	}
}

func (jc *JobController) GetFile(c *gin.Context) {
	job, err := jc.manager.FindJob(models.JobId(c.Param("jobId")))
	if err != models.NoError {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=result."+string(job.Format))
		c.Header("Content-Type", "application/octet-stream")
		c.File("./output/" + string(job.JobId) + "." + string(job.Format))
	}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	newId := uuid.New().String()
	jobId, err := jc.manager.CreateJob(models.NewJobRequest{JobId: models.JobId(newId), Format: models.Format(c.Query("outputFormat"))})
	if err != models.NoError {
		c.IndentedJSON(http.StatusTooManyRequests, jobId)
	} else {
		file, err := c.FormFile("file")
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		err = handleFileUpload(file, newId)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			c.IndentedJSON(http.StatusCreated, jobId)
		}
	}
}
func handleFileUpload(file *multipart.FileHeader, newId string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create("./input/" + newId)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (jc *JobController) RegisterRoutes(router *gin.Engine, base string) {
	router.POST(base+"/job", jc.CreateJob)
	router.GET(base+"/job/:jobId", jc.GetJob)
	router.GET(base+"/job/:jobId/file", jc.GetFile)
}
