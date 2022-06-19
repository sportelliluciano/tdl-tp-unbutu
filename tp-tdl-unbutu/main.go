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
	// var currentJob *job
	// for {
	// 	req := <-channel
	// 	if req.Code == Spawn {
	// 		if currentJob == nil {
	// 			currentJob = spawn()
	// 			channel <- request{Code: Ok, Result: "job started"}
	// 		} else {
	// 			channel <- request{Code: Error, Result: "job is already running"}
	// 		}
	// 	} else if req.Code == Progress {
	// 		if currentJob != nil {
	// 			channel <- request{Code: Ok, Result: "running: " + currentJob.currentProgress + "%"}
	// 		} else {
	// 			channel <- request{Code: Ok, Result: "not running"}
	// 		}
	// 	} else if req.Code == Output {
	// 		if currentJob != nil {
	// 			if currentJob.output != nil {
	// 				channel <- request{Code: Ok, Result: *currentJob.output}
	// 				currentJob = nil
	// 			} else {
	// 				channel <- request{Code: Error, Result: "job still running"}
	// 			}
	// 		} else {
	// 			channel <- request{Code: Error, Result: "job is not running"}
	// 		}
	// 	}
	// }
}
