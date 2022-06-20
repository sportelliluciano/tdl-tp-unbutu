package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
)

// API REST EN GO: https://go.dev/doc/tutorial/web-service-gin
func main() {
	router := gin.Default()
	manager := NewJobManager()
	router.GET("/date", func(c *gin.Context) {
		jobId := manager.CreateJob(NewJobRequest{})
		c.IndentedJSON(http.StatusOK, jobId)
	})
	router.GET("/date/:jobId/status", func(c *gin.Context) {
		jobId, err := strconv.ParseInt(c.Param("jobId"), 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			jobStatus, err := manager.GetJobStatus(JobId(jobId))
			if err != NoError {
				c.IndentedJSON(http.StatusBadRequest, err)
			} else {
				c.IndentedJSON(http.StatusOK, jobStatus)
			}
		}
	})
	router.GET("/date/:jobId/output", func(c *gin.Context) {
		jobId, err := strconv.ParseInt(c.Param("jobId"), 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			jobStatus, err := manager.GetJobOutput(JobId(jobId))
			if err != NoError {
				c.IndentedJSON(http.StatusBadRequest, err)
			} else {
				c.IndentedJSON(http.StatusOK, jobStatus)
			}
		}
	})
	router.Static("/ui", "./ui")

	go router.Run("localhost:8080")
	manager.Run()
	// W1(J1), W2(J2), BACKLOG(channel size=1)=J3; J4 falla porque size = 1 y lleno
}
