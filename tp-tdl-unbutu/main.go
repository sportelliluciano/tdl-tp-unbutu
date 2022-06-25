package main

import (
	"tp-tdl-unbutu/tp-tdl-unbutu/controllers"
	"tp-tdl-unbutu/tp-tdl-unbutu/repositories"
	"tp-tdl-unbutu/tp-tdl-unbutu/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config := LoadConfigFromEnv()

	db := ConnectToMongoDb(config.MongoDbUri)
	defer db.Disconnect()

	router := gin.Default()

	// Start repositories, services & controllers
	jobsRepository := repositories.NewJobRepository(db.Db.Collection("jobs"))
	manager := services.NewJobManager(&jobsRepository)
	controller := controllers.NewJobController(manager)

	// Register routes
	controller.RegisterRoutes(router, "/api")
	router.Static("/ui", "./ui")

	// Start http server in go routine
	go router.Run(config.Host + ":" + config.Port)

	// Start job manager in main thread
	manager.Run()
}
