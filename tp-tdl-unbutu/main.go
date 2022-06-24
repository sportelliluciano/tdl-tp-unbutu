package main

import (
	"context"
	"log"
	"os"
	"tp-tdl-unbutu/tp-tdl-unbutu/controllers"
	"tp-tdl-unbutu/tp-tdl-unbutu/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	manager := services.NewJobManager(coll)
	controller := controllers.NewJobController(manager)
	controller.RegisterRoutes(router, "/api")
	router.Static("/ui", "./ui")

	go router.Run("localhost:8080")
	manager.Run()
	// W1(J1), W2(J2), BACKLOG(channel size=1)=J3; J4 falla porque size = 1 y lleno
}
