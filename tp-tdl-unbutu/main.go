package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"strconv"
)

// API REST EN GO: https://go.dev/doc/tutorial/web-service-gin
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("TpTdl").Collection("jobs")
	log.Println("FUNCA")
	router := gin.Default()
	manager := NewJobManager(coll)
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
